package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var nodes = []string{
	"localhost:8080",
	"localhost:8081",
	"localhost:8082",
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		var cmd, key string
		fmt.Sscanf(text, "%s %s", &cmd, &key)

		node := getNode(key, nodes)
		fmt.Println("Routing to:", node)

		conn, err := net.Dial("tcp", node)
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		conn.Write([]byte(text))

		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(response)

		conn.Close()
	}
}