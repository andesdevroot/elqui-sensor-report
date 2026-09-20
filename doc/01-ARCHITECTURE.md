# 01 — Arquitectura

## Capas

```text
┌──────────────────────────────────────────────────────────────┐
│               web/ — FastAPI + HTMX + Leaflet                │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                  scripts/ — pipeline Python                  │
│               ndvi_probe.py · kcb.py · et0.py                │
│         irrigation.py · timeseries.py · deepseek.py          │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│           datos — Earth Search (STAC AWS, anónimo)           │
│                   Sentinel-2 L2A como COG                    │
└──────────────────────────────────────────────────────────────┘

Servicio externo lateral (lo llama el pipeline):

┌──────────────────────────────────────────────────────────────┐
│                 IA — DeepSeek (deepseek.py)                  │
│              números → lenguaje natural (es-CL)              │
└──────────────────────────────────────────────────────────────┘
```

**Regla de dependencias**: `web/` consume `scripts/`; los scripts hablan con Earth Search y con DeepSeek; los datos fluyen en una sola dirección (satélite → NDVI → Kcb → mm/día → lenguaje natural).

## Componentes

- `scripts/ndvi_probe.py` — búsqueda STAC (Earth Search) + NDVI por ventana del AOI.
- `scripts/timeseries.py` — serie temporal multitemporal de NDVI por parcela.
- `scripts/kcb.py` — NDVI → Kcb (método INIA).
- `scripts/et0.py` — ET0 Hargreaves-Samani (portado desde el Go archivado).
- `scripts/irrigation.py` — ET0 × Kcb → recomendación en mm/día.
- `scripts/deepseek.py` — traducción a lenguaje natural (es-CL).
- `web/` — FastAPI + HTMX + Leaflet.
- `internal-go-archive/` — motor Go archivado (histórico, sin mantenimiento).
- **Earth Search** — fuente de datos pública y anónima.

## Decisión: Python puro

- **Un solo lenguaje**: pipeline satelital, motor agronómico, modelo y web en Python. Se elimina el puente Go ↔ Python por subproceso, que agregaba superficie de error sin aportar.
- **Ecosistema**: `rasterio` y `numpy` resuelven la lectura por ventana de COG; `torch` cubre el modelo secuencial. Nada de esto existe en Go.
- **Sin binarios estáticos**: la web es un servicio FastAPI; no hay CLI compilada que distribuir.
- **Sin tokens**: Earth Search es anónimo.
- **Dependencias**: declaradas en `requirements.txt` y resueltas en un entorno virtual (ver `doc/04-CONVENTIONS.md`).

## Flujo de datos

```text
dirección/coordenadas → polígono → ndvi_probe.py (Earth Search) → NDVI → kcb.py → Kcb → irrigation.py (ET0 × Kcb) → mm/día → deepseek.py → lenguaje natural → web/
```

1. La web (o un script) recibe el polígono de la parcela.
2. `ndvi_probe.py` busca escenas en Earth Search y lee la ventana del AOI como COG → NDVI.
3. `kcb.py` convierte NDVI → Kcb con el método INIA (ver `doc/05-DATA-SOURCES.md`).
4. `irrigation.py` combina Kcb con ET0 FAO-56 (`et0.py`) → recomendación en mm/día.
5. `deepseek.py` traduce la recomendación a lenguaje natural (es-CL).
6. La web presenta el resultado.

Para el objetivo de investigación, `timeseries.py` construye la serie multitemporal que alimenta el baseline (Fase 5) y el modelo deep learning (Fase 6).
