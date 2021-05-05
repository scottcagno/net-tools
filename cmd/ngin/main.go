package main

import (
	"github.com/scottcagno/net-tools/pkg/ngin"
	"log"
)

func main() {

	// open engine
	en, err := ngin.OpenEngine("./test/data")
	if err != nil {
		log.Fatal(err)
	}

	n, err := en.Write([]byte(""))

	// don't forget to close the engine!
	defer func() {
		err = en.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
