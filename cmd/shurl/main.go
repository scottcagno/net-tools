package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/web"
	"log"
	"net/http"
	"os"
)

func main() {

	StoreTest()
	return
	s := NewSimpleStore("./")

	mux := http.NewServeMux()
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/", handleIndex(s))

	chain := web.Logger(mux)
	err := web.ListenAndServe(":8080", chain)
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(store Store) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// do stuff
		fmt.Fprintf(w, "index handler hit")
	}
	return http.HandlerFunc(fn)
}

func CreateDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0655); err != nil {
			fmt.Errorf("could not create static file path %q: %v\n", path, err)
		}
	}
	return nil
}

func CreateOrOpenFile(name string) *os.File {
	fd, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return fd
}
