package app

import (
	"runtime"
	"strings"
)

type Command struct {
	usage    string
	desc     string
	commands []*Command
	runFunc  RunCommandFunc
}

// RunCommandFunc defines the application's command startup callback function.
type RunCommandFunc func(args []string) error

func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
