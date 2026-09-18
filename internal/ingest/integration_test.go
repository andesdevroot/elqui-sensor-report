package ingest

import (
	"os"
	"path/filepath"
	"testing"
)

// fixturePath apunta al CSV de ejemplo en testdata/ de la raíz del repositorio.
var fixturePath = filepath.Join("..", "..", "testdata", "sensor_sample.csv")

// TestParseSensorCSV_IntegracionFixture ejercita el pipeline real archivo → models,
// sin mocks: abre el CSV de testdata/, lo parsea y valida las lecturas resultantes.
func TestParseSensorCSV_IntegracionFixture(t *testing.T) {
	f, err := os.Open(fixturePath)
	if err != nil {
		t.Fatalf("no se pudo abrir el fixture %q: %v", fixturePath, err)
	}
	defer f.Close()

	readings, err := ParseSensorCSV(f)
	if err != nil {
		t.Fatalf("ParseSensorCSV devolvió error: %v", err)
	}

	if len(readings) != 6 {
		t.Fatalf("Se esperaban 6 lecturas, se obtuvieron %d", len(readings))
	}

	primera := readings[0]
	if primera.ParcelaID != "PARCELA_01" {
		t.Errorf("Se esperaba parcela PARCELA_01, se obtuvo %q", primera.ParcelaID)
	}
	if primera.HumedadPct != 23.5 {
		t.Errorf("Se esperaba humedad 23.5, se obtuvo %f", primera.HumedadPct)
	}
	if primera.TempC != 18.2 {
		t.Errorf("Se esperaba temperatura 18.2, se obtuvo %f", primera.TempC)
	}
	if got := primera.Timestamp.UTC().Format("2006-01-02T15:04:05Z"); got != "2026-03-15T06:00:00Z" {
		t.Errorf("Se esperaba timestamp 2026-03-15T06:00:00Z, se obtuvo %q", got)
	}

	var parcela02 int
	for _, r := range readings {
		if r.ParcelaID == "PARCELA_02" {
			parcela02++
		}
	}
	if parcela02 != 3 {
		t.Errorf("Se esperaban 3 lecturas de PARCELA_02, se obtuvieron %d", parcela02)
	}
}
