// Package models define las estructuras de datos compartidas del proyecto
// (lecturas de sensores, resultados de ET0 y reportes de eficiencia).
package models

import "time"

// SensorReading representa una lectura puntual de un sensor de humedad de suelo.
type SensorReading struct {
	Timestamp  time.Time `json:"timestamp"`
	HumedadPct float64   `json:"humedad_pct"`
	TempC      float64   `json:"temp_c"`
	ParcelaID  string    `json:"parcela_id"`
}
