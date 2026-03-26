package main

import "strings"

func ApplyCommand(store *KVStore, line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "SET":
		if len(parts) >= 3 {
			store.Set(parts[1], parts[2])
		}
	case "DELETE":
		if len(parts) >= 2 {
			store.Delete(parts[1])
		}
	}
}
