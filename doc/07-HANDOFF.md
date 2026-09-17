# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- Fase 1 — Ingesta, en curso.
- T1.1 (parser CSV de sensores) completada y en verde: `internal/models/sensor.go` + `internal/ingest/csv.go`.
- `go test ./...` pasa. Repositorio publicado en https://github.com/andesdevroot/elqui-sensor-report (rama `main`).

## Último commit

- `b108196` — `feat(ingest): implementa parser CSV de sensores con validación de header`

## Siguiente tarea

- **T1.2 — Cliente DGA**. Criterio de aceptación aún sin definir: depende de documentar endpoint/formato de descarga en `doc/05-DATA-SOURCES.md` (TODO abierto).
- Alternativa si T1.2 sigue bloqueada: **T1.3** (tests de integración CSV → models sobre un archivo de `testdata/`) o ampliar los tests del parser a los caminos de error (header inválido, rangos, filas malformadas).

## Blockers

- **T1.2 bloqueada por falta de información de la fuente**: no hay endpoint ni formato de descarga documentado para la DGA. Requiere investigación (TODO en `doc/05-DATA-SOURCES.md`).
- Sin blockers técnicos: `go 1.23.0` operativo y `gh` autenticado como `andesdevroot`.

## Decisiones recientes

- **Validación de invariantes en `ingest`**: `ParseSensorCSV` rechaza humedad fuera de [0,100], temperatura fuera de [-10,50], `parcela_id` vacío, timestamp que no sea RFC3339 y filas con distinto número de columnas. Esto resuelve el TODO de `doc/02-DATA-MODEL.md` sobre *dónde* se validan los invariantes (decisión: en la capa de ingesta, devolviendo error).
- **Header exacto**: el CSV debe declarar exactamente `timestamp,humedad_pct,temp_c,parcela_id`; cualquier variación devuelve `ErrInvalidHeader`.
- **Sin dependencias externas**: el parser usa solo stdlib (`encoding/csv`, `errors`, `fmt`, `io`, `strconv`, `strings`, `time`).
