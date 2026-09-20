# Design — elqui-eye

> **Puntero**: la fuente de verdad de diseño vive en `doc/`:
> - Visión general → [doc/00-OVERVIEW.md](doc/00-OVERVIEW.md)
> - Arquitectura → [doc/01-ARCHITECTURE.md](doc/01-ARCHITECTURE.md)
> - Fuentes de datos → [doc/05-DATA-SOURCES.md](doc/05-DATA-SOURCES.md)
> - Roadmap → [doc/06-ROADMAP.md](doc/06-ROADMAP.md)

## Resumen en 5 líneas

1. **Qué es**: `elqui-eye`, copiloto de riego open source para el pequeño agricultor de Coquimbo: entra una dirección o coordenadas de la parcela y sale una recomendación de riego en mm/día con su explicación en lenguaje natural.
2. **Stack**: **Go** (motor ET0 FAO-56 + web con HTMX), **Python** (pipeline satelital en `scripts/`), **Earth Search** (Sentinel-2 L2A como COG, acceso anónimo) y **DeepSeek** (traducción a lenguaje natural).
3. **Flujo**: polígono → `scripts/ndvi_probe.py` → NDVI → `scripts/kcb.py` → Kcb → ET0 × Kcb → mm/día → DeepSeek → web.
4. **Metodología**: Spec-Driven Development + TDD estricto (RED → GREEN → REFACTOR), con commits atómicos.
5. **Principios**: open source, gratis, sin hardware y **sin tokens**; transparencia paso a paso y salida en lenguaje natural.

Nota: el repositorio y el módulo Go conservan el nombre `elqui-sensor-report`; `elqui-eye` es el nombre de trabajo del proyecto.
