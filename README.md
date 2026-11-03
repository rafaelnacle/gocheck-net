# gocheck-net 🌐

An interactive, concurrent terminal network status checker written in Go with **Bubble Tea** and **Lip Gloss**.

```
  GOCHECK-NET   Interactive Network Status Checker

 TARGET                             PROTOCOL   STATUS       CODE     LATENCY      DETAILS
 https://google.com                 HTTP          UP        200      45ms         OK
 https://github.com                 HTTP          UP        200      110ms        OK
 1.1.1.1:53                         TCP           UP        -        15ms         OK
 https://httpbin.org/status/500     HTTP         DOWN       500      85ms         OK
 localhost:59999                    TCP          DOWN       -        0s           dial tcp: connection refused

 ╭────────────────────────────────────────────────╮
 │ Total: 5  •  Up: 3  •  Down: 2  •  In-flight: 0│
 ╰────────────────────────────────────────────────╯

  [r] Refresh  [q/esc] Quit
```

## Features

- **Protocol Auto-detection**: Automatically handles HTTP/HTTPS URLs and TCP `host:port` pairs.
- **Concurrent Probing**: Targets are probed asynchronously in parallel using Goroutines and Bubble Tea commands (`tea.Cmd`).
- **Interactive UI**: Live terminal table rendered with Lip Gloss, color-coded latency, status badges, and real-time statistics.
- **Clean Architecture**: Decoupled core network engine (`internal/checker`) and TUI presentation layer (`internal/ui`).
- **Standard Library Networking**: Built with pure Go `net/http` and `net` without bloated network dependencies.

## Usage

### Run with default targets

```bash
go run main.go
```

### Check custom endpoints

```bash
go run main.go https://api.github.com google.com 1.1.1.1:53 localhost:8080
```

### Custom timeout

```bash
go run main.go -timeout 1.5s https://example.com 8.8.8.8:53
```

## Hotkeys

- `r`: Refresh / re-check all targets immediately
- `q` / `Esc` / `Ctrl+C`: Exit program

## Running Tests

Run table-driven unit tests with race detection:

```bash
go test -v -race ./...
```
