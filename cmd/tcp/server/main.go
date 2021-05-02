package main

import (
	"github.com/scottcagno/net-tools/pkg/tcp/server"
	"io"
	"log"
	"net"
)

func main() {

	// run example one
	ServerExample1()

	// run example two
	//ServerExample2()

	// run example three
	//ServerExample3()

	// run example four
	//ServerExample4()
}

func ServerExample1() {
	err := server.ListenAndServe(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}

func ServerExample2() {
	err := server.ListenAndServeTCP(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}

func ServerExample3() {
	err := server.ListenAndServeTCPWithHandler(":8080", server.HandleEcho())
	if err != nil {
		log.Fatalln(err)
	}
}

func ServerExample4() {
	// the power of closures
	customFn := func(conn net.Conn) {
		defer conn.Close()
		// echo all incoming data
		io.Copy(conn, conn)
	}
	err := server.ListenAndServeTCPWithHandler(":8080", customFn)
	if err != nil {
		log.Fatalln(err)
	}
}
