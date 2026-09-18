# 07 — Handoff entre sesiones

Documento de retoma: qué se hizo, qué sigue y qué está bloqueado. Se actualiza al cerrar cada sesión (ver `doc/03-METHODOLOGY.md`).

## Estado actual

- Fase 2 — Análisis, en curso. **T2.1 cerrada**: `ET0Hargreaves` implementada y validada contra FAO-56 (Ej. 20 y Ej. 8).
- Fase 1 — Ingesta: T1.1 y T1.3 cerradas; T1.2 bloqueada por falta de documentación de la fuente DGA.
- Suite en verde (`go test ./...`). Repositorio publicado en https://github.com/andesdevroot/elqui-sensor-report (rama `main`).

## Último commit

- `a13442f` — `test(analysis): agrega test RED de ET0 Hargreaves-Samani`

## Siguiente tarea

- **T2.2 — Balance hídrico**: consumo real vs. óptimo por período, con tests de tabla (criterio en `doc/06-ROADMAP.md`).
- T1.2 (cliente DGA) sigue bloqueada.

## Blockers

- **T1.2 bloqueada por falta de información de la fuente**: no hay endpoint ni formato de descarga documentado para la DGA (TODO en `doc/05-DATA-SOURCES.md`).
- Sin blockers técnicos.

## Decisiones recientes

- **Fuente de validación ET0 = FAO-56** (Allen et al. 1998): ec. 52 (Hargreaves-Samani) + ec. 21–25 para Ra. Documentada en `doc/05-DATA-SOURCES.md`.
- **Ancla externa**: FAO-56 Ej. 20 (Lyon, 45.7167°N, 15-jul, 26.6/14.8 °C) → Hargreaves = 5.0 mm/día; la implementación reproduce 5.035 (tol. 0.1). Ra validada con FAO-56 Ej. 8 (3-sep, 20°S → 13.1 mm/día).
- **API**: `analysis.ET0Hargreaves(date time.Time, latDeg, tMaxC, tMinC float64) (models.ET0Result, error)`, con `ErrInvalidInput` si la latitud sale de [-90, 90] o si Tmax < Tmin. Tmedia se deriva como (Tmax + Tmin) / 2.
- **`Confianza` fijada en 1.0 provisional**: la semántica del campo sigue como TODO en `doc/02-DATA-MODEL.md`.
- **Casos 2 y 3 (Elqui) son valores calculados** con la ec. 52, no mediciones publicadas; queda TODO validación cruzada con INIA para v2.
