package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dustin/go-humanize"
	"github.com/thebsdbox/BOOTy/pkg/plunderclient/types"
	"github.com/thebsdbox/BOOTy/pkg/utils"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

var data []byte

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

//PrintProgress -
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
	fmt.Println("")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {

	imageName := fmt.Sprintf("%s.img", r.RemoteAddr)

	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("BootyImage")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	out, err := os.OpenFile(imageName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer out.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(file, counter)); err != nil {
		log.Errorf("%v", err)
	}

	fmt.Printf("Beginning write of image [%s] to disk", imageName)

	w.WriteHeader(http.StatusOK)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(data)
}

// Serve will start the webserver for BOOTy
func main() {

	// Server Address
	rawAddress := flag.String("address", "", "The mac address of a server")
	var address string

	// Build configuration from flags
	var config types.BootyConfig
	flag.StringVar(&config.Action, "action", "", "The action that is being performed [readImage/writeImage]")
	flag.BoolVar(&config.DryRun, "dryRun", false, "Only demonstrate the output from the actions")
	flag.BoolVar(&config.DropToShell, "shell", false, "Start a shell")

	flag.StringVar(&config.SourceImage, "sourceImage", "", "The source for the image, typically a URL")
	flag.StringVar(&config.SourceDevice, "sourceDevice", "", "The device that will be the source of the image [/dev/sda]")

	flag.StringVar(&config.DesintationAddress, "destinationAddress", "", "The destination that the image will be writen too [url]")
	flag.StringVar(&config.DestinationDevice, "destinationDevice", "", "The destination devicethat the image will be writen too [/dev/sda]")
	flag.Parse()

	if *rawAddress == "" {
		log.Warnln("No Mac address passed for BOOTy configuration")
	} else {

		address = utils.DashMac(*rawAddress)
		http.HandleFunc(fmt.Sprintf("/booty/%s.bty", address), configHandler)
		log.Infof("handler for [%s.bty] generated", address)
		data, _ = json.Marshal(config)
	}

	switch config.Action {
	case types.ReadImage:
	case types.WriteImage:
	default:
		log.Fatalf("Unknown action [%s]", config.Action)
	}

	fs := http.FileServer(http.Dir("./images"))
	http.HandleFunc("/image", imageHandler)
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
