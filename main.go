package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	var port int
	var listener net.Listener
	var err error

	if str := os.Getenv("PORT"); len(str) > 0 {
		port, err = strconv.Atoi(str)
		if err != nil {
			log.Fatalln("PORT must be a valid integer", err)
		}
	} else {
		log.Fatalln("PORT must not be empty")
	}

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
	fmt.Printf("Serving %s ==> \n", addr)
	// for {
	// 	netData, err := bufio.NewReader(c).ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	temp := strings.TrimSpace(string(netData))
	// 	if temp == "STOP" {
	// 		break
	// 	}

	// 	result := strconv.Itoa(rand.Int()) + "\n"
	// 	c.Write([]byte(string(result)))
	// }
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		// fmt.Printf("%s: %s\n", addr, scanner.Text())
		c.Write([]byte(scanner.Text() + "\n"))
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
	fmt.Printf("<== DONE (%s)\n", addr)
	c.Close()
}
