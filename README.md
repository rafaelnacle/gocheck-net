# gocheck-net 🌐

An interactive, visual terminal network status playground and endpoint checker written in Go with **Bubble Tea** and **Lip Gloss**.

```text
  GOCHECK-NET   Interactive Network Playground
  ( ^‿^ ) All systems operational! 100% Health         [██████████████] 100%

   TARGET                       PROTO    STATUS       SPEED GAUGE               DETAILS
 ❯ https://google.com           HTTP        UP        [██████████] ⚡ 45ms        HTTP 200 OK
   https://github.com           HTTP        UP        [████████░░] 🚀 92ms        HTTP 200 OK
   1.1.1.1:53                   TCP         UP        [██████████] ⚡ 18ms        OK
   https://httpbin.org/status/2 HTTP        UP        [██████░░░░] ✨ 210ms       HTTP 200 OK
   localhost:59999              TCP        DOWN       [░░░░░░░░░░] 💥 failed       dial tcp: connection refused

 ╭─────────────────────────────────────────────────────────────────────────────╮
 │ Total: 5  •  Up: 4  •  Down: 1  •  In-flight: 0  •  Health: 80%             │
 ╰─────────────────────────────────────────────────────────────────────────────╯

  [↑/↓] Select  [+] Add  [x] Delete  [t] Theme (Default Modern)  [s] Bell (ON 🔔)  [r] Refresh  [q] Quit
```

---

## ✨ Features

- **Expressive Network Mascot & Health Meter**: ASCII companion header that reacts to your network state (`(^‿^)`, `(•_•)`, `(>_<)`) alongside a real-time health gauge.
- **Visual Speed Bars & Latency Ratings**: Colored graphical bars with speed tier emojis:
  - `⚡` Blazing (<80ms)
  - `🚀` Fast (<150ms)
  - `✨` Good (<300ms)
  - `🚶` Moderate (<600ms)
  - `🐌` Slow (>=600ms)
  - `💥` Failed / Down
- **Live Target Management**:
  - Press `+` or `a` to open an inline text input and add any website or `host:port` target on the fly.
  - Press `x` or `d` to delete the selected target.
- **4 Color Themes**:
  - Press `t` to cycle through palettes: **Default Modern**, **Matrix Hacker (CRT Green)**, **Cyberpunk Neon (Hot Pink & Cyan)**, and **Retro Amber (1980s Vintage)**.
- **Audio Outage Alerts**:
  - Terminal audio bell (`\a`) rings when an outage is detected on a probed target (toggle on/off with `s`).
- **Concurrent Asynchronous Probes**: Targets are probed in parallel without freezing the terminal using Bubble Tea commands (`tea.Cmd`).
- **Standard Library Networking**: Pure Go `net/http` and `net` without bloated network dependencies.

---

## 🚀 Usage

### Run with default targets
```bash
go run main.go
```

### Check custom endpoints
```bash
go run main.go https://api.github.com google.com 1.1.1.1:53 localhost:8080
```

### Custom timeout duration
```bash
go run main.go -timeout 1.5s https://example.com 8.8.8.8:53
```

---

## ⌨️ Hotkeys

| Hotkey | Description |
| :--- | :--- |
| `↑` / `k` | Move selection cursor up |
| `↓` / `j` | Move selection cursor down |
| `+` / `a` | Add new target interactively (URL or `host:port`) |
| `x` / `d` | Delete selected target |
| `t` | Cycle through color themes (*Default, Matrix, Cyberpunk, Amber*) |
| `s` | Toggle audio bell alerts for outages (*ON / OFF*) |
| `r` | Refresh and re-probe all targets immediately |
| `Esc` | Cancel input prompt or clear notices |
| `q` / `Ctrl+C` | Quit program |

---

## 🧪 Running Tests

Run table-driven unit tests with race condition detection:
```bash
go test -v -race ./...
```
