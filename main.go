package main

import (
	"time"

	"github.com/plunder-app/BOOTy/pkg/image"
	"github.com/plunder-app/BOOTy/pkg/plunderclient"
	"github.com/plunder-app/BOOTy/pkg/plunderclient/types"
	"github.com/plunder-app/BOOTy/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/plunder-app/BOOTy/pkg/realm"
	"github.com/plunder-app/BOOTy/pkg/ux"
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
	m.MountNamed("dev", true)

	// Create all devices
	d.CreateDevice()

	// Mount any additional mounts
	m.MountAll()

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
		err = image.Read(cfg.SourceDevice, cfg.DesintationAddress, mac, cfg.Compressed)
		if err != nil {
			log.Errorf("Read Image Error: [%v]", err)
			onError(cfg)
		}
		log.Infoln("Image written succesfully, restarting in 5 seconds")
		time.Sleep(time.Second * 5)
		realm.Reboot()

	case types.WriteImage:
		err = image.Write(cfg.SourceImage, cfg.DestinationDevice, cfg.Compressed)
		if err != nil {
			log.Errorf("Write Image Error: [%v]", err)
			onError(cfg)
		}
		// log.Infoln("Image written succesfully, restarting in 5 seconds")
		// time.Sleep(time.Second * 5)
		// realm.Reboot()

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

	rv, err := realm.MountRootVolume(cfg.LVMRootName)
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
		onError(cfg)
	}

	err = realm.GrowLVMRoot(cfg.DestinationDevice, cfg.LVMRootName, cfg.GrowPartition)
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
		onError(cfg)
	}

	// Start the networking configuration (UBUNTU ONLY)
	log.Infoln("Starting Networking configuration")
	err = realm.WriteNetPlan("/mnt", cfg)
	if err != nil {
		log.Errorf("Network Error: [%v]", err)
		onError(cfg)
	}

	// Apply the networking configuration (UBUNTU ONLY)
	log.Infoln("Applying Networking configuration")
	err = realm.ApplyNetplan("/mnt")
	if err != nil {
		log.Errorf("Network Error: [%v]", err)
		onError(cfg)
	}

	log.Infoln("Un Mounting boot volume")
	err = rv.UnMountNamed("dev")
	if err != nil {
		log.Errorf("UnMounting Error: [%v]", err)
		onError(cfg)
	}
	err = rv.UnMountNamed("proc")
	if err != nil {
		log.Errorf("UnMounting Error: [%v]", err)
		onError(cfg)
	}
	err = rv.UnMountAll()
	if err != nil {
		log.Errorf("UnMounting Error: [%v]", err)
		onError(cfg)
	}

	if cfg.DropToShell == true {
		realm.Shell()
	}

	log.Infoln("BOOTy is now exiting, system will reboot")
	time.Sleep(time.Second * 2)
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
	// Time to see the error
	time.Sleep(time.Second * 2)
	realm.Reboot()
}
