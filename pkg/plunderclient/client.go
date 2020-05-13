package plunderclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/plunder-app/BOOTy/pkg/plunderclient/types"
	log "github.com/sirupsen/logrus"
)

// GetConfigForAddress will retrieve the configuraiton for a server (mac address)
func GetConfigForAddress(mac string) (*types.BootyConfig, error) {
	// Attempt to find the Server URL
	url := os.Getenv("BOOTYURL")
	if url == "" {
		return nil, fmt.Errorf("The flag BOOTYURL is empty")
	}
	log.Infof("Connecting to provisioning server [%s]", url)

	// Address format

	// http:// address / booty / <mac> .bty

	// url = http://address/booty
	configURL := fmt.Sprintf("%s/%s.bty", url, mac)
	plunderClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	req, err := http.NewRequest(http.MethodGet, configURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "BOOTy-client")

	res, err := plunderClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 300 {
		// Customise response for the 404 to make degugging simpler
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("%s not found", configURL)
		}
		return nil, fmt.Errorf("%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var config types.BootyConfig

	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Errorf("Error reading [%s]", configURL)
		return nil, err
	}

	return &config, nil
}
