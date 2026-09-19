# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- **Fase 1 — Pipeline satelital, arrancando.** La Fase 0 (refactor de identidad `elqui-eye`) quedó cerrada.
- Motor de sensores intacto y en verde: parser CSV (`internal/ingest`) + ET0 FAO-56 (`internal/analysis`); es el fallback.
- Estructura base creada: `pkg/copernicus`, `internal/satellite`, `internal/ai`, `web` (cada una con `.gitkeep`).

## Último commit

- `1fb11f8` — `docs: redefine proyecto como herramienta satelital open source elqui-eye` (refactor de identidad; la estructura base entró después, en `4b807e6`)

## Siguiente tarea

- **`pkg/copernicus` — autenticación + búsqueda por AOI.** Cliente de Copernicus Data Space con `net/http`: autenticación, búsqueda de escenas Sentinel-2 L2A por polígono WKT y filtro por nubosidad. El TODO del endpoint exacto vive en `doc/05-DATA-SOURCES.md`.
- Después, en la misma fase: descarga de B04/B08 → NDVI → Kcb → PNG (ver `doc/06-ROADMAP.md`).

## Blockers

- **Endpoint exacto de Copernicus sin validar**: hay que confirmar el flujo de autenticación (token/OAuth) y el endpoint de búsqueda al implementar `pkg/copernicus` (TODO en `doc/05-DATA-SOURCES.md`).
- Sin blockers técnicos: `go 1.23.0` operativo, `gh` autenticado como `andesdevroot` y el registro en Copernicus ya está hecho.

## Decisiones recientes

- **Reencuadre a `elqui-eye`**: herramienta satelital open source (Sentinel-2 L2A + ET0 FAO-56 + DeepSeek); el motor de sensores pasa a ser **fallback**. Repo y módulo conservan el nombre `elqui-sensor-report`.
- **Ruta principal**: AOI (WKT) → Copernicus (B04, B08) → NDVI → Kcb → ET0 × Kcb → mm/día → DeepSeek → lenguaje natural.
- **Sin dependencias Go externas**: Copernicus y DeepSeek se consumen con `net/http` de la stdlib.
- **`Task.md`, `doc/07` y `Design.md` reencuadrados** al roadmap satelital de `doc/06`; el roadmap previo se documenta como re-encuadrado en la nota final de `doc/06`.
