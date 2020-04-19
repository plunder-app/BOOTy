package push

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/thebsdbox/BOOTy/pkg/utils"
)

// Image - will take a local disk and copy an image to a remote server
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

	// req, err := http.NewRequest("GET", src, nil)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// defer resp.Body.Close()

	// var out io.Writer
	// f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// out = f
	// defer f.Close()

	// bar := progressbar.NewOptions(
	// 	int(resp.ContentLength),
	// 	progressbar.OptionShowBytes(true),
	// )
	// out = io.MultiWriter(out, bar)
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("file", filepath.Base(src))
		if err != nil {
			log.Fatalf("%v", err)
		}
		file, err := os.Open(src)
		if err != nil {
			return
		}
		defer file.Close()
		if count, err := io.Copy(part, file); err != nil {
			log.Fatalf("Error writing %d bytes to [%s] -> %v", count, filepath.Base(src), err)
		}
	}()
	http.Post(dst, m.FormDataContentType(), r)

	fmt.Printf("\n\n\n")

	fmt.Printf("Beginning write of image [%s] to disk [%s]", filepath.Base(src), dst)
	fmt.Printf("\n\n\n")

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
