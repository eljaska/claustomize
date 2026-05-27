package main

import (
	"fmt"
	"os"

	"github.com/eljaska/claustomize/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
