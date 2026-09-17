# Task — cola de trabajo

> Ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md) para el detalle de fases y criterios de aceptación.

## Estado actual

- Fase 1 — Ingesta, en curso. **T1.1 completada**; siguiente: **T1.2 (cliente DGA)**.
- Retoma detallada (blockers y decisiones) en [doc/07-HANDOFF.md](doc/07-HANDOFF.md).

## Fase 0 — Setup

- [x] Toolchain de Go instalada (`go 1.23.0`)
- [x] Repositorio git inicializado (rama `main`)
- [x] `doc/` creada como fuente de verdad

## Fase 1 — Ingesta

- [x] T1.1: parser CSV de sensores — criterio cumplido: `go test ./internal/ingest/...` en verde (`internal/models/sensor.go`, `internal/ingest/csv.go`)
- [ ] T1.2: cliente DGA — criterio: TODO (se define al documentar la fuente) ← siguiente
- [ ] T1.3: tests de integración CSV → models — criterio: pipeline de punta a punta sin mocks

## Fase 2 — Análisis

- [ ] T2.1: ET0 Hargreaves-Samani — criterio: tests con casos conocidos (TODO: fuente de validación)
- [ ] T2.2: balance hídrico — criterio: consumo real vs. óptimo por período con tests de tabla
- [ ] T2.3: eficiencia y semáforo — criterio: `EficienciaPct` + `Semaforo` por umbrales documentados

## Fase 3 — Reporte

- [ ] T3.1: reporte Markdown con semáforo — criterio: salida `.md` legible generada desde un `EfficiencyReport`
- [ ] T3.2: recomendaciones es-CL — criterio: sección de recomendaciones en español de Chile

## Fase 4 — Release v0.1.0

- [ ] T4.1: README final + `install.sh` + LICENSE + tag `v0.1.0` — criterio: build limpio + tag publicado
