# 06 — Roadmap

Roadmap del proyecto (copiloto de riego + investigación) en 8 fases, más la Fase 0 de migración. Cada fase se cierra con un commit atómico y su criterio de aceptación cumplido (ver `doc/03-METHODOLOGY.md`).

## Fase 0 — Congelamiento y migración (completa ✅)

- **Tareas**: archivar el motor Go (`internal-go-archive/`, tag `go-motor-archive`); congelar `go.mod` como `go.mod.archived`; reescribir la documentación al stack Python puro.
- **Entregables**: repositorio sin código activo en Go; documentación coherente.
- **Criterios de aceptación**: el Go queda solo como referencia etiquetada y `doc/` describe un proyecto Python puro.
- **Tiempo estimado**: 1 sesión.

## Fase 1 — Pipeline satelital mínimo

- **Tareas**: `scripts/ndvi_probe.py` con `pystac-client` + `rasterio` (lectura por ventana del AOI sobre COG); `requirements.txt`; validación con la parcela La Serena `2W89+VG`.
- **Entregables**: `scripts/ndvi_probe.py`, `requirements.txt`, NDVI del AOI para la parcela de prueba.
- **Criterios de aceptación**: dado el polígono y una ventana de fechas, el script devuelve el NDVI medio del AOI **sin credenciales** y sin descargar la escena completa, indicando qué escena eligió.
- **Tiempo estimado**: 1–2 sesiones.

## Fase 2 — Motor agronómico en Python

- **Tareas**: `scripts/et0.py` (Hargreaves-Samani, portado desde el Go archivado y validado contra FAO-56), `scripts/kcb.py` (NDVI → Kcb, método INIA) y `scripts/irrigation.py` (ET0 × Kcb → mm/día).
- **Entregables**: los tres módulos con tests de `pytest`; casos FAO-56 reproducidos.
- **Criterios de aceptación**: los tres casos de referencia de `doc/05-DATA-SOURCES.md` pasan en `pytest` con la tolerancia documentada, y `irrigation.py` entrega la recomendación en mm/día.
- **Tiempo estimado**: 1–2 sesiones.

## Fase 3 — Capa DeepSeek

- **Tareas**: `scripts/deepseek.py` con template de prompt en español de Chile; manejo de errores y timeout; sin datos sensibles.
- **Entregables**: `scripts/deepseek.py`; texto de recomendación en es-CL.
- **Criterios de aceptación**: dada una recomendación numérica se obtiene un texto es-CL, con fallback explícito si la API no responde.
- **Tiempo estimado**: 1 sesión.

## Fase 4 — Serie temporal multitemporal

- **Tareas**: `scripts/timeseries.py` — descarga **2019–2026** para **10 parcelas** de Coquimbo; almacenamiento en Parquet; control de calidad (nubosidad y huecos).
- **Entregables**: dataset multitemporal propio y reproducible.
- **Criterios de aceptación**: el dataset se regenera con un comando, cubre 2019–2026 y las 10 parcelas, y reporta cobertura y huecos.
- **Tiempo estimado**: 2–3 sesiones.

## Fase 5 — Baseline FAO-56 + NDVI (método PLAS/INIA)

- **Tareas**: implementar el baseline agronómico (ET0 FAO-56 × Kcb desde NDVI) sobre la serie temporal; definir métricas de referencia.
- **Entregables**: baseline documentado + métricas (error y correlación frente a referencia).
- **Criterios de aceptación**: baseline reproducible sobre el dataset de la Fase 4 y métricas registradas como punto de comparación del modelo.
- **Tiempo estimado**: 1–2 sesiones.

## Fase 6 — Modelo deep learning

- **Tareas**: modelo secuencial (LSTM o Transformer temporal) sobre la serie multitemporal; partición **temporal** (no aleatoria); comparación contra el baseline.
- **Entregables**: script/notebook de entrenamiento, métricas y análisis de error.
- **Criterios de aceptación**: el modelo supera al baseline en la métrica acordada con partición temporal, y el resultado es reproducible.
- **Tiempo estimado**: 3–5 sesiones.

## Fase 7 — Web mínima (FastAPI + HTMX + Leaflet)

- **Tareas**: servicio FastAPI, HTMX y Leaflet; dibujar el polígono; mostrar NDVI, recomendación y explicación.
- **Entregables**: `web/` ejecutable en local.
- **Criterios de aceptación**: en local, dibujar un polígono y ver la recomendación y su explicación en pantalla.
- **Tiempo estimado**: 2 sesiones.

## Fase 8 — Paper y publicación

- **Tareas**: escritura del paper (método, dataset, baseline, modelo, resultados); release open source; difusión (LinkedIn, Reddit, universidades).
- **Entregables**: paper y release.
- **Criterios de aceptación**: paper enviado o publicado y repositorio liberado con tag de versión.
- **Tiempo estimado**: 3–4 sesiones.

## Deuda archivada

- **`internal-go-archive/`** — motor ET0 (Hargreaves-Samani) y parser CSV en Go, con sus tests. Archivado con el tag `go-motor-archive`; **se conserva como referencia**, sin mantenimiento.
  - **Por qué**: el autor es Python senior y todo el pipeline (rasters, series temporales, deep learning, web) vive en Python; mantener dos lenguajes añadía un puente por subproceso sin aportar valor.
  - **Qué se conserva**: los casos de validación FAO-56 (`doc/05-DATA-SOURCES.md`) y el parser CSV quedan como especificación de la migración a `scripts/et0.py`.
- **`archive/copernicus/`** — cliente Copernicus archivado con anterioridad (referencia histórica).
- **`go.mod.archived`** — módulo Go congelado; el `.gitignore` excluye `*.archived`.
