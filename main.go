package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thebsdbox/BOOTy/pkg/plunderclient"

	"github.com/thebsdbox/BOOTy/pkg/realm"
	"github.com/thebsdbox/BOOTy/pkg/ux"
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

	sys := m.GetMount("sys")
	sys.CreateMount = true
	sys.EnableMount = true

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
	log.Infoln("Starting BOOTy")
	time.Sleep(time.Second * 2)
	ux.Captain()
	ux.SysInfo()

	log.Infoln("Beginning provisioning process")

	err := plunderclient.GetServerConfig()
	if err != nil {
		fmt.Println("The BOOTYURL=x.x.x.x is missing from the boot flags, rebooting in 10 seconds")
		time.Sleep(time.Second * 10)
		realm.Reboot()
	}

	// bs, _ := utils.GetBlockDeviceSize("sda")
	// fmt.Printf("/dev/sda is %d bytes\n", bs)

	// cmdline, _ := utils.ParseCmdLine("")
	// if _, ok := cmdline["PLNDRSVR"]; ok {
	// 	fmt.Printf("Server has been set in boot flags \n")
	// } else {
	// 	fmt.Println("The PLNDRSVR=x.x.x.x is missing from the boot flags, rebooting in 10 seconds")
	// 	time.Sleep(time.Second * 10)
	// 	//realm.Reboot()
	// }
	// TODO - remove

	realm.Shell()
	realm.Reboot()
}
