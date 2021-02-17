package dynamic

import (
	"io/ioutil"
	"net/http"
)

const (
	// VatSimJSON3URL is the default VatSim data URL
	VatSimJSON3URL = "https://data.vatsim.net/v3/vatsim-data.json"
)

// Fetch fetches raw data from VATSIM
func Fetch(dataURL string) (*Data, error) {
	resp, err := http.Get(dataURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return newData(data)
}
