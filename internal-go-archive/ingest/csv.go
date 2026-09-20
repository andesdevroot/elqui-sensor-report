// Package ingest ingesta datos desde fuentes externas (CSV de sensores y, más
// adelante, APIs públicas) y los convierte en estructuras del dominio.
package ingest

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/andesdevroot/elqui-sensor-report/internal/models"
)

// csvHeader es el encabezado exacto que debe declarar el CSV de sensores.
const csvHeader = "timestamp,humedad_pct,temp_c,parcela_id"

// numColumns es el número de columnas esperadas en cada fila del CSV.
const numColumns = 4

// ErrInvalidHeader se devuelve cuando el encabezado del CSV no coincide con csvHeader.
var ErrInvalidHeader = errors.New("header del CSV inválido")

// ParseSensorCSV convierte un CSV de sensor de humedad de suelo en una lista de
// models.SensorReading. Tolera CSV exportados desde Excel (BOM UTF-8 y espacios
// alrededor del header), valida el header exacto y los invariantes de cada
// lectura: HumedadPct en [0,100], TempC en [-10,50], timestamp en RFC3339 y
// parcela_id no vacío.
func ParseSensorCSV(r io.Reader) ([]models.SensorReading, error) {
	reader := csv.NewReader(r)

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("parsear header: CSV vacío: %w", err)
		}
		return nil, fmt.Errorf("parsear header: %w", err)
	}
	header = normalizeHeader(header)
	if got := strings.Join(header, ","); got != csvHeader {
		return nil, fmt.Errorf("%w: se esperaba %q, se obtuvo %q", ErrInvalidHeader, csvHeader, got)
	}

	var readings []models.SensorReading
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parsear línea %d: %w", len(readings)+2, err)
		}

		reading, err := parseRecord(record, len(readings)+2)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}

	return readings, nil
}

// normalizeHeader limpia el header de los artefactos típicos de los exportadores
// tipo Excel: BOM UTF-8 al inicio del archivo y espacios alrededor de cada columna.
// El CRLF ya lo normaliza encoding/csv. Devuelve el mismo slice modificado in-place.
func normalizeHeader(header []string) []string {
	if len(header) == 0 {
		return header
	}
	header[0] = strings.TrimPrefix(header[0], "\ufeff")
	for i, col := range header {
		header[i] = strings.TrimSpace(col)
	}
	return header
}

// parseRecord convierte una fila del CSV (sin el header) en un models.SensorReading validado.
func parseRecord(record []string, line int) (models.SensorReading, error) {
	if len(record) != numColumns {
		return models.SensorReading{}, fmt.Errorf("línea %d: se esperaban %d columnas, se obtuvieron %d", line, numColumns, len(record))
	}

	ts, err := time.Parse(time.RFC3339, record[0])
	if err != nil {
		return models.SensorReading{}, fmt.Errorf("línea %d: timestamp inválido %q: %w", line, record[0], err)
	}

	humedad, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return models.SensorReading{}, fmt.Errorf("línea %d: humedad_pct inválido %q: %w", line, record[1], err)
	}
	if humedad < 0 || humedad > 100 {
		return models.SensorReading{}, fmt.Errorf("línea %d: humedad_pct fuera de rango [0,100]: %g", line, humedad)
	}

	temp, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return models.SensorReading{}, fmt.Errorf("línea %d: temp_c inválido %q: %w", line, record[2], err)
	}
	if temp < -10 || temp > 50 {
		return models.SensorReading{}, fmt.Errorf("línea %d: temp_c fuera de rango [-10,50]: %g", line, temp)
	}

	if record[3] == "" {
		return models.SensorReading{}, fmt.Errorf("línea %d: parcela_id vacío", line)
	}

	return models.SensorReading{
		Timestamp:  ts,
		HumedadPct: humedad,
		TempC:      temp,
		ParcelaID:  record[3],
	}, nil
}
