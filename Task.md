# Task — cola de trabajo

> Ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md) para el detalle de fases, entregables y criterios de aceptación.

## Estado actual

- Fase 0 — Refactor de identidad: **cerrada** (`elqui-eye`; el motor de sensores pasa a ser fallback).
- Fase 1 — Pipeline satelital: **en progreso**. T1.1 (auth Copernicus) y T1.2 (búsqueda STAC) cerradas y validadas en vivo; siguiente: **T1.3 (descarga de bandas)**.
- Retoma detallada (blockers, deuda técnica y decisiones) en [doc/07-HANDOFF.md](doc/07-HANDOFF.md).

## Fase 0 — Refactor de identidad

- [x] Documentación reencuadrada a `elqui-eye` (`README.md`, `doc/00`, `doc/01`, `doc/05`, `doc/06`)
- [x] Estructura base creada (`pkg/copernicus`, `internal/satellite`, `internal/ai`, `web`)

## Fase 1 — Pipeline satelital

- [x] T1.1: `pkg/copernicus` — autenticación OAuth2 client credentials con Copernicus Data Space — criterio cumplido: auth validado contra CDSE real, token JWT OK (`pkg/copernicus/client.go`)
- [x] T1.2: `pkg/copernicus` — búsqueda STAC por AOI (polígono WKT) con filtro de nubosidad — criterio cumplido: validado en vivo, 15 items `cloud_cover < 20` para la parcela La Serena `2W89+VG`; imagen más limpia 0.19 % de nubes (2026-09-01) (`pkg/copernicus/search.go`)
- [ ] T1.3: descarga de bandas B04 (rojo), B08 (NIR) y SCL (máscara de nubes) para el AOI ← siguiente
- [ ] T1.4: `internal/satellite` — NDVI y conversión NDVI → Kcb
- [ ] T1.5: exportación de un PNG de la parcela (NDVI)

## Fase 2 — Recomendación ET0 × Kcb

- [ ] T2.1: integrar ET0 (motor existente) con Kcb satelital
- [ ] T2.2: recomendación diaria en mm/día, con tests de tabla
- [ ] T2.3: CLI `cmd/elqui` — comando de recomendación y build del binario

## Fase 3 — Capa DeepSeek

- [ ] T3.1: `internal/ai` — cliente DeepSeek con `net/http`, timeouts y manejo de errores
- [ ] T3.2: prompt de traducción a lenguaje natural (es-CL) y fallback si la API no responde

## Fase 4 — Web mínima

- [ ] T4.1: `cmd/elqui-web` + `web/` (HTML + HTMX + Leaflet)
- [ ] T4.2: dibujar el polígono de la parcela y mostrar NDVI + recomendación

## Fase 5 — Validación y publicación

- [ ] T5.1: validación con al menos un agricultor/asesor de Coquimbo
- [ ] T5.2: contraste con RiegaBien / PLAS y publicación (LinkedIn, Reddit, universidades)

## Motor de sensores (fallback)

Ya construido y en verde; se usa cuando no hay imagen satelital utilizable.

- [x] Parser CSV de sensores con tolerancia a export de Excel (`internal/ingest`)
- [x] ET0 Hargreaves-Samani validada contra FAO-56 (`internal/analysis`)
