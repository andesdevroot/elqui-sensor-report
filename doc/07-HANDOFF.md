# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- **Fase 1 — Pipeline satelital, en progreso.** T1.1 (auth) y T1.2 (búsqueda STAC) cerradas; T1.2 validada **en vivo**: 15 items con `cloud_cover < 20` para la parcela La Serena `2W89+VG`.
- Motor de sensores intacto y en verde: parser CSV (`internal/ingest`) + ET0 FAO-56 (`internal/analysis`); es el fallback.
- Estructura: `pkg/copernicus` (auth + búsqueda); `internal/satellite`, `internal/ai` y `web` siguen vacías.

## Último commit

- `d6fa587` — `feat(copernicus): búsqueda STAC por AOI con filtro de nubosidad` (validado en vivo: 15 items para La Serena `2W89+VG`; imagen más limpia con 0.19 % de nubes, 2026-09-01)

## Siguiente tarea

- **T1.3 — Descarga de bandas B04, B08 y SCL** para el AOI: usar los `assets[].href` de los items que devuelve `SearchByAOI`. `SCL` (Scene Classification Layer) se descarga para enmascarar nubes y píxeles inválidos antes de calcular el NDVI.
- Después, en la misma fase: NDVI → Kcb → PNG (ver `doc/06-ROADMAP.md`).

## Blockers

- Sin blockers técnicos: endpoint STAC validado en vivo, auth OAuth2 operativa, `gh` autenticado como `andesdevroot` y registro en Copernicus hecho.
- Por confirmar al implementar T1.3: esquema de acceso de los `href` de los assets (HTTPS directo vs. requerir token también en la descarga).

## Deuda técnica

- **Parser WKT limitado**: `parseWKT` soporta solo `POINT` y `POLYGON` simple (sin `MULTIPOLYGON`, huecos anidados ni altitud). La **Fase 4 (web)** debe aceptar **GeoJSON directo** desde Leaflet, sin pasar por el parser WKT.
- **Paginación STAC**: `limit: 100` fijo y sin seguimiento del `next` del FeatureCollection; en AOIs grandes podría truncar resultados.

## Decisiones recientes

- **T1.2 cerrada y validada en vivo**: `SearchByAOI` consulta el catálogo STAC de CDSE por AOI (WKT → GeoJSON) y filtra la nubosidad **en Go** (no en el body), para mantener la consulta simple y debuggeable. Resultado real: 15 items, la imagen más limpia con 0.19 % de nubes (2026-09-01).
- **SCL en la etapa siguiente**: la descarga incluirá la banda SCL para enmascarar nubes antes del NDVI.
- **Reintentos**: backoff exponencial ante 5xx, configurable por cliente (`MaxRetries`, `RetryBaseDelay`).
