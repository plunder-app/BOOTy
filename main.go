package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/schollz/progressbar"
)

const cmdlinePath = "/proc/cmdline"

func parseCmdLine(path string) (src, dst string, err error) {

	// Read the file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Split by whitespace
	entries := strings.Fields(string(b))

	// find k=v entries
	for x := range entries {
		kv := strings.Split(entries[x], "=")
		if len(kv) == 2 {
			// find the entries we care about
			if kv[0] == "BOOTYSRC" {
				src = kv[1]
			}
			if kv[0] == "BOOTYDST" {
				dst = kv[1]
			}
		}
	}
	return
}

func clear() {
	fmt.Print("\033[2J")
}

func main() {

	if os.Getenv("SERVER") != "" {
		fs := http.FileServer(http.Dir("./images"))
		http.Handle("/images/", http.StripPrefix("/images/", fs))
		log.Println("Listening on :3000...")
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	path := os.Getenv("CMDLINEPATH")
	if path == "" {
		path = cmdlinePath
	}
	clear()
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Starting BOOTy \n")
	fmt.Printf("\n\n")
	fmt.Printf("Parsing config from [%s]\n", path)
	src, dst, err := parseCmdLine(path)
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
		log.Fatalf("Error writing %d bytes to [%s] -> %v", count, filepath.Base(src), err)
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
