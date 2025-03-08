package netutils

import (
	"fmt"
	"net"
)

func StartListening(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	fmt.Println("Server Listening on port:", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Received: %s \n", buff)
}

func SendMessage(port string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	fmt.Println("Sending data to port", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	return
}
