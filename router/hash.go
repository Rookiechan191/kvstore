package main

import "hash/fnv"

func hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func getNode(key string, nodes []string) string {
	hash := hashKey(key)
	index := int(hash % uint32(len(nodes)))
	return nodes[index]
}