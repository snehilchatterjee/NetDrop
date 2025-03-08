package netutils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func StartListening(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server listening on port:", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read file extension size
	var extSize int64
	err := binary.Read(conn, binary.LittleEndian, &extSize)
	if err != nil {
		fmt.Println("Error reading extension size:", err)
		return
	}

	// Read file extension
	extBytes := make([]byte, extSize)
	_, err = io.ReadFull(conn, extBytes)
	if err != nil {
		fmt.Println("Error reading file extension:", err)
		return
	}
	fileExt := string(extBytes)

	// Read file size
	var fileSize int64
	err = binary.Read(conn, binary.LittleEndian, &fileSize)
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return
	}

	// Create file with preserved extension
	fileName := fmt.Sprintf("received_file.%s", fileExt)
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	// Receive file data
	written, err := io.CopyN(outFile, conn, fileSize)
	if err != nil {
		fmt.Println("Error receiving file:", err)
		return
	}

	fmt.Printf("File received successfully: %s (%d bytes)\n", fileName, written)
}

func SendMessage(port string, filename string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Get file details
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	// Extract file extension
	fileNameSplit := strings.Split(filename, ".")
	fileExt := fileNameSplit[len(fileNameSplit)-1]

	// Send extension size
	extSize := int64(len(fileExt))
	err = binary.Write(conn, binary.LittleEndian, extSize)
	if err != nil {
		fmt.Println("Error sending extension size:", err)
		return
	}

	// Send extension
	_, err = conn.Write([]byte(fileExt))
	if err != nil {
		fmt.Println("Error sending file extension:", err)
		return
	}

	// Send file size
	err = binary.Write(conn, binary.LittleEndian, fileSize)
	if err != nil {
		fmt.Println("Error sending file size:", err)
		return
	}

	// Send file content
	written, err := io.CopyN(conn, file, fileSize)
	if err != nil {
		fmt.Println("Error sending file:", err)
		return
	}

	fmt.Printf("File sent successfully (%d bytes)\n", written)
}
