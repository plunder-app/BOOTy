package plunderclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetServerConfig will retrieve the configuraiton for a server (mac address)
func GetServerConfig() error {
	// Attempt to find the Server URL
	url := os.Getenv("BOOTYURL")
	if url == "" {
		return fmt.Errorf("The flag BOOTYURL is empty")
	}
	log.Infof("Connecting to provisioning server [%s]", url)

	plunderClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "BOOTy-client")

	res, err := plunderClient.Do(req)
	if err != nil {
		return err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return err
	}

	fmt.Printf("\n%s\n", string(body))
	return nil
}
