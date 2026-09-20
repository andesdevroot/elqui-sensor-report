# 01 — Arquitectura

## Capas

```text
┌──────────────────────────────────────────────────────────────┐
│                  frontend — Leaflet + HTMX                   │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│        backend Go — cmd/elqui-web + internal/analysis        │
│        ET0 FAO-56 · Kcb (INIA) · irrigation → mm/día         │
│       orquesta scripts/ (subproceso) y DeepSeek (HTTP)       │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│               scripts/ Python — pipeline NDVI                │
│      ndvi_probe.py · kcb.py (pystac-client + rasterio)       │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│          Earth Search — STAC público AWS (anónimo)           │
│                   Sentinel-2 L2A como COG                    │
└──────────────────────────────────────────────────────────────┘

Servicio externo lateral (lo llama el backend Go):

┌──────────────────────────────────────────────────────────────┐
│                 DeepSeek — servicio externo                  │
│              números → lenguaje natural (es-CL)              │
└──────────────────────────────────────────────────────────────┘
```

**Regla de dependencias**: el frontend solo habla con el backend; el backend Go orquesta (`internal/analysis` para el cálculo, `scripts/` por subproceso, DeepSeek por HTTP); los scripts Python hablan con Earth Search; `internal/analysis` es cálculo puro (sin red).

## Componentes

- `web/` + `cmd/elqui-web` — frontend (Leaflet + HTMX) y servidor.
- `internal/analysis` — motor agronómico en Go: ET0 FAO-56, Kcb e `irrigation` (recomendación en mm/día).
- `internal/ingest` — parser CSV de sensores (motor existente; fallback).
- `scripts/` — pipeline satelital en Python: `ndvi_probe.py`, `kcb.py`.
- `archive/copernicus/` — cliente Copernicus archivado: referencia histórica, ya no se usa.
- **Earth Search** — catálogo STAC público, sin credenciales.
- **DeepSeek** — traducción a lenguaje natural.

## Decisión: Python para rasters, Go para el motor y la web

- **Python** es el ecosistema real para rasters: `pystac-client`, `rasterio` y `numpy` cubren la búsqueda STAC y la lectura por ventana de un COG. Go no tiene lectura de GeoTIFF en la stdlib y reimplementarla no aporta valor.
- **Go** se queda con lo que ya funciona y está probado: el motor ET0 (validado contra FAO-56), la orquestación, la API y la web.
- **Sin dependencias Go externas**: Go sigue solo con stdlib; lo pesado se delega a subprocesos Python (mismo criterio que se usó con `gdal_translate`).
- **Sin tokens**: Earth Search es anónimo, así que desaparece la fricción de OAuth que bloqueó la descarga anterior.

## Flujo de datos

```text
dirección/coordenadas → polígono → ndvi_probe.py (Earth Search) → NDVI → kcb.py → Kcb → ET0 × Kcb → mm/día → DeepSeek → lenguaje natural → web
```

1. El frontend (Leaflet) o la CLI entrega el polígono de la parcela.
2. `scripts/ndvi_probe.py` busca escenas en Earth Search y lee la ventana del AOI como COG → NDVI.
3. `scripts/kcb.py` convierte NDVI → Kcb con el método INIA (ver `doc/05-DATA-SOURCES.md`).
4. `internal/analysis` combina Kcb con ET0 FAO-56 → recomendación en mm/día.
5. DeepSeek traduce la recomendación a lenguaje natural (es-CL).
6. La web (o la CLI) presenta el resultado.

Ruta de sensor (fallback, ya implementada):

```text
CSV del sensor → []SensorReading → ET0Result → balance hídrico → EfficiencyReport → Markdown
```

Cuando no hay imagen satelital utilizable (por ejemplo, nubosidad persistente), el motor existente produce el reporte a partir de las lecturas del sensor.
