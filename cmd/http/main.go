package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/web"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/index", getIndex())
	mux.Handle("/home", getHome())
	mux.Handle("/login", getLogin())

	// add a logger
	chain := web.Logger(mux)

	// server
	err := web.ListenAndServe(":8080", chain)
	log.Fatal(err)
}

func getIndex() http.Handler {
	// only allow GET
	return web.GetFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /index hit!")
	})
}

func getHome() http.Handler {
	// only allow GET
	return web.GetFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /home hit!")
	})
}

func getLogin() http.Handler {
	// only allow GET
	return web.GetFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /login hit!")
	})
}
