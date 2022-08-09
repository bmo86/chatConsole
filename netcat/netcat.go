package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("p", 3090, "Port")
	host = flag.String("h", "localhost", "host")
)

/*
	-> host:port
	Writer -> host:port
	Reader -> host:port
	> [hola] -> host:port -> 	[hola]

*/

func main() {

	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn) // recibe dos interfaces 1 de lectura y 1 de escritura
		done <- struct{}{}
	}()

	CopyContent(conn, os.Stdin) //actua como escritor
	conn.Close()
	<-done
}

func CopyContent(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}
