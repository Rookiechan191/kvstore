# kvstore

A distributed, in-memory key-value store written in Go. Clients talk plain TCP text, writes are durably logged to a WAL before being applied, and a separate router process uses FNV-32a hashing to shard keys across multiple nodes. Each node can asynchronously replicate writes to a single peer.

---

## Folder structure

```
kvstore/
├── benchmark/
│   └── benchmark.go       # Standalone load tester
├── router/
│   ├── hash.go            # FNV-32a hashing + node selection
│   └── router.go          # Interactive stdin router
├── server/
│   ├── command.go         # ApplyCommand: parses SET / DELETE
│   ├── main.go            # TCP server, connection handler, replication
│   ├── store.go           # KVStore: RWMutex-protected map
│   ├── wal.go             # Write-Ahead Log: append, fsync, replay
│   ├── wal_8080.log       # Generated at runtime
│   ├── wal_8081.log       # Generated at runtime
│   └── wal_8082.log       # Generated at runtime
├── .gitignore
├── go.mod
└── README.md
```

Each subdirectory is its own `package main` — they are separate binaries and cannot be compiled together.

---

## How it works

### Server (`server/`)

Accepts TCP connections and handles commands line-by-line using `bufio.Scanner`. Each connection runs in its own goroutine. For every mutating command (`SET` / `DELETE`), the server:

1. Writes the raw command line to the WAL and calls `fsync`
2. Applies the mutation to the in-memory store
3. Spawns a goroutine to forward the command to the configured replica

`REPL` commands (incoming replication messages) bypass the WAL entirely — they are applied directly via `ApplyCommand` with no response written back. This prevents replication loops and avoids double-logging.

On startup, the WAL file is seeked to offset 0 and replayed line-by-line to restore state. The file is opened with `O_CREATE|O_APPEND|O_RDWR` — the `O_RDWR` flag is required so the startup seek works; `O_APPEND` alone would be write-only.

### Store (`server/store.go`)

A plain `map[string]string` protected by a `sync.RWMutex`. `Get` takes a read lock; `Set` and `Delete` take a write lock.

### Router (`router/`)

Reads a command from stdin, extracts the key with `fmt.Sscanf`, runs it through FNV-32a (`hash/fnv`), and picks a node via `int(hash) % len(nodes)`. The full raw input line is forwarded over a new TCP connection and the first response line is printed. Nodes are hardcoded as `localhost:8080`, `8081`, `8082`.

### Benchmark (`benchmark/`)

Opens a new TCP connection for each of 10,000 `SET` commands, writes the command, then closes the connection immediately without reading the `OK` response back. The reported RPS measures connection-setup + write latency, not round-trip latency.

---

## Getting started

### Prerequisites

- Go 1.21+

### Run the server

```bash
cd server
go run . 8080 8081
```

Starts a TCP server on `:8080` that replicates writes to `localhost:8081`.

### Run a replica

```bash
cd server
go run . 8081 8080
```

### Use the interactive router

```bash
cd router
go run .
```

```
> SET foo bar
Routing to: localhost:8081
OK
> GET foo
Routing to: localhost:8081
bar
```

### Run the benchmark

```bash
cd benchmark
go run .
```

Fires 10,000 `SET key<n> value<n>` commands sequentially against `localhost:8080`.

---

## Commands

| Command | Syntax | Response |
|---|---|---|
| `SET` | `SET <key> <value>` | `OK` |
| `GET` | `GET <key>` | `<value>` or `NULL` |
| `DELETE` | `DELETE <key>` | `OK` |

All messages are newline-delimited plain text. Values cannot contain spaces — commands are parsed with `strings.Fields`, so `SET key hello world` silently drops `world`.

---

## Known issues and limitations

**Correctness**

- **`int(hash) % len(nodes)` sign bug** — `hashKey` returns `uint32`. Casting to `int` can go negative on 32-bit systems, causing a runtime panic. Fix: `int(hash % uint32(len(nodes)))`.
- **No read-your-writes guarantee** — replication is fire-and-forget in a goroutine; a `GET` issued immediately after a `SET` may return a stale value from the replica.
- **Values cannot contain spaces** — `strings.Fields` splits on all whitespace, so multi-word values are silently truncated.
- **WAL is never compacted** — the log grows indefinitely. Replay time grows with it; there is no snapshotting or truncation.
- **No replication acknowledgment** — the primary does not verify that the replica applied the write.

**Router**

- **`GET` may route to the wrong node** — if a key was written via a direct connection (bypassing the router), it may live on a different node than where the router sends the `GET`.
- **Nodes are hardcoded** — `localhost:8080`, `8081`, `8082` are hardcoded in `router.go`; there is no config file or flag.

**Benchmark**

- **Does not read responses** — the connection closes before the server `OK` arrives, which can produce write errors on the server side.
- **Sequential, single-threaded** — one connection at a time; does not measure concurrent throughput.

---

## License

MIT
