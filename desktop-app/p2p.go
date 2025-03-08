package main

import (
	"bufio"
	"netdrop/desktop-app/netutils"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	go netutils.StartListening(args[0])

	for true {
		reader := bufio.NewReader(os.Stdin)
		val, _ := reader.ReadString('\n')
		val = strings.TrimSpace(val)

		if val == "send" {
			netutils.SendMessage(args[1])
		}
	}

}
