package models

import "time"

// ET0Result contiene la ET0 (mm) calculada para una fecha y el método usado.
type ET0Result struct {
	Fecha     time.Time `json:"fecha"`
	ET0mm     float64   `json:"et0_mm"`
	Metodo    string    `json:"metodo"`
	Confianza float64   `json:"confianza"`
}
