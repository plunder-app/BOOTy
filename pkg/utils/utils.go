package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
)

//CmdlinePath is the default location for the cmdline
const CmdlinePath = "/proc/cmdline"

// ParseCmdLine will read through the command line and return the source and destination
func ParseCmdLine(path string) (src, dst string, err error) {

	// Read the file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Split by whitespace
	entries := strings.Fields(string(b))

	// find k=v entries
	for x := range entries {
		kv := strings.Split(entries[x], "=")
		if len(kv) == 2 {
			// find the entries we care about
			if kv[0] == "BOOTYSRC" {
				src = kv[1]
			}
			if kv[0] == "BOOTYDST" {
				dst = kv[1]
			}
		}
	}
	return
}

//ClearScreen will clear the screen of all text
func ClearScreen() {
	fmt.Print("\033[2J")
}
