# 06 — Roadmap

Roadmap satelital en 5 fases, más la Fase 0 de refactor de identidad. Cada tarea se cierra con un commit atómico y su criterio de aceptación cumplido (ver `doc/03-METHODOLOGY.md`).

## Fase 0 — Refactor de identidad (en progreso ✅)

Redefinir el proyecto como herramienta satelital open source `elqui-eye`, manteniendo el motor existente como fallback.

- **Tareas**: actualizar `README.md`, `doc/00`, `doc/01`, `doc/05` y `doc/06`; crear la estructura base (`pkg/copernicus`, `internal/satellite`, `internal/ai`, `web`).
- **Entregables**: documentación coherente con la nueva misión; carpetas creadas con `.gitkeep`.
- **Criterios de aceptación**: `go test ./...` en verde (el motor no se rompe) y `doc/` sin referencias a un CLI de sensores como misión principal.

## Fase 1 — Pipeline satelital

- **Tareas**: cliente `pkg/copernicus` (autenticación, búsqueda por AOI en WKT, filtro por nubosidad, descarga de B04/B08); `internal/satellite` (cálculo de NDVI y conversión NDVI → Kcb); exportación de un PNG de la parcela.
- **Entregables**: `pkg/copernicus/`, `internal/satellite/`, un PNG de NDVI para una parcela de prueba en Coquimbo.
- **Criterios de aceptación**: dado un AOI real, el pipeline obtiene las bandas y produce NDVI y PNG reproducibles; tests con bandas de ejemplo (sin red) más un test de integración contra la API marcado como opcional.

## Fase 2 — Recomendación ET0 × Kcb

- **Tareas**: integrar ET0 (motor existente) con Kcb satelital; calcular la recomendación diaria en mm/día; balance por período.
- **Entregables**: `internal/analysis` extendido; recomendación por parcela.
- **Criterios de aceptación**: recomendación en mm/día calculada con tests de tabla y casos documentados con fuente.

## Fase 3 — Capa DeepSeek

- **Tareas**: cliente `internal/ai` (net/http); prompt de traducción a lenguaje natural (es-CL); manejo de errores y timeouts; no enviar datos sensibles.
- **Entregables**: `internal/ai/`, texto de recomendación en español de Chile.
- **Criterios de aceptación**: dada una recomendación numérica, la capa produce un texto es-CL; fallback explícito si la API no responde.

## Fase 4 — Web mínima

- **Tareas**: servidor Go + HTMX + Leaflet; dibujar el polígono de la parcela; mostrar NDVI y la recomendación.
- **Entregables**: `cmd/elqui-web`, `web/`.
- **Criterios de aceptación**: en local, dibujar un polígono y recibir la recomendación en pantalla.

## Fase 5 — Validación y publicación

- **Tareas**: validar con al menos un agricultor/asesor de Coquimbo; contrastar con RiegaBien / PLAS; publicar (LinkedIn, Reddit, universidades).
- **Entregables**: caso de validación documentado; anuncio público.
- **Criterios de aceptación**: validación de campo registrada y release publicada.

## Nota sobre el roadmap anterior

Las tareas del roadmap previo se re-encuadran: el parser CSV y ET0 quedan **cerrados** y pasan a ser el motor de fallback; el cliente DGA (T1.2) queda **parked**, porque la ruta principal ahora es satelital.
