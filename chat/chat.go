package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Client chan<- string

var (
	inComingClients = make(chan Client) //canal de canales
	leavingClients  = make(chan Client) //canal de canales
	msg             = make(chan string)
)

var (
	port = flag.Int("p", 3090, "Port")
	host = flag.String("h", "localhost", "host")
)

func handlerConnection(conn net.Conn) {
	defer conn.Close()
	msgs := make(chan string)
	go MsgWriter(conn, msgs)

	clientName := conn.RemoteAddr().String()

	msgs <- fmt.Sprintf("Welcome to the server, your Name: %s\n", clientName)
	msg <- fmt.Sprintf("New Client is here, name: %s\n", clientName)
	inComingClients <- msgs

	inputMsg := bufio.NewScanner(conn)
	for inputMsg.Scan() {
		msg <- fmt.Sprintf("%s => %s\n", clientName, inputMsg.Text())
	}

	leavingClients <- msg
	msg <- fmt.Sprintf("%s said goodbye", clientName)
}

func MsgWriter(conn net.Conn, message <-chan string) {
	for msg := range message {
		fmt.Fprintln(conn, msg)
	}
}

func BroadCast() {
	clients := make(map[Client]bool)
	for {
		select {
		case msgs := <-msg:
			for client := range clients {
				client <- msgs
			}
		case newClient := <-inComingClients:
			clients[newClient] = true

		case leavingClient := <-leavingClients:
			delete(clients, leavingClient)
			close(leavingClient) //cerrar el canal
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	go BroadCast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handlerConnection(conn)
	}

}
