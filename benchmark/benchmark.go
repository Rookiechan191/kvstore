package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	start := time.Now()

	total := 10000

	for i := 0; i < total; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Connection error:", err)
			return
		}

		command := fmt.Sprintf("SET key%d value%d\n", i, i)
		conn.Write([]byte(command))

		_, _ = bufio.NewReader(conn).ReadString('\n')

		conn.Close()
	}

	elapsed := time.Since(start)
	fmt.Println("Total requests:", total)
	fmt.Println("Time taken:", elapsed)
	fmt.Println("Requests per second:", float64(total)/elapsed.Seconds())
}