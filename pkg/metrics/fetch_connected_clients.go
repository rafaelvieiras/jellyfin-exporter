package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/api"
)

var (
	connectedClientsCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_connected_clients_count",
			Help: "Number of connected clients to Jellyfin",
		},
		[]string{"user_id", "username", "client", "device"},
	)
)

func init() {
	prometheus.MustRegister(connectedClientsCount)
}

// FetchConnectedClients retrieves the connected clients and updates the metric
func FetchConnectedClients(jellyfinApiUrl, jellyfinToken string) {
	sessionsUrl := jellyfinApiUrl + "/Sessions"

	// Make the request to the Jellyfin API
	rawData, err := api.MakeRequest(sessionsUrl, jellyfinToken)

	if err != nil {
		log.Printf("Error fetching API: %v", err)
		return
	}

	// Type assertion for []interface{}
	data, ok := rawData.([]interface{})
	if !ok {
		log.Printf("Error parsing data: expected an array")
		return
	}

	// Update the metric with the number of connected clients (this depends on the structure of the JSON response)
	connectedClientsCount.Reset()
	for _, session := range data {
		sessionMap, ok := session.(map[string]interface{})
		if !ok {
			continue
		}

		connectedClientsCount.With(prometheus.Labels{
			"user_id":  sessionMap["UserId"].(string),
			"username": sessionMap["UserName"].(string),
			"client":   sessionMap["Client"].(string),
			"device":   sessionMap["DeviceName"].(string),
		}).Inc()
	}
}
