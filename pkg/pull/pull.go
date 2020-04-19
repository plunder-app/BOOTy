package pull

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/schollz/progressbar"
	"github.com/thebsdbox/BOOTy/pkg/utils"
)

// Image will pull an image and write it to local storage device
func Image() {
	path := os.Getenv("CMDLINEPATH")
	if path == "" {
		path = utils.CmdlinePath
	}
	utils.ClearScreen()
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Starting BOOTy \n")
	fmt.Printf("\n\n")
	fmt.Printf("Parsing config from [%s]\n", path)
	src, dst, err := utils.ParseCmdLine(path)
	if err != nil {
		log.Fatalf("%v", err)
	}

	envSrc := os.Getenv("SRC")
	envDst := os.Getenv("DST")
	if envSrc == "" {
		//fmt.Printf("The \"SRC\" environment variable wasn't set")
	} else {
		src = envSrc
	}

	if envDst == "" {
		//fmt.Printf("The \"DST\" environment variable wasn't set")
	} else {
		dst = envDst
	}

	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	var out io.Writer
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	out = f
	defer f.Close()

	bar := progressbar.NewOptions(
		int(resp.ContentLength),
		progressbar.OptionShowBytes(true),
	)
	out = io.MultiWriter(out, bar)
	fmt.Printf("\n\n\n")

	fmt.Printf("Beginning write of image [%s] to disk [%s]", filepath.Base(src), dst)
	fmt.Printf("\n\n\n")

	count, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Error writing %d bytes to disk [%s] -> %v", count, dst, err)
	}
	fmt.Printf("\n\n\n")

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("\n\n\n\n")

	// TODO - reboot
	fmt.Println("This is where the reboot happens :-D")

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return
	}

	time.Sleep(time.Second * 5)

	//set keys
	kb.HasCTRL(true)
	kb.HasALT(true)
	kb.SetKeys(111) // Delete

	//launch
	kb.Launching()
}
