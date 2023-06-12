package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/config"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/metrics"
)

func main() {
	// Load the environment variables
	err := config.LoadEnvironment(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Get the environment variables
	jellyfinApiUrl := os.Getenv("JELLYFIN_API_URL")
	jellyfinToken := os.Getenv("JELLYFIN_TOKEN")
	metricsPath := os.Getenv("METRICS_PATH")
	serverAddr := os.Getenv("SERVER_ADDR")
	serverPort := os.Getenv("SERVER_PORT")

	// Validate the environment variables
	if jellyfinApiUrl == "" {
		log.Fatal("JELLYFIN_API_URL not set")
	}
	if jellyfinToken == "" {
		log.Fatal("JELLYFIN_TOKEN not set")
	}
	if metricsPath == "" {
		log.Fatal("METRICS_PATH not set")
	}
	if serverPort == "" {
		log.Fatal("SERVER_PORT not set")
	}

	// Collect metrics periodically
	go func() {
		for {
			metrics.FetchConnectedClients(jellyfinApiUrl, jellyfinToken)
			metrics.FetchMediaCounts(jellyfinApiUrl, jellyfinToken)
			metrics.FetchStreamCounts(jellyfinApiUrl, jellyfinToken)
			metrics.FetchScheduledTasks(jellyfinApiUrl, jellyfinToken)

			// Wait before collecting again (ex: 15 seconds)
			time.Sleep(15 * time.Second)
		}
	}()

	// Configure the server and route
	http.Handle(metricsPath, promhttp.Handler())
	server := &http.Server{
		Addr:         serverAddr + ":" + serverPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server
	log.Printf("Starting server on port %s", serverPort)
	log.Fatal(server.ListenAndServe())
}
