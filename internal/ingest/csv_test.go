package ingest

import (
	"errors"
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

func TestParseSensorCSV_HeaderInvalido(t *testing.T) {
	csvData := `timestamp,humedad,temp_c,parcela_id
2026-03-15T06:00:00Z,23.5,18.2,PARCELA_01`

	_, err := ParseSensorCSV(strings.NewReader(csvData))
	if !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("Se esperaba ErrInvalidHeader, se obtuvo: %v", err)
	}
}

func TestParseSensorCSV_CSVVacio(t *testing.T) {
	_, err := ParseSensorCSV(strings.NewReader(""))
	if err == nil {
		t.Fatal("Se esperaba error para un CSV vacío, se obtuvo nil")
	}
}

func TestParseSensorCSV_FilaInvalida(t *testing.T) {
	tests := []struct {
		name string
		row  string
	}{
		{"humedad sobre 100", "2026-03-15T06:00:00Z,150,18.2,PARCELA_01"},
		{"humedad bajo 0", "2026-03-15T06:00:00Z,-5,18.2,PARCELA_01"},
		{"temperatura sobre 50", "2026-03-15T06:00:00Z,23.5,60,PARCELA_01"},
		{"temperatura bajo -10", "2026-03-15T06:00:00Z,23.5,-15,PARCELA_01"},
		{"timestamp no RFC3339", "15-03-2026,23.5,18.2,PARCELA_01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csvData := "timestamp,humedad_pct,temp_c,parcela_id\n" + tt.row

			_, err := ParseSensorCSV(strings.NewReader(csvData))
			if err == nil {
				t.Fatalf("Se esperaba error para %q, se obtuvo nil", tt.row)
			}
		})
	}
}

// TestParseSensorCSV_ExportadoExcel cubre las variantes que introducen los
// exportadores tipo Excel: BOM UTF-8, espacios alrededor del header y CRLF.
func TestParseSensorCSV_ExportadoExcel(t *testing.T) {
	const header = "timestamp,humedad_pct,temp_c,parcela_id"
	const row = "2026-03-15T06:00:00Z,23.5,18.2,PARCELA_01"

	tests := []struct {
		name string
		data string
	}{
		{"BOM UTF-8", "\ufeff" + header + "\n" + row + "\n"},
		{"espacios en header", " timestamp , humedad_pct , temp_c , parcela_id \n" + row + "\n"},
		{"CRLF", header + "\r\n" + row + "\r\n"},
		{"BOM + espacios + CRLF", "\ufeff timestamp , humedad_pct , temp_c , parcela_id \r\n" + row + "\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readings, err := ParseSensorCSV(strings.NewReader(tt.data))
			if err != nil {
				t.Fatalf("Se esperaba nil error, se obtuvo: %v", err)
			}
			if len(readings) != 1 {
				t.Fatalf("Se esperaba 1 lectura, se obtuvo %d", len(readings))
			}
			if readings[0].ParcelaID != "PARCELA_01" {
				t.Errorf("Se esperaba PARCELA_01, se obtuvo %q", readings[0].ParcelaID)
			}
		})
	}
}
