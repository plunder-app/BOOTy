package types

// BOOTy example 1

// Action - writeImage
// SourceImage - ubuntu.img
// DestinationDevice - /dev/sda

// BOOTy example 2

// Action - readImage
// SourceDevice - ubuntu.img

// BOOTy example 3

// Action - writeImage
// SourceImage - centos.img
// DestinationDevice - /dev/sda
// GrowDisk - true

// ------------------ //

const (
	//ReadImage means that this is a read only action
	ReadImage = "readImage"

	// WriteImage means that this is a read/write action
	WriteImage = "writeImage"
)

// BootyConfig defines the data passed to the BOOTy initramdisk
type BootyConfig struct {
	// Defines what action the deployment will take
	Action string `json:"action"`

	// Disk actions ->

	// Data should be compressed or is compressed
	Compressed bool `json:"compressed"`

	// Write image to disk from remote address
	SourceImage       string `json:"sourceImage,omitempty"`
	DestinationDevice string `json:"destinationDevice,omitempty"`

	// Read Image from Disk and write to remote address
	SourceDevice       string `json:"sourceDevice,omitempty"`
	DesintationAddress string `json:"desintationAddress,omitempty"`

	// Post tasks - Once the image has been deployed

	// Disk modifications
	GrowDisk bool `json:"growDisk"`

	// Volume modifications (LVM2)
	GrowPartition int    `json:"growPartition"`
	LVMRootName   string `json:"lvmRootName"`

	// Network modifcations
	Address    string `json:"address,omitempty"`
	Gateway    string `json:"gateway,omitempty"`
	NameServer string `json:"nameserver,omitempty"`

	// Debugging, troubleshooting
	DryRun      bool `json:"dryRun"`
	DropToShell bool `json:"dropToShell"`
	WipeDevice  bool `json:"wipeDevice"`
}
