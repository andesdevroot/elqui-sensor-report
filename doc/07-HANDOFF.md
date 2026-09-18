# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- Fase 1 — Ingesta: **T1.1 y T1.3 cerradas**; T1.2 bloqueada por falta de documentación de la fuente DGA.
- Suite: 5 tests en verde (`go test ./...`), incluido el test de integración sobre `testdata/sensor_sample.csv`.
- Repositorio publicado en https://github.com/andesdevroot/elqui-sensor-report (rama `main`).

## Último commit

- `f87995c` — `test(ingest): agrega test de integración con fixture CSV de sensores`

## Siguiente tarea

- **T2.1 — ET0 Hargreaves-Samani**: definir `ET0Result` en `internal/models` (contrato en `doc/02`) e implementar el cálculo con TDD. Falta fijar la fuente de validación de los casos conocidos (TODO). El sensor entrega `TempC`, insumo principal de Hargreaves-Samani.
- T1.2 (cliente DGA) sigue bloqueada.

## Blockers

- **T1.2 bloqueada por falta de información de la fuente**: no hay endpoint ni formato de descarga documentado para la DGA (TODO en `doc/05-DATA-SOURCES.md`).
- Sin blockers técnicos.

## Decisiones recientes

- **Fixture en la raíz** (`testdata/sensor_sample.csv`, no por paquete): sirve como dataset de ejemplo además de insumo del test, y el README ya declara ese directorio.
- **Test de integración sin mocks**: `internal/ingest/integration_test.go` abre el CSV por ruta relativa y valida 6 lecturas (2 parcelas × 3), incluidos `Timestamp`, `HumedadPct` y `TempC`.
- **T1.3 cerrada**: Fase 1 queda completa salvo T1.2.
