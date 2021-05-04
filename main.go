package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
)

func main() {

	listener, err := net.Listen("tcp", ":5757")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	var connBuffer = &sync.Map{}

	fmt.Println("listening on :5757")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		id := uuid.New().String()
		connBuffer.Store(id, conn)
		go handleClient(id, conn, connBuffer)
	}

}

func handleClient(id string, conn net.Conn, buffer *sync.Map) {
	defer conn.Close()
	defer buffer.Delete(id)

	for {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		buffer.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if id != key.(string) {
					if _, err := conn.Write([]byte(input)); err != nil {
						fmt.Println(err)
					}
				}
			}
			return true
		})
	}
}
