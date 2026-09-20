# Task — cola de trabajo

> Ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md) para el detalle de fases, entregables y criterios de aceptación.

## Estado actual

- Fase 0 — Congelamiento y migración: **cerrada** (motor Go archivado en `internal-go-archive/`, tag `go-motor-archive`; stack Python puro).
- Fase 1 — Pipeline satelital mínimo: **siguiente**. `scripts/ndvi_probe.py` ya existe en el working tree sin commitear; falta `requirements.txt` y la validación en la parcela La Serena `2W89+VG`.
- Retoma detallada (blockers, deuda y decisiones) en [doc/07-HANDOFF.md](doc/07-HANDOFF.md).

## Fase 0 — Congelamiento y migración

- [x] Archivar el motor Go: `internal/` → `internal-go-archive/` (tag `go-motor-archive`)
- [x] Congelar el módulo: `go.mod` → `go.mod.archived` (`.gitignore` excluye `*.archived`)
- [x] Documentación reescrita a Python puro (`doc/00`, `doc/01`, `doc/04`, `doc/06`)
- [x] Eliminar directorios vestigiales de Go (`cmd/` y `pkg/` vacíos)

## Fase 1 — Pipeline satelital mínimo

- [ ] T1.1: `scripts/ndvi_probe.py` — Earth Search (anónimo) + lectura por ventana con `rasterio` → NDVI del AOI ← siguiente (el archivo ya existe, sin commitear)
- [ ] T1.2: `scripts/requirements.txt` (`pystac-client`, `rasterio`, `numpy`)
- [ ] T1.3: validación documentada con la parcela La Serena `2W89+VG`

## Fase 2 — Motor agronómico en Python

- [ ] T2.1: `scripts/et0.py` — Hargreaves-Samani portado del Go archivado, validado contra FAO-56
- [ ] T2.2: `scripts/kcb.py` — NDVI → Kcb (método INIA: `Kcb = 1.51 × NDVI − 0.23`)
- [ ] T2.3: `scripts/irrigation.py` — ET0 × Kcb → mm/día, con tests de `pytest`

## Fase 3 — Capa DeepSeek

- [ ] T3.1: `scripts/deepseek.py` — cliente con timeout y manejo de errores
- [ ] T3.2: template de prompt en español de Chile + fallback si la API no responde

## Fase 4 — Serie temporal multitemporal

- [ ] T4.1: `scripts/timeseries.py` — 2019–2026, 10 parcelas de Coquimbo, salida Parquet

## Fase 5 — Baseline FAO-56 + NDVI (método PLAS/INIA)

- [ ] T5.1: baseline agronómico sobre la serie temporal + métricas de referencia

## Fase 6 — Modelo deep learning

- [ ] T6.1: LSTM o Transformer temporal con partición temporal, comparado contra el baseline

## Fase 7 — Web mínima (FastAPI + HTMX + Leaflet)

- [ ] T7.1: `web/` — dibujar el polígono y mostrar NDVI + recomendación + explicación

## Fase 8 — Paper y publicación

- [ ] T8.1: paper (método, dataset, baseline, modelo, resultados) + release open source

## Histórico Go

El motor anterior (ET0 Hargreaves-Samani + parser CSV) está archivado y **no se mantiene**:

- Código: [`internal-go-archive/`](internal-go-archive/) (`analysis`, `ingest`, `models`) y `go.mod.archived`.
- Tag: `go-motor-archive`.
- Motivo: Python puro por expertise del autor y ecosistema científico (rasters, series temporales, deep learning).
- Uso: los casos FAO-56 de `doc/05-DATA-SOURCES.md` son la especificación de `scripts/et0.py`; el parser CSV queda como referencia.
- También histórico: `archive/copernicus/` (cliente Copernicus, tag `copernicus-archive`).
