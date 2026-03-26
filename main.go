package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, store *KVStore, wal *WAL) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("ERROR: SET needs key and value\n"))
				continue
			}

			wal.Write(line)

			store.Set(parts[1], parts[2])
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERROR: GET needs a key\n"))
				continue
			}
			value, exists := store.Get(parts[1])
			if exists {
				conn.Write([]byte(value + "\n"))
			} else {
				conn.Write([]byte("NULL\n"))
			}

		case "DELETE":
			if len(parts) < 2 {
				conn.Write([]byte("ERROR: DELETE needs a key\n"))
				continue
			}

			wal.Write(line)

			store.Delete(parts[1])
			conn.Write([]byte("OK\n"))

		default:
			conn.Write([]byte("ERROR: unknown command\n"))
		}
	}
}

func main() {
	store := NewKVStore()

	wal, err := NewWAL("wal.log")
	if err != nil {
		fmt.Println("WAL error:", err)
		return
	}

	err = wal.Load(store)
	if err != nil {
		fmt.Println("WAL load error:", err)
		return
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("failed to start server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("KV store listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("connection error:", err)
			continue
		}

		go handleConnection(conn, store, wal)
	}
}