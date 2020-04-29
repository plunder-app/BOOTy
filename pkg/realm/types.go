package realm

import "os"

// Mount contains the configuration for a single mount within the initramfs
type Mount struct {
	// Create the location on disk
	CreateMount bool
	// Enable the mount Source -> Path w/options
	EnableMount bool

	// Configurations
	Name    string
	Source  string
	Path    string
	Mode    os.FileMode
	FSType  string
	Flags   uintptr
	Options string
}

// Mounts are the paths that can be mounted or created on boot
type Mounts struct {
	Mount []Mount
}

// Device contains the configuration for a single device within the initramfs
type Device struct {
	// Create the device within the ramdisk
	CreateDevice bool

	// Configuration for the device
	Name  string
	Path  string
	Mode  uint32
	Major int64
	Minor int64
}

// Devices are the devices that can be created on boot
type Devices struct {
	Device []Device
}