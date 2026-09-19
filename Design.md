# Design — elqui-eye

> **Puntero**: la fuente de verdad de diseño vive en `doc/`:
> - Visión general → [doc/00-OVERVIEW.md](doc/00-OVERVIEW.md)
> - Arquitectura → [doc/01-ARCHITECTURE.md](doc/01-ARCHITECTURE.md)
> - Modelo de datos → [doc/02-DATA-MODEL.md](doc/02-DATA-MODEL.md)
> - Roadmap → [doc/06-ROADMAP.md](doc/06-ROADMAP.md)

## Resumen en 5 líneas

1. **Qué es**: `elqui-eye`, herramienta satelital open source de recomendación de riego para pequeños agricultores de Coquimbo (Sentinel-2 L2A + ET0 FAO-56 + DeepSeek).
2. **Cómo se estructura**: `cmd/elqui` (CLI motor) + `cmd/elqui-web` (web) → `internal/{satellite, ai, analysis, ingest, report}` → `pkg/copernicus` → `internal/models`; `web/` con HTML + HTMX + Leaflet.
3. **Lenguaje**: Go, solo biblioteca estándar, binario estático.
4. **Metodología**: Spec-Driven Development + TDD estricto (RED → GREEN → REFACTOR), con commits atómicos.
5. **Rol del motor de sensores**: **fallback** cuando no hay imagen satelital utilizable (CSV del sensor → ET0 FAO-56 → reporte).

Nota: el repositorio y el módulo Go conservan el nombre `elqui-sensor-report`; `elqui-eye` es el nombre de trabajo del proyecto.
