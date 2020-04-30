// +build linux

package ux

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/zcalusic/sysinfo"
)

func SysInfo() {
	var si sysinfo.SysInfo

	si.GetSysInfo()
	fmt.Println("")
	fmt.Println("------------ BOOTy System Information ------------")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	// 44  	fmt.Fprintln(w, "a\tb\tc")
	// 45  	fmt.Fprintln(w, "aa\tbb\tcc")
	// 46  	fmt.Fprintln(w, "aaa\t") // trailing tab
	// 47  	fmt.Fprintln(w, "aaaa\tdddd\teeee")
	// 48  	w.Flush()

	fmt.Fprintf(w, "CPU:\t %s\n", si.CPU.Model)
	fmt.Fprintf(w, "CPU speed:\t %dMHz\n", si.CPU.Speed)
	fmt.Fprintf(w, "MEM size:\t %dMB\n", si.Memory.Size)
	for x := range si.Network {
		fmt.Fprintf(w, "Network device:\t %s\n", si.Network[x].Name)
		fmt.Fprintf(w, "Network driver:\t %s\n", si.Network[x].Driver)
		fmt.Fprintf(w, "Network address:\t %s\n", si.Network[x].MACAddress)
	}
	for x := range si.Storage {
		fmt.Fprintf(w, "Storage device:\t %s\n", si.Storage[x].Name)
		fmt.Fprintf(w, "Storage driver:\t %s\n", si.Storage[x].Driver)
		fmt.Fprintf(w, "Storage size:\t %dGB\n", si.Storage[x].Size)
	}
	w.Flush()

	fmt.Println("--------------------------------------------------")
	fmt.Println("")

}
