package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	port  uint
	check bool
)

func init() {
	flag.UintVar(&port, "port", 0, "port to listen on")
	flag.BoolVar(&check, "check", false, "perform a health check and exit")
	flag.Parse()

	if check {
		log.Println("Healthcheck passed!")
		os.Exit(0)
	}
}

func main() {
	var (
		listener net.Listener
		err      error
	)

	listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	defer listener.Close()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	addr := c.RemoteAddr().String()

	log.Printf("==> Serving %s\n", addr)

	for scanner := bufio.NewScanner(c); scanner.Scan(); {
		// fmt.Printf("%s: %s\n", addr, scanner.Text())
		c.Write([]byte(scanner.Text() + "\n"))
		if err := scanner.Err(); err != nil {
			log.Println("reading standard input:", err)
		}
	}

	log.Printf("<== DONE (%s)\n", addr)
	c.Close()
}
