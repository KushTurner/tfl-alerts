package tfl

import (
	"encoding/json"
	"net/http"
)

type TflService struct {
	client *http.Client
	url    string
}

type TrainDisruption struct {
	Description string `json:"description"`
	ClosureText string `json:"closureText"`
}

func NewService() TflService {
	return TflService{&http.Client{}, "https://api.tfl.gov.uk"}
}

func (s *TflService) AllCurrentDisruptions() ([]TrainDisruption, error) {
	req, _ := http.NewRequest("GET", s.url+"/Line/Mode/tube/Disruption", nil)

	resp, _ := s.client.Do(req)
	defer resp.Body.Close()

	var t []TrainDisruption

	_ = json.NewDecoder(resp.Body).Decode(&t)

	return t, nil
}
