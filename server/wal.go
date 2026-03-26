package main

import (
	"bufio"
	"os"
)

type WAL struct {
	file *os.File
}

func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: file}, nil
}

func (w *WAL) Write(command string) error {
	_, err := w.file.WriteString(command + "\n")
	if err != nil {
		return err
	}
	return w.file.Sync()
}

func (w *WAL) Load(store *KVStore) error {
	_, err := w.file.Seek(0, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(w.file)
	for scanner.Scan() {
		line := scanner.Text()
		ApplyCommand(store, line)
	}
	return scanner.Err()
}