# Task — cola de trabajo

> Ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md) para el detalle de fases y criterios de aceptación.

**Estado actual**: Fase 0 — Setup (en progreso).

## Fase 0 — Setup

- [x] Toolchain de Go instalada (`go 1.23.0`)
- [x] Repositorio git inicializado
- [ ] `doc/` creada como fuente de verdad — en progreso (este commit)

## Fase 1 — Ingesta

- [ ] T1.1: parser CSV de sensores (test RED ya escrito en `internal/ingest/csv_test.go`) — criterio: `go test ./internal/ingest/...` en verde contra `testdata/`
- [ ] T1.2: cliente DGA — criterio: TODO (se define al documentar la fuente)
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
