package tcpmessage

import (
	"fmt"
	"log"
	"net"
)

var conn net.Conn

func init() {
	var err error
	ip := "localhost"
	// ip := "192.168.185.77"
	port := "8080"
	conn, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	//	defer conn.Close()
}

func Send(message string) {
	fmt.Fprintf(conn, message+"\n")
}
