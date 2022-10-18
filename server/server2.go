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
	//fmt.Println(err)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
		}
		conns <- conn
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
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
		}
		msgs <- Message{clientid, msg}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above. DONE

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		handleError(err)
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	newestID := 0
	//shutDown := false
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection DONE
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
			clients[newestID] = conn
			go handleClient(conn, newestID, msgs)
			newestID++

		case msg := <-msgs:
			//TODO Deal with a new message DONE
			// Send the message to all clients that aren't the sender
			senderID := msg.sender
			for i, _ := range clients {
				if i != senderID {
					//fmt.Println("PRINTING TO ID ", i)
					fmt.Fprintf(clients[i], msg.message)
				}
			}
			if msg.message == "shutdown" {
				break
			}

		}
	}
}
