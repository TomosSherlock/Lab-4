package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	fmt.Println(err)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
		} else {
			conns <- conn
		}
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')

		msgs <- Message{clientid, msg}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)

	if err != nil {
		handleError(err)
	} else {
		//Create a channel for connections
		conns := make(chan net.Conn)
		//Create a channel for messages
		msgs := make(chan Message)
		//Create a mapping of IDs to connections
		clients := make(map[int]net.Conn)

		//Start accepting connections
		go acceptConns(ln, conns)
		for {
			select {
			case conn := <-conns:
				//TODO Deal with a new connection
				// - assign a client ID
				// - add the client to the clients channel
				// - start to asynchronously handle messages from this client
				clientID := len(clients)
				clients[clientID] = conn
				go handleClient(conn, clientID, msgs)
			case msg := <-msgs:
				//TODO Deal with a new message
				// Send the message to all clients that aren't the sender
				for _, conn := range clients {
					//if clients[i] != conn {
					fmt.Fprintln(conn, msg)
					//}
				}
			}
		}
	}
}
