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
		_, _ = w.Write([]byte(`[{"description": "Elizabeth line: Something happened", "closureText": "severeDelays"}]`))
	}))

	defer server.Close()

	td := []TrainDisruption{{Description: "Elizabeth line: Something happened", ClosureText: "severeDelays"}}

	tfl := FakeTflClient(server.URL)
	resp, _ := tfl.AllCurrentDisruptions()

	assert.Equal(t, td, resp)
}

func TestGetDisruptedTrain(t *testing.T) {
	td := TrainDisruption{Description: "Elizabeth Line: Something happened", ClosureText: "severeDelays"}

	assert.Equal(t, "Elizabeth Line", td.getDisruptedTrain())
}
