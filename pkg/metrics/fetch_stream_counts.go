package metrics

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/api"
)

var streamCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "jellyfin_stream_count",
		Help: "Count of streams in Jellyfin",
	},
	[]string{"user_id", "username", "client", "device", "play_method", "audio_codec", "video_codec", "container", "audio_channels", "transcode_reasons"},
)

func init() {
	prometheus.MustRegister(streamCount)
}

// FetchStreamCounts retrieves the stream counts and updates the metric
func FetchStreamCounts(jellyfinApiUrl, jellyfinToken string) {
	sessionsUrl := fmt.Sprintf("%s/Sessions", jellyfinApiUrl)

	// Faça a requisição para a API do Jellyfin
	rawData, err := api.MakeRequest(sessionsUrl, jellyfinToken)
	if err != nil {
		log.Printf("Error fetching API: %v", err)
		return
	}

	// Faça um type assertion para []interface{}
	data, ok := rawData.([]interface{})
	if !ok {
		log.Printf("Error parsing data: expected an array")
		return
	}

	// Resetar a métrica antes de populá-la
	streamCount.Reset()

	// Iterar sobre as sessões em andamento
	for _, session := range data {
		sessionMap := session.(map[string]interface{})
		playStateMap := sessionMap["PlayState"].(map[string]interface{})

		playMethod := ""
		audioCodec := ""
		videoCodec := ""
		container := ""
		audioChannels := 0
		transcodeReasons := ""

		if playStateMap["PlayMethod"] != nil {
			playMethod = playStateMap["PlayMethod"].(string)
		}

		if sessionMap["AudioCodec"] != nil {
			audioCodec = sessionMap["AudioCodec"].(string)
		}

		if sessionMap["VideoCodec"] != nil {
			videoCodec = sessionMap["VideoCodec"].(string)
		}

		if sessionMap["Container"] != nil {
			container = sessionMap["Container"].(string)
		}

		if sessionMap["AudioChannels"] != nil {
			audioChannels = int(sessionMap["AudioChannels"].(float64))
		}

		if playMethod == "Transcode" && sessionMap["TranscodeReasons"] != nil {
			transcodeReasonsArr := sessionMap["TranscodeReasons"].([]interface{})
			for _, reason := range transcodeReasonsArr {
				transcodeReasons += reason.(string) + ","
			}
			transcodeReasons = strings.TrimRight(transcodeReasons, ",")
		}

		if playMethod != "" {
			// Atualizar a métrica com os dados obtidos
			streamCount.With(prometheus.Labels{
				"user_id":           sessionMap["UserId"].(string),
				"username":          sessionMap["UserName"].(string),
				"client":            sessionMap["Client"].(string),
				"device":            sessionMap["DeviceName"].(string),
				"play_method":       playMethod,
				"audio_codec":       audioCodec,
				"video_codec":       videoCodec,
				"container":         container,
				"audio_channels":    strconv.Itoa(audioChannels),
				"transcode_reasons": transcodeReasons,
			}).Set(1)
		}
	}
}
