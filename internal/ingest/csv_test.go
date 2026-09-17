package ingest

import (
    "strings"
    "testing"
)

func TestParseSensorCSV_ValidInput(t *testing.T) {
    csvData := `timestamp,humedad_pct,temp_c,parcela_id
2026-03-15T06:00:00Z,23.5,18.2,PARCELA_01
2026-03-15T12:00:00Z,19.8,26.1,PARCELA_01`

    readings, err := ParseSensorCSV(strings.NewReader(csvData))
    if err != nil {
        t.Fatalf("Se esperaba nil error, se obtuvo: %v", err)
    }

    if len(readings) != 2 {
        t.Fatalf("Se esperaban 2 lecturas, se obtuvieron %d", len(readings))
    }

    if readings[0].HumedadPct != 23.5 {
        t.Errorf("Se esperaba humedad 23.5, se obtuvo %f", readings[0].HumedadPct)
    }
}

