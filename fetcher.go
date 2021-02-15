package vatspy

import (
	"io/ioutil"
	"net/http"
)

// Load loads and parses data from a local file
func Load(dataFilename string, boundariesFilename string) (*Data, error) {
	rawData, err := ioutil.ReadFile(dataFilename)
	if err != nil {
		return nil, err
	}
	rawBoundaries, err := ioutil.ReadFile(boundariesFilename)
	if err != nil {
		return nil, err
	}
	return parse(rawData, rawBoundaries)
}

// Fetch fetches and parses data from an HTTP url
func Fetch(dataURL string, boundariesURL string) (*Data, error) {
	resp, err := http.Get(dataURL)
	if err != nil {
		return nil, err
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp, err = http.Get(boundariesURL)
	if err != nil {
		return nil, err
	}
	rawBoundaries, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parse(rawData, rawBoundaries)
}
