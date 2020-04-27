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

	// TODO - consider timeouts

	client := &http.Client{}
	_, err = UploadMultipartFile(client, dst, "BootyImage", src)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\n\n\n")

	fmt.Printf("Reading of disk [%s], and sending to [%s]", filepath.Base(src), dst)
	fmt.Printf("\n\n\n")

	fmt.Printf("\n\n\n")

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("\n\n\n\n")

	// TODO - reboot
	fmt.Println("This is where the push reboot happens :-D")

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

//UploadMultipartFile -
func UploadMultipartFile(client *http.Client, uri, key, path string) (*http.Response, error) {
	body, writer := io.Pipe()

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}

	mwriter := multipart.NewWriter(writer)
	req.Header.Add("Content-Type", mwriter.FormDataContentType())

	errchan := make(chan error)

	go func() {
		defer close(errchan)
		defer writer.Close()
		defer mwriter.Close()

		w, err := mwriter.CreateFormFile(key, path)
		if err != nil {
			errchan <- err
			return
		}

		in, err := os.Open(path)
		if err != nil {
			errchan <- err
			return
		}

		defer in.Close()

		if written, err := io.Copy(w, in); err != nil {
			errchan <- fmt.Errorf("error copying %s (%d bytes written): %v", path, written, err)
			return
		}

		if err := mwriter.Close(); err != nil {
			errchan <- err
			return
		}
	}()

	resp, err := client.Do(req)
	merr := <-errchan

	if err != nil || merr != nil {
		return resp, fmt.Errorf("http error: %v, multipart error: %v", err, merr)
	}

	return resp, nil
}
