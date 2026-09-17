# 06 — Roadmap

Cada tarea se cierra con un commit atómico y su criterio de aceptación cumplido (ver `doc/03-METHODOLOGY.md`).

## Fase 0 — Setup (en progreso ✅)

- Toolchain de Go instalada (`go 1.23.0`).
- Repositorio git inicializado (rama `master`, aún sin commits ni remoto — TODO: crear repositorio remoto).
- `doc/` creada como fuente de verdad (este commit).

## Fase 1 — Ingesta

**T1.1 — Parser CSV de sensores**
Criterio de aceptación: `go test ./internal/ingest/...` en verde; `ParseSensorCSV` convierte el CSV del sensor en `[]models.SensorReading` y rechaza lecturas que violan los invariantes de `doc/02-DATA-MODEL.md`.

**T1.2 — Cliente DGA**
Criterio de aceptación: TODO — se define al documentar la fuente en `doc/05-DATA-SOURCES.md`.

**T1.3 — Tests de integración**
Criterio de aceptación: pipeline `CSV → models` de punta a punta sobre un archivo de `testdata/`, sin mocks.

## Fase 2 — Análisis

**T2.1 — ET0 Hargreaves-Samani**
Criterio de aceptación: `ET0Result` calculado desde las temperaturas del sensor; tests con casos conocidos (TODO: definir fuente de validación).

**T2.2 — Balance hídrico**
Criterio de aceptación: consumo real vs. óptimo por período, con tests de tabla.

**T2.3 — Eficiencia y semáforo**
Criterio de aceptación: `EficienciaPct` y `Semaforo` calculados según los umbrales que se definan en `doc/02-DATA-MODEL.md` (TODO).

## Fase 3 — Reporte

**T3.1 — Reporte Markdown con semáforo**
Criterio de aceptación: la CLI genera un `.md` legible a partir de un `EfficiencyReport`.

**T3.2 — Recomendaciones en es-CL**
Criterio de aceptación: sección de recomendaciones en español de Chile, derivada del semáforo.

## Fase 4 — Release v0.1.0

**T4.1 — Release**
Entregables: README final, `install.sh`, LICENSE, tag `v0.1.0`.
Criterio de aceptación: `go build` limpio en una máquina limpia y tag publicado con notas de release.
