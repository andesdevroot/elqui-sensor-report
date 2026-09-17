# Design — elqui-sensor-report

> **Puntero**: la fuente de verdad de diseño vive en `doc/`:
> - Arquitectura → [doc/01-ARCHITECTURE.md](doc/01-ARCHITECTURE.md)
> - Modelo de datos → [doc/02-DATA-MODEL.md](doc/02-DATA-MODEL.md)
> - Metodología → [doc/03-METHODOLOGY.md](doc/03-METHODOLOGY.md)

## Resumen en 5 líneas

1. **Qué es**: CLI en Go que convierte lecturas de sensores de humedad de suelo y datos públicos en reportes de eficiencia hídrica por parcela.
2. **Cómo se estructura**: `cmd/elqui` → `internal/{ingest, analysis, report}` → `pkg/{dga, sentinel}` → `internal/models`.
3. **Lenguaje**: Go, solo biblioteca estándar en v1, binario estático.
4. **Metodología**: Spec-Driven Development + TDD estricto (RED → GREEN → REFACTOR), con commits atómicos.
5. **Alcance v1**: CSV del sensor + datos DGA + reporte Markdown con semáforo.
