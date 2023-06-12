package utils

import (
	"log"
	"time"
)

// CalculateTaskDuration recebe strings de data/hora de início e fim no formato RFC3339 e retorna a duração em segundos.
func CalculateTaskDuration(startTimeUtc, endTimeUtc string) float64 {
	const layout = time.RFC3339Nano // Layout que corresponde ao formato "2023-06-10T08:23:30.5098827Z"

	// Convertendo strings para objetos de tempo
	startTime, err := time.Parse(layout, startTimeUtc)
	if err != nil {
		log.Println("Erro ao converter a hora de início:", err)
		return 0
	}

	endTime, err := time.Parse(layout, endTimeUtc)
	if err != nil {
		log.Println("Erro ao converter a hora de término:", err)
		return 0
	}

	// Calculando a duração e retornando em segundos
	duration := endTime.Sub(startTime).Seconds()
	return duration
}
