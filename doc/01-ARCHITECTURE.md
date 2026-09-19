# 01 — Arquitectura

## Capas

```text
┌──────────────────────────────────────────────────────────────┐
│        cmd/elqui (CLI motor)  ·  cmd/elqui-web (web)         │
│                 web/ (HTML + HTMX + Leaflet)                 │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                          internal/                           │
│                                                              │
│          satellite (NDVI, Kcb)   ►   ai (DeepSeek)           │
│           ingest ► analysis (ET0 FAO-56) ► report            │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                             pkg/                             │
│                                                              │
│ copernicus — Sentinel-2 L2A (AOI WKT · nubosidad · B04/B08)  │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                       internal/models                        │
│         SensorReading · ET0Result · EfficiencyReport         │
└──────────────────────────────────────────────────────────────┘
```

**Regla de dependencias**: `cmd` orquesta e importa `internal/` y `pkg/`; `internal/` importa `pkg/` y `models`; `pkg/` no importa `internal/`; `internal/models` no importa nada del proyecto.

## Componentes

- `cmd/elqui` — CLI motor (análisis y recomendación).
- `cmd/elqui-web` — servidor web mínimo.
- `internal/satellite` — índices de vegetación: NDVI y conversión a Kcb.
- `internal/ai` — capa DeepSeek: traduce los números a lenguaje natural (es-CL).
- `internal/analysis` — ET0 FAO-56 y balance hídrico (motor existente).
- `internal/ingest` — parser CSV de sensores (motor existente; fallback).
- `internal/report` — render de reportes.
- `pkg/copernicus` — cliente de Copernicus Data Space (Sentinel-2 L2A).
- `web/` — HTML + HTMX + Leaflet.

## Decisión: Go

- **stdlib robusta**: `encoding/csv`, `encoding/json`, `net/http` y `time` cubren el motor y los clientes HTTP (Copernicus, DeepSeek).
- **Binario estático**: `go build -ldflags="-s -w"` produce un único binario, sin runtime ni instalación.
- **Consistencia con promptc**: mismo stack, estructura y metodología que el proyecto hermano del mismo autor (`github.com/andesdevroot/promptc`).

## Principio: sin dependencias Go externas

`go.mod` no declara dependencias. Las APIs externas (Copernicus, DeepSeek) se consumen con `net/http` de la stdlib. Si una tarea parece requerir una librería externa, primero se busca una alternativa con stdlib; si no existe, la decisión se documenta aquí antes de agregarla.

## Flujo de datos

Ruta satelital (principal):

```text
AOI (WKT) → Copernicus (B04, B08) → NDVI → Kcb → ET0 (FAO-56) × Kcb → mm/día → DeepSeek → lenguaje natural → CLI / web
```

1. `pkg/copernicus` busca la escena Sentinel-2 L2A del AOI, filtra por nubosidad y descarga B04/B08.
2. `internal/satellite` calcula NDVI y lo convierte a Kcb (coeficiente de cultivo basal).
3. `internal/analysis` combina ET0 con Kcb → recomendación en mm/día.
4. `internal/ai` traduce la recomendación a lenguaje natural (es-CL).
5. `internal/report` y `web/` presentan el resultado.

Ruta de sensor (fallback, ya implementada):

```text
CSV del sensor → []SensorReading → ET0Result → balance hídrico → EfficiencyReport → Markdown
```

Cuando no hay imagen satelital utilizable (p. ej. nubosidad persistente), el motor existente produce el reporte a partir de las lecturas del sensor.

Las fuentes y su estado se documentan en `doc/05-DATA-SOURCES.md`.
