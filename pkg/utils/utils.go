package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//CmdlinePath is the default location for the cmdline
const CmdlinePath = "/proc/cmdline"

// ParseCmdLine will read through the command line and return the source and destination
func ParseCmdLine(path string) (m map[string]string, err error) {
	// allow path override
	if path == "" {
		path = CmdlinePath
	}

	m = make(map[string]string)
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
			m[kv[0]] = kv[1]
		}
	}
	return
}

//ClearScreen will clear the screen of all text
func ClearScreen() {
	fmt.Print("\033[2J")
}

// GetBlockDeviceSize will read the size from the /sys/block for a specific block device
func GetBlockDeviceSize(device string) (int64, error) {

	// This should return the path to the block device and it's size (in sectores)
	// Each sector is 512 bytes

	path := fmt.Sprintf("/sys/block/%s/size", device)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	parsedData := strings.TrimSpace(string(data))
	size, _ := strconv.ParseInt(parsedData, 10, 64)
	return size * 512, nil
}
