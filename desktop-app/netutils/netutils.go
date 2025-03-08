package netutils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
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

	var size int64
	binary.Read(conn, binary.LittleEndian, &size)

	outFile, err := os.Create("received_file")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	written, err := io.CopyN(outFile, conn, size)
	if err != nil {
		fmt.Println("Error receiving file:", err)
		return
	}

	fmt.Printf("File received successfully (%d bytes)\n", written)
}

func SendMessage(port string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	fmt.Println("Sending data to port", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open("sample_file")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	err = binary.Write(conn, binary.LittleEndian, fileSize)
	if err != nil {
		fmt.Println("Error sending file size:", err)
		return
	}

	written, err := io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file:", err)
		return
	}

	fmt.Printf("File sent successfully (%d bytes)\n", written)

	defer conn.Close()
	return
}
