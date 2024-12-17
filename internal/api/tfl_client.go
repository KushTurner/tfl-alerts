package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TflConfig struct {
	Url string
}

type TflClient struct {
	client *http.Client
	url    string
}

type TrainDisruption struct {
	Description string `json:"description"`
	ClosureText string `json:"closureText"`
}

func NewTflClient(cfg *TflConfig) (*TflClient, error) {
	return &TflClient{&http.Client{}, cfg.Url}, nil
}

func (c *TflClient) AllCurrentDisruptions() ([]TrainDisruption, error) {
	trainType := "tube"

	resp, err := c.get(c.url + "/Line/Mode/" + trainType + "/Disruption")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch disruptions %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var t []TrainDisruption

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode disruptions: %v", err)
	}

	return t, nil
}

func (c *TflClient) get(url string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "")

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
