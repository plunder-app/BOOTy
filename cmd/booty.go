package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Release - this struct contains the release information populated when building booty
var Release struct {
	Version string
	Build   string
}

var bootyCmd = &cobra.Command{
	Use:   "booty",
	Short: "This is a tool for managing boot images for bare-metal servers",
}

func init() {
	// bootyCmd.AddCommand(bootyPull)
	// bootyCmd.AddCommand(bootyPush)
	// bootyCmd.AddCommand(bootyServer)
	bootyCmd.AddCommand(bootyVersion)

}

// Execute - starts the command parsing process
func Execute() {
	if err := bootyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// var bootyPull = &cobra.Command{
// 	Use:   "pull",
// 	Short: "This is will direct BOOTy to pull and image from a remote server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		pull.Image()
// 	},
// }

// var bootyPush = &cobra.Command{
// 	Use:   "push",
// 	Short: "This is will direct BOOTy to push the contents of a disk to a remote server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		push.Image()
// 	},
// }

// var bootyServer = &cobra.Command{
// 	Use:   "server",
// 	Short: "This is for starting BOOTy as a simple (test) web server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		server.Serve()
// 	},
// }

var bootyVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the BOOTy image manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BOOTy Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}
