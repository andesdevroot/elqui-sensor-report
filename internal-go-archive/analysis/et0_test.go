package analysis

import (
	"math"
	"testing"
	"time"
)

// TestET0Hargreaves_CasosFAO56 usa los casos de referencia documentados en
// doc/05-DATA-SOURCES.md (FAO-56, ec. 52).
func TestET0Hargreaves_CasosFAO56(t *testing.T) {
	cases := []struct {
		name string
		date time.Time
		lat  float64
		tMax float64
		tMin float64
		want float64
		tol  float64
	}{
		{
			"FAO-56 Ej.20 (Lyon, 45.7167N, 15-jul, 26.6/14.8)",
			time.Date(2026, time.July, 15, 0, 0, 0, 0, time.UTC),
			45.7167, 26.6, 14.8, 5.0, 0.1,
		},
		{
			"Elqui verano (30S, 15-ene, 30.0/15.0)",
			time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
			-30.0, 30.0, 15.0, 6.32, 0.01,
		},
		{
			"Elqui invierno (30S, 15-jul, 18.0/6.0)",
			time.Date(2026, time.July, 15, 0, 0, 0, 0, time.UTC),
			-30.0, 18.0, 6.0, 1.90, 0.01,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ET0Hargreaves(c.date, c.lat, c.tMax, c.tMin)
			if err != nil {
				t.Fatalf("ET0Hargreaves devolvió error: %v", err)
			}
			if math.Abs(got.ET0mm-c.want) > c.tol {
				t.Errorf("ET0 = %.3f mm/día, se esperaba %.2f ± %.2f", got.ET0mm, c.want, c.tol)
			}
			if got.Metodo == "" {
				t.Error("Metodo no debe estar vacío")
			}
		})
	}
}

// TestRaMM_FAO56Ejemplo8 valida el término radiativo con el Ejemplo 8 de FAO-56:
// 3 de septiembre (J = 246) a 20°S → Ra = 32.2 MJ m-2 d-1 = 13.1 mm/día.
func TestRaMM_FAO56Ejemplo8(t *testing.T) {
	date := time.Date(2026, time.September, 3, 0, 0, 0, 0, time.UTC)

	got := raMM(date, -20.0)
	if math.Abs(got-13.1) > 0.05 {
		t.Errorf("Ra = %.3f mm/día, se esperaba 13.1 ± 0.05 (FAO-56, ec. 21-25)", got)
	}
}

func TestET0Hargreaves_EntradaInvalida(t *testing.T) {
	cases := []struct {
		name string
		lat  float64
		tMax float64
		tMin float64
	}{
		{"Tmax menor que Tmin", -30.0, 10.0, 20.0},
		{"latitud fuera de rango", 91.0, 20.0, 10.0},
	}

	date := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := ET0Hargreaves(date, c.lat, c.tMax, c.tMin)
			if err == nil {
				t.Fatal("Se esperaba error de entrada inválida, se obtuvo nil")
			}
		})
	}
}
