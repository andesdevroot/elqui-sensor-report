# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- Fase 1 — Ingesta: T1.1, T1.3 y T1.4 cerradas; T1.2 bloqueada por falta de documentación de la fuente DGA.
- Fase 2 — Análisis: T2.1 cerrada (ET0 Hargreaves-Samani validada contra FAO-56); siguiente T2.2.
- Suite en verde (`go test ./...`). Repositorio publicado en https://github.com/andesdevroot/elqui-sensor-report (rama `main`).

## Último commit

- `bb32ecd` — `docs(conventions): marca fixtures de testdata como inmutables`

## Siguiente tarea

- **T2.2 — Balance hídrico**: consumo real vs. óptimo por período, con tests de tabla (criterio en `doc/06-ROADMAP.md`).
- T1.2 (cliente DGA) sigue bloqueada.

## Blockers

- **T1.2 bloqueada por falta de información de la fuente**: no hay endpoint ni formato de descarga documentado para la DGA (TODO en `doc/05-DATA-SOURCES.md`).
- Sin blockers técnicos.

## Decisiones recientes

- **Fixtures de `testdata/` inmutables** (`doc/04-CONVENTIONS.md`): si se necesita otro caso se crea un archivo nuevo. Regla añadida tras detectar un drift local en `sensor_sample.csv` (un espacio antes del header) que git no puede atribuir porque nunca se commiteó.
- **T1.4 — tolerancia a CSV de Excel**: `ParseSensorCSV` ahora limpia BOM UTF-8 y espacios alrededor del header (`normalizeHeader`). El CRLF ya lo resolvía `encoding/csv`, así que no requirió código.
- **Fixture nuevo**: `testdata/sensor_sample_excel.csv` (BOM + espacios en header + CRLF) con su test de integración.
