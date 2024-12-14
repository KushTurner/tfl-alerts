package tfl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllDisruptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"description": "Elizabeth line: Something happened", "closureText": "severeDelays"}]`))
	}))

	defer server.Close()

	td := []TrainDisruption{{Description: "Elizabeth line: Something happened", ClosureText: "severeDelays"}}

	tfl := TflService{&http.Client{}, server.URL}
	resp, _ := tfl.AllCurrentDisruptions()

	assert.Equal(t, td, resp)
}
