package metrics

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/api"
	"github.com/rafaelvieiras/jellyfin-exporter/pkg/utils"
)

type ScheduledTask struct {
	Id              string  `json:"Id"`
	Name            string  `json:"Name"`
	State           string  `json:"State"`
	LastExecutionMs float64 `json:"LastExecutionResultDurationMs"`
}

var (
	scheduledTasks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_scheduled_tasks",
			Help: "Information about Jellyfin scheduled tasks",
		},
		[]string{"id", "name", "state"},
	)
)

func FetchScheduledTasks(jellyfinApiUrl, jellyfinToken string) {
	scheduledTasksUrl := fmt.Sprintf("%s/ScheduledTasks", jellyfinApiUrl)

	rawData, err := api.MakeRequest(scheduledTasksUrl, jellyfinToken)
	if err != nil {
		log.Println("Error fetching scheduled tasks:", err)
		return
	}

	data, ok := rawData.([]interface{})
	if !ok {
		log.Println("Error parsing scheduled tasks response:", err)
		return
	}
	scheduledTasks.Reset()

	for _, task := range data {
		taskMap, ok := task.(map[string]interface{})
		if !ok {
			log.Println("Error parsing task as map")
			continue
		}

		if taskMap["LastExecutionResult"] != nil {
			lastExecutionMap, ok := taskMap["LastExecutionResult"].(map[string]interface{})
			if !ok {
				log.Println("Error parsing LastExecutionResult as map")
				continue
			}

			startTime, ok := lastExecutionMap["StartTimeUtc"].(string)
			if !ok {
				log.Println("Error parsing StartTimeUtc as string or StartTime might not exist")
				continue
			}

			endTime, ok := lastExecutionMap["StartTimeUtc"].(string)
			if !ok {
				log.Println("Error parsing StartTimeUtc as string or EndTime might not exist")
				continue
			}

			id, ok := taskMap["Id"].(string)
			if !ok {
				log.Println("Error parsing Id as string")
				continue
			}

			name, ok := taskMap["Name"].(string)
			if !ok {
				log.Println("Error parsing Name as string")
				continue
			}

			state, ok := taskMap["State"].(string)
			if !ok {
				log.Println("Error parsing State as string")
				continue
			}

			lastExecutionMs := utils.CalculateTaskDuration(startTime, endTime)
			scheduledTasks.WithLabelValues(id, name, state).Set(lastExecutionMs)
		}
	}

}
