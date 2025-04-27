package tfl

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Url string
}

type Client struct {
	httpClient *http.Client
	url        string
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

const trainTypes string = "tube,overground,national-rail,elizabeth-line,dlr"

func NewClient(cfg *Config) (Client, error) {
	return Client{&http.Client{}, cfg.Url}, nil
}

func (c *Client) AllCurrentDisruptions() ([]TrainStatus, error) {
	resp, err := c.get(c.url + "/Line/Mode/" + trainTypes + "/Status")
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

func (c *Client) get(url string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "")

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
