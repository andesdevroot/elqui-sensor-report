# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- **Migración a Python puro: completa.** El motor Go quedó archivado (`internal-go-archive/`, tag `go-motor-archive`) y el módulo congelado (`go.mod.archived`).
- **Fase 1 — Pipeline satelital mínimo: en curso.** `scripts/ndvi_probe.py` ya existe en el working tree (sin commitear): Earth Search anónimo + lectura por ventana con `rasterio` + NDVI.
- Sin código Go activo: ya no aplica `go test`; los tests pasan a `pytest` (desde la Fase 2).

## Último commit

- `b4ad769` — `chore: archiva motor Go y prepara migración a Python puro` (tag `go-motor-archive`)

## Siguiente tarea

- **Validar y commitear `scripts/ndvi_probe.py`** (ya escrito, sin commitear) contra la parcela La Serena `2W89+VG` (`--lat -29.90453 --lon -71.24894`), y añadir `scripts/requirements.txt` (`pystac-client`, `rasterio`, `numpy`).
- Después, Fase 2: `scripts/et0.py` portado desde el Go archivado.

## Blockers

- Ninguno.

## Deuda técnica

- **`internal-go-archive/` y `go.mod.archived`**: referencia histórica, sin mantenimiento (el `.gitignore` excluye `*.archived`).
- **Cita INIA pendiente**: `doc/05-DATA-SOURCES.md` marca como TODO la cita exacta del paper que valida `Kcb = 1.51 × NDVI − 0.23` en Coquimbo.
- **Sin `requirements.txt` todavía**: lo agrega la Fase 1; hasta entonces las dependencias se instalan a mano en `.venv/`.

## Decisiones recientes

- **Python puro**: por expertise del autor (Python senior) y porque el ecosistema científico (rasterio, numpy, torch) vive en Python; un solo lenguaje elimina el puente por subproceso entre Go y Python.
- **Go archivado, no borrado**: `git mv` preserva la historia completa y el tag `go-motor-archive` marca el punto de congelamiento.
- **`scripts/ndvi_probe.py` ya escrito sin commitear**: el pipeline satelital arrancó antes del cierre de la migración; falta su validación y su commit.
