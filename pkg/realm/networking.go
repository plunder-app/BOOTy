//+build linux

package realm

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/digineo/go-dhclient"
	"github.com/google/gopacket/layers"
)

const ifname = "eth0"

// LeasedAddress is the currently leased address
var LeasedAddress string

// GetMAC will return a mac address
func GetMAC() (string, error) {
	// retrieve interface from name
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return "", err
	}
	return iface.HardwareAddr.String(), nil
}

// DHCPClient starts the DHCP client listening for a lease
func DHCPClient() error {

	// Bring up interface
	ifaceDev, err := netlink.LinkByName(ifname)
	if err != nil {
		log.Errorf("Error finding adapter [%v]", err)

		return err
	}

	if err := netlink.LinkSetUp(ifaceDev); err != nil {
		log.Errorf("Error bringing up adapter [%v]", err)
	}

	// Setup interface to recieve DHCP traffic
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Errorf("Error finding interface by name [%v]", err)

		return err
	}
	client := dhclient.Client{
		Iface: iface,
		OnBound: func(lease *dhclient.Lease) {
			// Set the lease string to be used in other places
			LeasedAddress = lease.FixedAddress.String()

			link, _ := netlink.LinkByName(iface.Name)

			// Set address / netmask into cidr we can use to apply to interface
			cidr := net.IPNet{
				IP:   lease.FixedAddress,
				Mask: lease.Netmask,
			}
			addr, _ := netlink.ParseAddr(cidr.String())

			err = netlink.AddrAdd(link, addr)
			if err != nil {
				log.Errorf("Error adding %s to link %s", cidr.String(), iface.Name)
			} else {
				log.Printf("Adding address %s to link %s", cidr.String(), iface.Name)
			}

			// Apply default gateway so we can route outside
			route := netlink.Route{
				Scope: netlink.SCOPE_UNIVERSE,
				Gw:    lease.ServerID,
			}
			if err := netlink.RouteAdd(&route); err != nil {
				log.Errorf("Error setting gateway [%v]", err)
			} else {
				log.Printf("Adding gateway %s to link %s", lease.ServerID.String(), iface.Name)
			}
		},
	}

	// Add requests for default options
	for _, param := range dhclient.DefaultParamsRequestList {
		log.Printf("Requesting default option %d", param)
		client.AddParamRequest(layers.DHCPOpt(param))
	}

	// // Add requests for custom options
	// for _, param := range requestParams {
	// 	log.Printf("Requesting custom option %d", param)
	// 	client.AddParamRequest(layers.DHCPOpt(param))
	// }

	// Add hostname option
	hostname, _ := os.Hostname()
	client.AddOption(layers.DHCPOptHostname, []byte(hostname))

	// // Add custom options
	// for _, option := range options {
	// 	log.Printf("Adding option %d=0x%x", option.Type, option.Data)
	// 	client.AddOption(option.Type, option.Data)
	// }

	client.Start()
	defer client.Stop()

	// Below will sit

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)
	for {
		sig := <-c
		log.Println("received", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			return nil
		case syscall.SIGHUP:
			log.Println("renew lease")
			client.Renew()
		case syscall.SIGUSR1:
			log.Println("acquire new lease")
			client.Rebind()
		}
	}
	//log.Errorf("DHCP client has ended")
	//return nil
}
