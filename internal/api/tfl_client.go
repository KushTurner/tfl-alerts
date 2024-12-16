package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TflClient struct {
	client *http.Client
	url    string
}

type TrainDisruption struct {
	Description string `json:"description"`
	ClosureText string `json:"closureText"`
}

func NewTflClient() *TflClient {
	return &TflClient{&http.Client{}, "https://api.tfl.gov.uk"}
}

func (c *TflClient) AllCurrentDisruptions() ([]TrainDisruption, error) {
	trainType := "tube"

	req, err := http.NewRequest("GET", c.url+"/Line/Mode/"+trainType+"/Disruption", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch disruptions %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var t []TrainDisruption

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode disruptions %v", err)
	}

	return t, nil
}
