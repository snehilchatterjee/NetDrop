package netutils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

	buff := new(bytes.Buffer)
	var size int64
	binary.Read(conn, binary.LittleEndian, &size)

	for {
		_, err := io.CopyN(buff, conn, size)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received: %s \n", buff.Bytes())
		buff.Reset()
	}
}

func SendMessage(port string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	fmt.Println("Sending data to port", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	file := []byte("Hello, server!")
	// file := make([]byte, 4096)
	size := len(file)
	// size := 4096

	binary.Write(conn, binary.LittleEndian, int64(size))
	_, err = io.CopyN(conn, bytes.NewReader(file), int64(size))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	return
}
