package tfl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllDisruptions(t *testing.T) {
	t.Run("Can make request TFL to get all disruptions", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/Line/Mode/tube,overground,national-rail,elizabeth-line,dlr/Status", r.URL.Path)
			assert.Equal(t, "true", r.URL.Query().Get("detail"))
			assert.Equal(t, "", r.Header.Get("User-Agent"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(
				`[
				  {
					"name": "Avanti West Coast",
					"lineStatuses": [
					  {
						"statusSeverity": 9,
						"statusSeverityDescription": "Minor Delays",
						"reason": "https://www.nationalrail.co.uk/service-disruptions/polesworth-20241204/",
						"disruption": {
						  "summary": "Minor delays between London Euston and Birmingham due to an earlier fault"
						}
					  }
					]
				  }
				]`))
		}))
		defer s.Close()

		expectedTrainStatus := []TrainStatus{
			{
				Name: "Avanti West Coast",
				LineStatuses: []LineStatus{
					{
						StatusSeverity:            9,
						StatusSeverityDescription: "Minor Delays",
						Reason:                    "https://www.nationalrail.co.uk/service-disruptions/polesworth-20241204/",
						Disruption: Disruption{
							Summary: "Minor delays between London Euston and Birmingham due to an earlier fault",
						},
					},
				},
			},
		}

		tfl, _ := NewClient(&Config{Url: s.URL})
		actualTrainStatus, _ := tfl.AllCurrentDisruptions()

		assert.Equal(t, expectedTrainStatus, actualTrainStatus)
	})
}
