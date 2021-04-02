package main

import (
	http2 "github.com/scottcagno/net-tools/pkg/http"
	"github.com/scottcagno/net-tools/pkg/httputil"
	"log"
	"net/http"
)

func main() {

	httputil.InfoLogger.Println("Something noteworthy happened")
	httputil.WarningLogger.Println("There is something you should know about")
	httputil.ErrorLogger.Println("Something went wrong")

	mux := http2.NewServer(nil)

	// example middleware
	final := http.HandlerFunc(finalMiddleware)
	mux.HandleGet("/", middlewareOne(middlewareTwo(final)))

	// server
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		if r.URL.Path == "/foo" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}

func finalMiddleware(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	w.Write([]byte("OK"))
}
