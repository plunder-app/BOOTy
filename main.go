package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thebsdbox/BOOTy/pkg/image"
	"github.com/thebsdbox/BOOTy/pkg/plunderclient"
	"github.com/thebsdbox/BOOTy/pkg/plunderclient/types"
	"github.com/thebsdbox/BOOTy/pkg/utils"

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

	// What is needed

	// 1. Disk to read/write to
	// 2. Source/Destination to read/write from
	// 3. Post tasks
	// --- 1. Disk stretch
	// --- 2. Post config?
	mac, err := realm.GetMAC()
	if err != nil {
		log.Errorln(err)
		realm.Shell()
	}

	cfg, err := plunderclient.GetConfigForAddress(utils.DashMac(mac))

	if err != nil {
		log.Errorf("Error with remote server [%v]", err)
		log.Errorln("Rebooting in 10 seconds")
		time.Sleep(time.Second * 10)
		realm.Reboot()
	}

	switch cfg.Action {
	case types.ReadImage:
		err = image.Read(cfg.SourceDevice, cfg.DesintationAddress)
		if err != nil {
			log.Errorf("Read Image Error: [%v]", err)
		}

	case types.WriteImage:
		err = image.Write(cfg.SourceImage, cfg.DestinationDevice)
		if err != nil {
			log.Errorf("Write Image Error: [%v]", err)

		}

	default:
		log.Errorf("Unknown action [%s] passed to deployment image, restarting in 10 seconds", cfg.Action)
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

	if cfg.DropToShell == true {
		realm.Shell()
	}

	realm.Reboot()
}
