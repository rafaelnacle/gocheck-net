package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"gocheck-net/internal/checker"
	"gocheck-net/internal/ui"
)

var defaultTargets = []string{
	"https://google.com",
	"https://github.com",
	"https://cloudflare.com",
	"1.1.1.1:53",               // Cloudflare DNS TCP
	"8.8.8.8:53",               // Google DNS TCP
	"https://httpbin.org/status/200",
	"https://httpbin.org/status/500", // Will show DOWN (server error)
	"localhost:59999",                 // Unopened local port, will show DOWN
}

func main() {
	timeout := flag.Duration("timeout", 3*time.Second, "Timeout duration for each probe (e.g. 2s, 500ms)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gocheck-net [flags] [targets...]\n\n")
		fmt.Fprintf(os.Stderr, "Targets can be HTTP/HTTPS URLs or TCP host:port addresses.\n")
		fmt.Fprintf(os.Stderr, "If no targets are specified, a default sample set is checked.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	rawTargets := flag.Args()
	if len(rawTargets) == 0 {
		rawTargets = defaultTargets
	}

	targets := make([]checker.Target, len(rawTargets))
	for i, raw := range rawTargets {
		targets[i] = checker.ParseTarget(raw)
	}

	model := ui.New(targets, *timeout)
	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running gocheck-net: %v\n", err)
		os.Exit(1)
	}
}
