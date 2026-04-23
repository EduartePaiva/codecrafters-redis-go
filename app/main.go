package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

const (
	PONG = "+PONG\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		input, err := io.ReadAll(conn)
		if err != nil {
			fmt.Println("Error reading content: ", err.Error())
			continue
		}
		fmt.Println(string(input))

		_, err = conn.Write([]byte(PONG))
		if err != nil {
			fmt.Println("Error sending PONG response: ", err.Error())
			os.Exit(1)
		}

	}
}
