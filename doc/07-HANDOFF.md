# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- **Fase 1 — Pipeline satelital, en progreso.** T1.1 cerrada: autenticación OAuth2 validada **en vivo** contra Copernicus Data Space (token JWT OK).
- Motor de sensores intacto y en verde: parser CSV (`internal/ingest`) + ET0 FAO-56 (`internal/analysis`); es el fallback.
- Estructura: `pkg/copernicus` (ya con código), `internal/satellite`, `internal/ai`, `web`.

## Último commit

- `893dfbd` — `feat(copernicus): cliente OAuth2 client credentials para Copernicus Data Space` (validado en vivo: token JWT OK). Después entró `chore(copernicus): elimina .gitkeep vestigial`.

## Siguiente tarea

- **`pkg/copernicus` — búsqueda STAC por AOI.** Implementar la consulta al catálogo STAC de CDSE: polígono **WKT** de la parcela (referencia: La Serena `2W89+VG`), filtro por **nubosidad** y colección **Sentinel-2 L2A**, para obtener las escenas candidatas.
- Después, en la misma fase: descarga de B04/B08 → NDVI → Kcb → PNG (ver `doc/06-ROADMAP.md`).

## Blockers

- **Endpoint STAC sin validar**: falta confirmar la URL del catálogo STAC de CDSE y el formato exacto de la consulta (colección, intersección por AOI, propiedad de nubosidad). TODO en `doc/05-DATA-SOURCES.md`.
- Sin blockers técnicos: `go 1.23.0` operativo, `gh` autenticado como `andesdevroot`, registro en Copernicus hecho y **auth validado en vivo**.

## Decisiones recientes

- **T1.1 cerrada**: `pkg/copernicus/client.go` obtiene token OAuth2 `client_credentials`; credenciales por entorno (`CDSE_CLIENT_ID`, `CDSE_CLIENT_SECRET`); validado contra CDSE real (token JWT OK). Tests con `httptest` (200 / 401 / sin credenciales / entorno).
- **Búsqueda vía STAC**: las escenas Sentinel-2 L2A se obtendrán consultando un catálogo STAC (no scraping ni descarga manual), filtrando por AOI y nubosidad.
- **Parcela de referencia**: La Serena `2W89+VG` (Plus Code) para las pruebas del pipeline.
- **`.gitkeep` eliminado** de `pkg/copernicus` al existir código real en la carpeta.
