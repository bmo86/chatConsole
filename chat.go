package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client chan<- string

var (
	inComingClients = make(chan Client) //canal de canales
	leavingClients  = make(chan Client) //canal de canales
	msg             = make(chan string)
)

/* var (
	host = flag.String("h", "localhost", "Host")
	port = flag.Int("p", 3090, "Port")
)
*/
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
	msg <- fmt.Sprintf("%s soid goodbye", clientName)
}

func MsgWriter(conn net.Conn, message <-chan string) {
	for msg := range message {
		fmt.Fprintln(conn, msg)
	}
}
