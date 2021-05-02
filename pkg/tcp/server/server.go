package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// ListenAndServe creates a generic listening socket on the TCP protocol
// for the provided address and/port and handles TCP connections. You can
// in theory change the network type of this listener to a different
// protocol if you know how to properly handle that protocol.
func ListenAndServe(host string) error {
	// initialize a listening socket
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		// wait (block) for a connection, accepting it when it arrives.
		conn, err := ln.Accept()
		if err != nil {
			// if we can't accept a connection, then finish loop and wait for next connection.
			log.Printf("error accepting connection: %s\n", err)
			continue
		}
		// got a connection--lets hand it off on it's own goroutine and get back to waiting for
		// the next potential connection and do it all again.
		log.Printf("ACCEPTED CONNECTION: %q\n", conn.RemoteAddr())
		go HandleConn(conn)
	}
}

// HandleConn is a server handler that takes a generic net.Conn. A new one is
// created in a separate goroutine for every incoming connection that is accepted.
func HandleConn(conn net.Conn) {
	// make a buffered reader to read incoming data
	r := bufio.NewReader(conn)
	// make a buffered writer to write back to client
	w := bufio.NewWriter(conn)
	defer conn.Close()
	for {
		// read data until we get to the end of a line
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			// if we find the end of the data before
			// we find a newline char, then just stop
			break
		} else if err != nil {
			// if some other error happens, log and return
			log.Printf("Error reading: %s\n", err)
			HandleConnClose(conn)
			return
		}
		fmt.Printf("RECEIVED: %q FROM [%s], REPLYING...", data, conn.RemoteAddr())
		reply := fmt.Sprintf("ECHO: %q\n", data)
		if _, err = w.WriteString(reply); err != nil {
			log.Printf("Error writing: %s\n", err)
			fmt.Printf(" ERR!\n")
		} else {
			fmt.Printf(" OK!\n")
		}
	}
}

/*
 * This next section is a more specific type of server. It is STRICTLY a TCP
 * server only. It is mostly the same--the main differences are that this has
 * more specific methods you are able to access that you don't have available
 * to you with the more generic net.Listen and net.Conn.
 */

// ListenAndServeTCP creates a listening socket on the TCP protocol
// for the provided address and/port and handles TCP connections.
func ListenAndServeTCP(host string) error {
	// resolve tcp address
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}
	// initialize a listening socket
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		// wait (block) for a connection, accepting it when it arrives.
		conn, err := ln.AcceptTCP()
		if err != nil {
			// if we can't accept a connection, then finish loop and wait for next connection.
			log.Printf("error accepting connection: %s\n", err)
			continue
		}
		// got a connection--lets hand it off on it's own goroutine and get back to waiting for
		// the next potential connection and do it all again.
		log.Printf("ACCEPTED CONNECTION: %q\n", conn.RemoteAddr())
		go HandleTCPConn(conn)
	}
}

// HandleTCPConn is a server handler that takes a *net.TCPConn. A new one is
// created in a separate goroutine for every incoming connection that is accepted.
func HandleTCPConn(conn *net.TCPConn) {
	// make a buffered reader to read incoming data
	r := bufio.NewReader(conn)
	// make a buffered writer to write back to client
	w := bufio.NewWriter(conn)
	defer conn.Close()
	for {
		// read data until we get to the end of a line
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			// if we find the end of the data before
			// we find a newline char, then just stop
			break
		} else if err != nil {
			// if some other error happens, log and return
			log.Printf("Error reading: %s\n", err)
			HandleConnClose(conn)
			return
		}
		fmt.Printf("RECEIVED: %q FROM [%s], REPLYING...", data, conn.RemoteAddr())
		reply := fmt.Sprintf("ECHO: %q\n", data)
		if _, err = w.WriteString(reply); err != nil {
			log.Printf("Error writing: %s\n", err)
			fmt.Printf(" ERR!\n")
		} else {
			fmt.Printf(" OK!\n")
		}
	}
}

// Handler is a generic type definition
type Handler func(net.Conn)

// HandleEcho is a generic echo handler that conforms to the
// custom handler type we created. Just to show you that there
// are a TON of ways to implement different ways to handle
// connections and requests.
func HandleEcho() Handler {
	fn := func(conn net.Conn) {
		defer conn.Close()
		// echo all incoming data
		io.Copy(conn, conn)
	}
	return fn
}

// ListenAndServeTCPWithHandler creates a listening socket on the TCP
// protocol for the provided address and/port and handles TCP connections.
func ListenAndServeTCPWithHandler(host string, handle Handler) error {
	// resolve tcp address
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}
	// initialize a listening socket
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		// wait (block) for a connection, accepting it when it arrives.
		conn, err := ln.AcceptTCP()
		if err != nil {
			// if we can't accept a connection, then finish loop and wait for next connection.
			log.Printf("error accepting connection: %s\n", err)
			continue
		}
		// got a connection--lets hand it off on it's own goroutine and get back to waiting for
		// the next potential connection and do it all again.
		log.Printf("ACCEPTED CONNECTION: %q\n", conn.RemoteAddr())
		go handle(conn)
	}
}

// HandleConnClose is a semi generic handler that closes down a remote connection
// if something happens to the connected client and it is unrecoverable.
func HandleConnClose(conn net.Conn) {
	conn.Write([]byte("GOODBYE!\n"))
	log.Printf("CLOSED CONNECTION: %q\n", conn.RemoteAddr())
	conn.Close()
	conn = nil
	return
}
