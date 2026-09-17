# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- Fase 1 — Ingesta, en curso. **T1.1 cerrada**: parser CSV con caminos de error cubiertos por tests (4 tests, 5 subtests de tabla).
- `go test ./...` pasa. Repositorio publicado en https://github.com/andesdevroot/elqui-sensor-report (rama `main`).

## Último commit

- `703573b` — `docs: indexa 07-HANDOFF y documenta validación en ingest`

## Siguiente tarea

- **T1.3 — Tests de integración CSV → models**: pipeline de punta a punta sobre un archivo de `testdata/`, sin mocks.
- T1.2 (cliente DGA) queda bloqueada por falta de documentación de la fuente.

## Blockers

- **T1.2 bloqueada por falta de información de la fuente**: no hay endpoint ni formato de descarga documentado para la DGA (TODO en `doc/05-DATA-SOURCES.md`).
- Sin blockers técnicos.

## Decisiones recientes

- **Tests de caminos de error**: `internal/ingest/csv_test.go` cubre header inválido (`errors.Is` sobre `ErrInvalidHeader`), CSV vacío, humedad fuera de [0,100] (150 y -5), temperatura fuera de [-10,50] (60 y -15) y timestamp no RFC3339. No fue necesario tocar `csv.go`.
- **`doc/02` actualizado**: el TODO sobre *dónde* se validan los invariantes quedó resuelto con una referencia a `internal/ingest/csv.go`.
- **README indexa `doc/07`**: la lista de documentación ahora cubre 00–07.
- **gofmt**: `internal/ingest/csv_test.go` quedó formateado; `gofmt -l internal/` está limpio.
