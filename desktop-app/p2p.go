package main

import (
	"bufio"
	"fmt"
	"log"
	"netdrop/desktop-app/netutils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
)

func main() {
	// args := os.Args[1:]

	// if len(args) < 2 {
	// 	log.Fatalf("Usage: go run p2p.go <listening port> <name of machine>")
	// }

	port := "8000"

	name, err := os.Hostname()
	// name := args[1]

	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	service, err := mdns.NewMDNSService(
		name,
		"_custom._tcp",
		"",
		"",
		parsePort(port),
		nil,
		[]string{"registered-via-mdns"},
	)
	if err != nil {
		log.Fatalf("Failed to create mDNS service: %v", err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		log.Fatalf("Failed to start mDNS server: %v", err)
	}
	defer server.Shutdown()

	log.Printf("Registered '%s' on port %s", name, port)

	go netutils.StartListening(port)

	for true {
		reader := bufio.NewReader(os.Stdin)
		val, _ := reader.ReadString('\n')
		targetName := strings.Split(val, " ")[0]
		fileName := strings.TrimSpace(strings.Split(val, " ")[1])

		entriesCh := make(chan *mdns.ServiceEntry, 4)

		found := false
		go func() {
			for entry := range entriesCh {
				name := strings.Split(entry.Name, ".")[0]

				if name == targetName {
					fmt.Printf("Found IP for %s: %s:%d\n", targetName, entry.AddrV4, entry.Port)
					netutils.SendMessage(entry.AddrV4.String(), strconv.Itoa(entry.Port), fileName)
					found = true
				}
			}
		}()

		mdns.Lookup("_custom._tcp", entriesCh)
		close(entriesCh)
		time.Sleep(2 * time.Second)

		if !found {
			fmt.Println("No matching service found.")
		}
	}

}

func parsePort(p string) int {
	var port int
	_, err := fmt.Sscanf(p, "%d", &port)
	if err != nil {
		log.Fatalf("Invalid port number: %s", p)
	}
	return port
}
