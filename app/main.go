package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

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

		go func() {
			defer conn.Close()
			buf := make([]byte, 1024)
			for {
				_, err := conn.Read(buf)

				fmt.Println("did it read?")
				if err != nil {
					fmt.Println("Error reading content: ", err.Error())
					break
				}
				fmt.Println(string(buf))

				_, rsp := resp.ReadNextRESP(buf)

				response, err := DispatchCommand(rsp)
				if err != nil {
					fmt.Println("Error processing command: ", err.Error())
					break
				}

				_, err = conn.Write(response)
				if err != nil {
					fmt.Println("Error sending response: ", err.Error())
					break
				}
			}
		}()

	}

}
