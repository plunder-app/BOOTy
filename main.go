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
			onError(cfg)
		}

	case types.WriteImage:
		err = image.Write(cfg.SourceImage, cfg.DestinationDevice)
		if err != nil {
			log.Errorf("Write Image Error: [%v]", err)
			onError(cfg)
		}

	default:
		log.Errorf("Unknown action [%s] passed to deployment image, restarting in 10 seconds", cfg.Action)
		time.Sleep(time.Second * 10)
		realm.Reboot()
	}

	log.Infoln("Beginning Disk Management")

	err = realm.PartProbe(cfg.DestinationDevice)
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
		onError(cfg)
	}

	err = realm.EnableLVM()
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
		onError(cfg)
	}

	err = realm.MountRootVolume(cfg.LVMRootName)
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
		onError(cfg)
	}

	err = realm.GrowRoot(cfg.DestinationDevice, cfg.LVMRootName, cfg.GrowPartition)
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
		onError(cfg)
	}

	err = realm.UnMount("/mnt")
	if err != nil {
		log.Errorf("UnMounting Error: [%v]", err)
		onError(cfg)
	}

	if cfg.DropToShell == true {
		realm.Shell()
	}

	realm.Reboot()

}

// on Error we will execute the following steps
func onError(cfg *types.BootyConfig) {

	if cfg.WipeDevice == true {
		realm.Wipe(cfg.DestinationDevice)
	}

	if cfg.DropToShell == true {
		realm.Shell()
	}

	realm.Reboot()
}
