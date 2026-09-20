# Task — cola de trabajo

> Ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md) para el detalle de fases, entregables y criterios de aceptación.

## Estado actual

- Fase 0 — Reset y limpieza: **cerrada** (`pkg/copernicus` archivado con el tag `copernicus-archive`; Earth Search como fuente satelital).
- Fase 1 — Pipeline satelital mínimo: **siguiente**. Arranca por `scripts/ndvi_probe.py`.
- Retoma detallada (blockers, deuda técnica y decisiones) en [doc/07-HANDOFF.md](doc/07-HANDOFF.md).

## Fase 0 — Reset y limpieza

- [x] Archivar `pkg/copernicus` en `archive/copernicus/` (tag `copernicus-archive`)
- [x] Documentación reencuadrada al copiloto de riego (`doc/00`, `doc/01`, `doc/05`, `doc/06`)
- [x] Eliminar directorios vestigiales (`internal/satellite`, `internal/ai`, `cmd/elqui`)

## Fase 1 — Pipeline satelital mínimo

- [ ] T1.1: `scripts/ndvi_probe.py` — búsqueda STAC en **Earth Search** (anónimo) + lectura **por ventana** del AOI con `rasterio` → NDVI medio del predio ← siguiente
- [ ] T1.2: validación sobre la parcela La Serena `2W89+VG` (POINT `-71.24894 -29.90453`)
- [ ] T1.3: `scripts/requirements.txt` (`pystac-client`, `rasterio`, `numpy`)

## Fase 2 — Integración NDVI → Kcb → ET0 (Go)

- [ ] T2.1: `scripts/kcb.py` — NDVI → Kcb con el método INIA (`Kcb = 1.51 × NDVI − 0.23`)
- [ ] T2.2: `internal/analysis` — recomendación `ET0 × Kcb` en mm/día, con tests de tabla

## Fase 3 — Capa DeepSeek

- [ ] T3.1: cliente DeepSeek en Go (`net/http`, timeouts, manejo de errores)
- [ ] T3.2: prompt de traducción a lenguaje natural (es-CL) con fallback si la API no responde

## Fase 4 — Web mínima (Go + HTMX + Leaflet)

- [ ] T4.1: `cmd/elqui-web` + `web/`
- [ ] T4.2: dibujar el polígono de la parcela y mostrar NDVI + recomendación

## Fase 5 — Validación con agricultores y publicación

- [ ] T5.1: validación con al menos un agricultor o asesor de Coquimbo
- [ ] T5.2: contraste con PLAS (INIA) y RiegaBien (UC) + publicación

## Motor existente (se mantiene)

- [x] Parser CSV de sensores (`internal/ingest`) — fallback
- [x] ET0 Hargreaves-Samani validada contra FAO-56 (`internal/analysis`)
- [x] `archive/copernicus/` — cliente Copernicus archivado; referencia histórica, sin mantenimiento
