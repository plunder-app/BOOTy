package realm

import (
"os"
"syscall"

log "github.com/sirupsen/logrus"
)

// This contains all methods for managing the final steps with a host

// Reboot a host
func Reboot() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		log.Printf("reboot off failed: %v", err)
		Shell()
	}
	// Should cause a panic
	os.Exit(1)
}

// PowerOff will result in the host using an ACPI power off
func PowerOff() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		log.Printf("power off failed: %v", err)
		Shell()
	}
	// Should cause a panic
	os.Exit(1)
}

// Halt will instuct the CPU to enter a halt state (no-power off (usually))
func Halt() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
	if err != nil {
		log.Printf("halt failed: %v", err)
		Shell()
	}
	// Should cause a panic
	os.Exit(1)
}
