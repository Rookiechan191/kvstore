# kvstore

A simple in-memory key-value store server written in Go. It listens on a TCP port and accepts text-based commands to set, get, and delete key-value pairs.

## Getting Started

### Prerequisites

- Go 1.21 or later

### Build

```bash
go build -o kvstore .
```

### Run

```bash
./kvstore
```

The server starts listening on port **8080**.

## Usage

Connect to the server using any TCP client, such as `nc` (netcat):

```bash
nc localhost 8080
```

### Commands

| Command | Syntax | Description |
|---------|--------|-------------|
| `SET` | `SET <key> <value>` | Store a key-value pair |
| `GET` | `GET <key>` | Retrieve the value for a key |
| `DELETE` | `DELETE <key>` | Remove a key-value pair |

### Examples

```
SET name Alice
OK

GET name
Alice

GET unknown
NULL

DELETE name
OK

GET name
NULL
```

## Project Structure

```
kvstore/
├── main.go    # TCP server and command handling
├── store.go   # In-memory KVStore implementation
└── go.mod     # Go module definition
```
