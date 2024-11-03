package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces)")
var port = flag.Int("port", 8080, "The port to listen on; default is 8080.")

func main() {
	flag.Parse()

	fmt.Print("Starting Server... \n")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected from %s\n", remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		go handleMessage(scanner.Text(), conn)
	}
	fmt.Printf("client at %s disconnected", remoteAddr)
}

func handleMessage(message string, conn net.Conn) {
	fmt.Print("> \n" + message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/time":
			resp := "It is " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))
		case message == "/quit":
			fmt.Print("Quitting.\n")
			conn.Write([]byte("I'm shutting down now. \n"))
			fmt.Print("< \n" + "%quit%")
			conn.Write([]byte("%quit% \n"))
			os.Exit(0)
		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}
