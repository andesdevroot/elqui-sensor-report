// Package analysis implementa los cálculos agronómicos del proyecto:
// evapotranspiración de referencia (ET0), balance hídrico y eficiencia.
package analysis

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/andesdevroot/elqui-sensor-report/internal/models"
)

// et0Method identifica el método de cálculo registrado en models.ET0Result.
const et0Method = "Hargreaves-Samani (FAO-56)"

// ErrInvalidInput se devuelve cuando los datos de entrada no son utilizables.
var ErrInvalidInput = errors.New("entrada inválida")

// ET0Hargreaves calcula la evapotranspiración de referencia diaria (mm/día) con el
// método Hargreaves-Samani (FAO-56, Allen et al. 1998, ec. 52):
//
//	ET0 = 0.0023 (Tmedia + 17.8) (Tmax - Tmin)^0.5 Ra
//
// donde Tmedia = (Tmax + Tmin) / 2 y Ra (mm/día) es la radiación extraterrestre
// calculada con las ec. 21-25 de FAO-56 a partir de la latitud y el día del año.
func ET0Hargreaves(date time.Time, latDeg, tMaxC, tMinC float64) (models.ET0Result, error) {
	if latDeg < -90 || latDeg > 90 {
		return models.ET0Result{}, fmt.Errorf("%w: latitud %g fuera de [-90, 90]", ErrInvalidInput, latDeg)
	}
	if tMaxC < tMinC {
		return models.ET0Result{}, fmt.Errorf("%w: Tmax (%g) menor que Tmin (%g)", ErrInvalidInput, tMaxC, tMinC)
	}

	tMean := (tMaxC + tMinC) / 2
	et0 := 0.0023 * (tMean + 17.8) * math.Sqrt(tMaxC-tMinC) * raMM(date, latDeg)

	return models.ET0Result{
		Fecha:  date,
		ET0mm:  et0,
		Metodo: et0Method,
		// Confianza = 1.0: sin incertidumbre cuantificada (ver doc/02-DATA-MODEL.md).
		Confianza: 1.0,
	}, nil
}

// raMM calcula la radiación extraterrestre Ra en mm/día (FAO-56, ec. 21-25).
// El valor en MJ m-2 día-1 se recupera multiplicando el resultado por 2.45.
func raMM(date time.Time, latDeg float64) float64 {
	j := float64(date.YearDay())
	phi := latDeg * math.Pi / 180

	// Distancia relativa Tierra-Sol (ec. 23) y declinación solar (ec. 24).
	dr := 1 + 0.033*math.Cos(2*math.Pi/365*j)
	delta := 0.409 * math.Sin(2*math.Pi/365*j-1.39)

	// Ángulo horario de puesta de sol (ec. 25); se acota a [-1, 1] para evitar NaN.
	ws := math.Acos(math.Max(-1, math.Min(1, -math.Tan(phi)*math.Tan(delta))))

	// Radiación extraterrestre (ec. 21) → evaporación equivalente (λ = 2.45 MJ/kg).
	ra := 24 * 60 / math.Pi * 0.0820 * dr * (ws*math.Sin(phi)*math.Sin(delta) + math.Cos(phi)*math.Cos(delta)*math.Sin(ws))
	return ra / 2.45
}
