# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- **Fase 0 — Reset y limpieza: completa.** `pkg/copernicus` quedó archivado en `archive/copernicus/` (tag `copernicus-archive`) y **Earth Search** reemplaza a Copernicus Data Space como fuente satelital.
- **Fase 1 — Pipeline satelital mínimo: arrancando.**
- Motor existente en verde: parser CSV (`internal/ingest`) + ET0 FAO-56 (`internal/analysis`).

## Último commit

- `aa847f1` — `docs(roadmap): reescribe roadmap como copiloto de riego`

## Siguiente tarea

- **Crear `scripts/ndvi_probe.py`** con `pystac-client` + `rasterio`: búsqueda STAC en Earth Search (**anónimo**), lectura **por ventana** del AOI sobre COG y cálculo del NDVI medio. Validar con la parcela La Serena `2W89+VG`.
- Después, en la misma fase: `scripts/requirements.txt` y la validación documentada de la parcela.

## Blockers

- Ninguno.

## Deuda técnica

- **`archive/copernicus/` archivado, sin mantenimiento**: su descarga rechazaba con 401 el token de cuenta de servicio y el redirect perdía el header de autorización. Referencia histórica; no se usa.
- **Cita INIA pendiente**: `doc/05` marca como TODO la cita exacta del paper que valida `Kcb = 1.51 × NDVI − 0.23` en Coquimbo.

## Decisiones recientes

- **Earth Search reemplaza a Copernicus**: catálogo STAC público y **anónimo** (sin OAuth, sin tokens, sin redirects rotos). Se adoptó **Python** (`pystac-client` + `rasterio`) para el pipeline satelital, porque es el ecosistema real para rasters; **Go** se queda con el motor ET0 y la web.
- **`pkg/copernicus` archivado, no borrado**, con el tag `copernicus-archive` como referencia.
- **Directorios vestigiales eliminados**: `internal/satellite`, `internal/ai` y `cmd/elqui`.
