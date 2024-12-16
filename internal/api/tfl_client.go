package api

import (
	"encoding/json"
	"fmt"
	"github.com/kushturner/tfl-alerts/internal/config"
	"net/http"
)

type TflClient struct {
	client *http.Client
	url    string
	appId  string
}

type TrainDisruption struct {
	Description string `json:"description"`
	ClosureText string `json:"closureText"`
}

func NewTflClient(cfg *config.TflConfig) *TflClient {
	return &TflClient{&http.Client{}, "https://api.tfl.gov.uk", cfg.AppId}
}

func (c *TflClient) AllCurrentDisruptions() ([]TrainDisruption, error) {
	trainType := "tube"

	req, err := http.NewRequest("GET", c.url+"/Line/Mode/"+trainType+"/Disruption", nil)
	req.Header.Set("User-Agent", "tfl-alerts")

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
