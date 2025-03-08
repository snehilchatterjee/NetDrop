package main

import (
	"netdrop/desktop-app/netutils"
	"os"
)

func main() {
	args := os.Args[1:]
	netutils.StartListening(args[0])
}
