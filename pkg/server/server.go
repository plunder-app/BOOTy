package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/schollz/progressbar"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {

	imageName := fmt.Sprintf("%s.img", r.Host)

	var out io.Writer
	f, err := os.OpenFile(imageName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	out = f
	defer f.Close()

	bar := progressbar.NewOptions(
		int(r.ContentLength),
		progressbar.OptionShowBytes(true),
	)
	out = io.MultiWriter(out, bar)
	fmt.Printf("\n\n\n")

	fmt.Printf("Beginning write of image [%s] to disk", imageName)
	fmt.Printf("\n\n\n")

	count, err := io.Copy(out, r.Body)
	if err != nil {
		log.Fatalf("Error writing %d bytes to [%s] -> %v", count, imageName, err)
	}
	w.WriteHeader(http.StatusOK)
}

// Serve will start the webserver for BOOTy
func Serve() {

	fs := http.FileServer(http.Dir("./images"))
	http.HandleFunc("/image", imageHandler)
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
