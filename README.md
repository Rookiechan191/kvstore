# kvstore

A simple TCP key-value store in Go with:

- In-memory storage (`map[string]string` protected by `sync.RWMutex`)
- Text protocol over TCP on port `8080`
- Write-ahead log (WAL) persistence to `wal.log`
- Recovery on startup by replaying WAL commands

## Requirements

- Go version compatible with [go.mod](go.mod)

## Run

Because `benchmark.go` also defines a `main` function, run the server using explicit files:

```bash
go run main.go store.go wal.go command.go
```

Or build the server binary from those same files:

```bash
go build -o kvstore main.go store.go wal.go command.go
```

Server output:

```text
KV store listening on port 8080
```

## Protocol

Connect with a TCP client (for example, `nc`):

```bash
nc localhost 8080
```

Supported commands:

| Command | Syntax | Response |
|---|---|---|
| `SET` | `SET <key> <value>` | `OK` |
| `GET` | `GET <key>` | `<value>` or `NULL` |
| `DELETE` | `DELETE <key>` | `OK` |

Error responses:

- `ERROR: SET needs key and value`
- `ERROR: GET needs a key`
- `ERROR: DELETE needs a key`
- `ERROR: unknown command`

Example session:

```text
SET name Alice
OK
GET name
Alice
DELETE name
OK
GET name
NULL
```

## Persistence (WAL)

- `SET` and `DELETE` are appended to `wal.log` before being applied in-memory.
- On startup, `wal.log` is replayed to rebuild in-memory state.
- `GET` is not logged.

## Benchmark

`benchmark.go` is a simple load generator that sends `10000` `SET` requests over TCP.

Start server first, then run:

```bash
go run benchmark.go
```

## Known Limitations

- Values are parsed with `strings.Fields`, so spaces are not preserved (only the third token is used as value).
- No graceful shutdown hook to close WAL explicitly.
- WAL file grows indefinitely (no compaction/snapshotting).
- Protocol is line-based and minimal (no auth, no TTL, no transactions).

## Project Files

- `main.go`: TCP server, command dispatch, startup/recovery flow
- `store.go`: thread-safe in-memory key-value store
- `wal.go`: WAL append and replay logic
- `command.go`: WAL command application (`SET`, `DELETE`)
- `benchmark.go`: standalone benchmark client
