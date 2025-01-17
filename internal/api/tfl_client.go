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

type TrainStatus struct {
	Name         string       `json:"name"`
	LineStatuses []LineStatus `json:"lineStatuses"`
}

type LineStatus struct {
	StatusSeverity            int    `json:"statusSeverity"`
	StatusSeverityDescription string `json:"statusSeverityDescription"`
	Reason                    string `json:"reason"`
}

func NewTflClient(cfg *TflConfig) (TflClient, error) {
	return TflClient{&http.Client{}, cfg.Url}, nil
}

func (c *TflClient) AllCurrentDisruptions() ([]TrainStatus, error) {
	trainType := "tube,overground,national-rail,elizabeth-line,dlr"

	resp, err := c.get(c.url + "/Line/Mode/" + trainType + "/Status")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch status: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var t []TrainStatus

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
