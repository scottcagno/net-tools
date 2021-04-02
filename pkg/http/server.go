package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	GET         = http.MethodGet
	HEAD        = http.MethodHead
	POST        = http.MethodPost
	PUT         = http.MethodPut
	DELETE      = http.MethodDelete
	ANY         = "ANY"
	STATIC_PATH = "static"
)

type Server struct {
	routes       []*Route
	static, errs http.Handler
	logger       *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	if _, err := os.Stat(STATIC_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(STATIC_PATH, 0755); err != nil {
			log.Fatalf("could not create static file path %q: %v\n", STATIC_PATH, err)
		}
	}
	if logger == nil {
		logger = log.Default()
	}
	s := &Server{
		routes: make([]*Route, 0),
		static: http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))),
		logger: logger,
	}
	s.Handle("/error/*", http.HandlerFunc(handleError))
	return s
}

func (s *Server) Handle(path string, handler http.Handler) {
	s.HandleMethod(ANY, path, handler)
}

func (s *Server) HandleMethod(method string, path string, handler http.Handler) {
	s.routes = append(s.routes, &Route{method, path, handler})
}

func (s *Server) HandleGet(path string, handler http.Handler) {
	s.HandleMethod(http.MethodGet, path, handler)
}

func (s *Server) HandlePost(path string, handler http.Handler) {
	s.HandleMethod(http.MethodPost, path, handler)
}

type Route struct {
	Method string
	Path   string
	http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions || r.URL.Path == "/favicon.ico" {
		return
	}
	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/static/") {
		s.static.ServeHTTP(w, r)
		return
	}
	if r.Method == http.MethodPost && !strings.Contains(r.Referer(), r.Host) {
		log.Printf("Error: %v\n", "hmm, i dunno")
		goto errNotFound
	}
	for _, route := range s.routes {
		if route.Method != r.Method || !strings.Contains(route.Method, "ANY") {
			continue // skip to next handler
		}
		found, err := filepath.Match(route.Path, r.URL.Path)
		if !found || err != nil && err == filepath.ErrBadPattern {
			log.Printf("Error: %v\n", err)
			goto errNotFound
		}
		route.ServeHTTP(w, r)
		return
	}
errNotFound:
	http.Redirect(w, r, "/error/404", http.StatusTemporaryRedirect)
	return
}

var handleFavicon = func() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})
}

var handleError = func(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.Atoi(r.FormValue(":code"))
	if err != nil {
		code = 500
	}
	w.Header().Set("Content-Type", "text/html; utf-8")
	errHtmlStr := `<html>
	<body>
		<head><title>%d</title></head>
		<center>
			<br/>
			<h1>HTTP Status %d %s</h1>
			<p>Default Error Handler</p>
		</center>
	</body>
</html>`
	fmt.Fprintf(w, errHtmlStr, code, code, http.StatusText(code))
	return
}

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}
