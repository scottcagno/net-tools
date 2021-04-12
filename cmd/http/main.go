package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/index", http.HandlerFunc(getIndex))
	mux.Handle("/home", http.HandlerFunc(getHome))
	mux.Handle("/login", http.HandlerFunc(getOrPostLogin))

	// server
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func allowMethod(w http.ResponseWriter, r *http.Request, method ...string) {
	var allowed bool
	for _, m := range method {
		if m == r.Method {
			allowed = true
			break
		}
	}
	if allowed {
		return
	}
	code := http.StatusMethodNotAllowed
	http.Error(w, http.StatusText(code), code)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	allowMethod(w, r, http.MethodGet)
	fmt.Fprintf(w, "GET /index hit!")
}

func getHome(w http.ResponseWriter, r *http.Request) {
	allowMethod(w, r, http.MethodGet)
	fmt.Fprintf(w, "GET /home hit!")
}

func getOrPostLogin(w http.ResponseWriter, r *http.Request) {
	allowMethod(w, r, http.MethodGet, http.MethodPost)
	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "POST /login hit!")
		return
	}
	fmt.Fprintf(w, "GET /login hit!")
}
