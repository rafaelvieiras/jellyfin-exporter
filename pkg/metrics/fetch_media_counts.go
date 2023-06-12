package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/api"
)

var mediaCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "jellyfin_media_count",
		Help: "Count of media items in Jellyfin",
	},
	[]string{"title", "type"},
)

func init() {
	prometheus.MustRegister(mediaCount)
}

// FetchMediaCounts retrieves the media counts and updates the metric
func FetchMediaCounts(jellyfinApiUrl, jellyfinToken string) {
	countsUrl := jellyfinApiUrl + "/Items/Counts"

	// Make the request to the Jellyfin API
	rawData, err := api.MakeRequest(countsUrl, jellyfinToken)

	if err != nil {
		log.Printf("Error fetching API: %v", err)
		return
	}

	// Type assertion for map[string]interface{}
	data, ok := rawData.(map[string]interface{})
	if !ok {
		log.Printf("Error parsing data: expected a map")
		return
	}

	// Mapping key names to more readable titles and types
	mapping := map[string]map[string]string{
		"MovieCount":      {"title": "Movies", "type": "movie"},
		"SeriesCount":     {"title": "Shows", "type": "show"},
		"EpisodeCount":    {"title": "Shows - Episodes", "type": "show_episode"},
		"ArtistCount":     {"title": "Artists", "type": "artist"},
		"ProgramCount":    {"title": "Programs", "type": "program"},
		"TrailerCount":    {"title": "Trailers", "type": "trailer"},
		"SongCount":       {"title": "Songs", "type": "song"},
		"AlbumCount":      {"title": "Albums", "type": "album"},
		"MusicVideoCount": {"title": "Music Videos", "type": "music_video"},
		"BoxSetCount":     {"title": "Box Sets", "type": "box_set"},
		"BookCount":       {"title": "Books", "type": "book"},
		"ItemCount":       {"title": "Items", "type": "item"},
	}

	// Update the metrics with the values obtained from the API
	mediaCount.Reset()
	for key, values := range mapping {
		value, exists := data[key].(float64)
		if !exists {
			log.Printf("Error parsing JSON for key %s", key)
			continue
		}
		mediaCount.With(prometheus.Labels{
			"title": values["title"],
			"type":  values["type"],
		}).Set(value)
	}
}
