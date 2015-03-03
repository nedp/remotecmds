package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/nedp/remotecmds/router"
	"bitbucket.org/nedp/remotecmds/routes"
)

const defaultPort = 7832
const defaultNSlots = 16
const defaultMaxSlots = 32

func main() {
	println("Hello, world!")

	// Get the port to listen on.
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = defaultPort
	}

	// Get slot details
	nSlotsStr := os.Getenv("NSLOTS")
	nSlots, err := strconv.Atoi(nSlotsStr)
	if err != nil {
		nSlots = defaultNSlots
	}
	maxSlotsStr := os.Getenv("MAXSLOTS")
	maxSlots, err := strconv.Atoi(maxSlotsStr)
	if err != nil {
		maxSlots = defaultMaxSlots
	}

	// Listen on the specified port.
	laddr, err := net.ResolveTCPAddr("tcp6", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Couldn't resolve local TCP address: %s", err.Error())
	}
	listener, err := net.ListenTCP(laddr.Network(), laddr)
	defer listener.Close()
	if err != nil {
		log.Fatalf("Couldn't listen on local TCP address: %s", err.Error())
	}

	// Make a CommandRouter and specify the routes.
	cmdr := router.New(nSlots, maxSlots)
	routes.AddRoutesTo(cmdr)

	// Accept all connections.
	for true {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("error accepting connection: %s", err.Error())
			continue
		}
		handle(cmdr, conn)
	}
}

// handle reads, routes, and responds to the request,
// then closes the connection.
func handle(cmdr router.Interface, conn net.Conn) {
	const bufferSize = 1024
	var lines []string

	// Read the request
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading from connection: %s", err.Error())
		conn.Close()
		return
	}

	// Route the request
	req := strings.Join(lines, "\n")
	out, err := cmdr.OutputFor(req)
	if err != nil {
		log.Printf("error routing request: %s", err.Error())
		conn.Close()
		fmt.Fprintf(conn, "ERROR: couldn't route the request.\n")
		return
	}

	// Respond to the request.
	// Put off the creation of a new goroutine until everything's
	// verified to save resources.
	go respondAndClose(out, conn)
}

func respondAndClose(out <-chan string, conn net.Conn) {
	for s := range out {
		n, err := fmt.Fprintf(conn, "%s\n", s)
		if err != nil {
			log.Printf("WARN: main: failed to write to connection: %s", err.Error())
			continue
		}
		if n < len([]byte(s)) {
			log.Printf("WARN: response incomplete - sent %d of %d bytes",
				n, len([]byte(s)))
			continue
		}
	}
	conn.Close()
}
