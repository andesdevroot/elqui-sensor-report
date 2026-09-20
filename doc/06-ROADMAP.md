# 06 — Roadmap

Roadmap del copiloto de riego en 6 fases. Cada tarea se cierra con un commit atómico y su criterio de aceptación cumplido (ver `doc/03-METHODOLOGY.md`).

## Fase 0 — Reset y limpieza (en progreso ✅)

- **Tareas**: archivar `pkg/copernicus` en `archive/copernicus/` (tag `copernicus-archive`); adoptar **Earth Search** como fuente satelital; reescribir `doc/00`, `doc/01`, `doc/05` y `doc/06`.
- **Entregables**: repositorio sin dependencia operativa de Copernicus Data Space; documentación alineada al copiloto de riego.
- **Criterios de aceptación**: `go test ./...` en verde (incluido el código archivado) y `doc/` sin la ruta Copernicus como fuente principal.

## Fase 1 — Pipeline satelital mínimo

- **Tareas**: `scripts/ndvi_probe.py` con `pystac-client` + `rasterio` (lectura por ventana del AOI sobre COG); validación sobre la parcela La Serena `2W89+VG`; `requirements.txt`.
- **Entregables**: `scripts/ndvi_probe.py`, NDVI del AOI para la parcela de prueba, `requirements.txt`.
- **Criterios de aceptación**: dado el polígono y una ventana de fechas, el script devuelve el NDVI medio del AOI **sin credenciales** y **sin descargar la escena completa** (lectura por ventana), indicando la escena elegida.

## Fase 2 — Integración NDVI → Kcb → ET0 (Go)

- **Tareas**: `scripts/kcb.py` (NDVI → Kcb con el método INIA); `internal/analysis` calcula la recomendación `ET0 × Kcb` en mm/día e integra la salida de los scripts.
- **Entregables**: `internal/analysis` extendido con la recomendación (`irrigation`); mm/día por parcela.
- **Criterios de aceptación**: caso documentado NDVI → Kcb → mm/día, con test de tabla en Go y la fuente de la fórmula citada.

## Fase 3 — Capa DeepSeek

- **Tareas**: cliente en Go (`net/http`) con timeouts y manejo de errores; prompt de traducción a lenguaje natural (es-CL); no enviar datos sensibles.
- **Entregables**: texto de recomendación en español de Chile.
- **Criterios de aceptación**: dada una recomendación numérica, la capa produce un texto es-CL y hay fallback explícito si la API no responde.

## Fase 4 — Web mínima (Go + HTMX + Leaflet)

- **Tareas**: servidor Go, HTMX y Leaflet; dibujar el polígono de la parcela; mostrar NDVI y la recomendación.
- **Entregables**: `cmd/elqui-web` + `web/`.
- **Criterios de aceptación**: en local, dibujar un polígono y ver la recomendación en pantalla.

## Fase 5 — Validación con agricultores y publicación

- **Tareas**: validar con al menos un agricultor o asesor de Coquimbo; contrastar con PLAS (INIA) y RiegaBien (UC); publicar (LinkedIn, Reddit, universidades).
- **Entregables**: caso de validación documentado; anuncio público.
- **Criterios de aceptación**: validación de campo registrada y release publicada.

## Deuda técnica archivada

- **`archive/copernicus/`** — cliente OAuth2 + búsqueda STAC + descarga de bandas + conversión a GeoTIFF para Copernicus Data Space (con sus tests). Se **archiva, no se borra**: queda como referencia histórica.
  - **Motivo del abandono**: el servicio de descarga rechazaba con **401** el token de cuenta de servicio (`client_credentials`) y exige login de usuario; además su redirect pierde el header de autorización. Earth Search es anónimo y elimina esa fricción.
  - **Estado**: etiquetado con el tag `copernicus-archive`, sin mantenimiento.
  - **Qué se conserva de esa etapa**: el criterio de "sin dependencias Go externas" y el patrón de invocar herramientas del sistema por subproceso (usado con `gdal_translate` y ahora con Python).
