package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func FakeTflClient(url string) *TflClient {
	return &TflClient{&http.Client{}, url}
}

func TestGetAllDisruptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(
			`[{"name":"Avanti West Coast","lineStatuses":[{"statusSeverity":0,"statusSeverityDescription":"Special Service","reason":"https://www.nationalrail.co.uk/service-disruptions/polesworth-20241204/"}]}]`))
	}))

	defer server.Close()

	td := []TrainStatus{{Name: "Avanti West Coast", LineStatuses: []LineStatus{{StatusSeverity: 0, StatusSeverityDescription: "Special Service", Reason: "https://www.nationalrail.co.uk/service-disruptions/polesworth-20241204/"}}}}

	tfl := FakeTflClient(server.URL)
	resp, _ := tfl.AllCurrentDisruptions()

	assert.Equal(t, td, resp)
}
