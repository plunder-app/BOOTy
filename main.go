package main

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/thebsdbox/BOOTy/pkg/realm"
)

func main() {

	// Fuck it

	//cmd.Execute()
	m := realm.DefaultMounts()
	d := realm.DefaultDevices()
	dev := m.GetMount("dev")
	dev.CreateMount = true
	dev.EnableMount = true

	proc := m.GetMount("proc")
	proc.CreateMount = true
	proc.EnableMount = true

	tmp := m.GetMount("tmp")
	tmp.CreateMount = true
	tmp.EnableMount = true

	// Create all folders
	m.CreateFolder()
	// Ensure that /dev is mounted (first)
	m.CreateNamedMount("dev", true)

	// Create all devices
	d.CreateDevice()

	// Mount any additional mounts
	m.CreateMount()

	log.Println("Starting DHCP client")
	go realm.DHCPClient()

	// HERE IS WHERE THE MAIN CODE GOES

	// Shell stuff
	log.Println("Starting Shell")
	cmd := exec.Command("/bin/sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Errorf("Shell error [%v]", err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		log.Errorf("Shell error [%v]", err)
	}
	log.Infoln()
	exit()
}

func exit() {
	var err error

	err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		log.Printf("reboot off failed: %v", err)
	}

	// err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	// if err != nil {
	// 	log.Printf("power off failed: %v", err)
	// }

	// err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
	// if err != nil {
	// 	log.Printf("halt failed: %v", err)
	// }

	os.Exit(1)
	select {}
}
