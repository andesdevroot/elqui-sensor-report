# elqui-sensor-report

```text
         _                   _
   ___  | |   __ _   _   _  (_)
  / _ \ | |  / _` | | | | | | |
 |  __/ | | | (_| | | |_| | | |
  \___| |_|  \__, |  \__,_| |_|
                |_|

   Eficiencia hídrica para el Valle del Elqui
   v0.1.0-dev • by Cesar Rivas
```

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8.svg?logo=go)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-WIP-orange.svg)

`elqui-sensor-report` es un CLI open source escrito en Go que procesa datos de sensores de humedad de suelo y genera reportes de eficiencia hídrica para agricultores del Valle del Elqui, Chile. Responde una pregunta concreta: **¿está cada parcela regando con el agua que realmente necesita?**

El proyecto cruza las lecturas de los sensores con la evapotranspiración de referencia (ET0) y datos públicos (DGA, Sentinel-2) para calcular consumo real vs. óptimo por período, y entrega un reporte Markdown con semáforo de eficiencia. Está construido sin dependencias externas: solo biblioteca estándar de Go.

## Quick Start

```bash
# 1. Clonar el repositorio (TODO: repo aún sin publicar; URL provisional)
git clone https://github.com/andesdevroot/elqui-sensor-report.git
cd elqui-sensor-report

# 2. Ejecutar la suite de tests
go test -v ./...

# 3. Compilar el binario optimizado
go build -ldflags="-s -w" -o build/elqui ./cmd/elqui/main.go
```

> **Estado WIP**: v1 en desarrollo (Fase 0). La CLI todavía no procesa datos reales; ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md).

## Estructura del proyecto

Estructura objetivo de v1 (se va poblando con TDD en las próximas fases):

```text
elqui-sensor-report/
├── cmd/
│   └── elqui/            # entry point del binario (main.go)
├── internal/
│   ├── ingest/           # CSV de sensores → models
│   ├── analysis/         # ET0, balance hídrico, eficiencia
│   ├── report/           # render del reporte Markdown
│   └── models/           # structs compartidos
├── pkg/
│   ├── dga/              # datos públicos DGA
│   └── sentinel/         # NDVI/NDWI (Sentinel-2)
├── doc/                  # documentación — fuente de verdad
├── testdata/             # CSV de ejemplo para tests
├── Design.md
├── Task.md
├── LICENSE
└── README.md
```

## Documentación

La documentación completa vive en `doc/` y es la fuente de verdad del proyecto:

- [00 — Visión general](doc/00-OVERVIEW.md)
- [01 — Arquitectura](doc/01-ARCHITECTURE.md)
- [02 — Modelo de datos](doc/02-DATA-MODEL.md)
- [03 — Metodología](doc/03-METHODOLOGY.md)
- [04 — Convenciones](doc/04-CONVENTIONS.md)
- [05 — Fuentes de datos](doc/05-DATA-SOURCES.md)
- [06 — Roadmap](doc/06-ROADMAP.md)

## Autor

Cesar Rivas — Senior Software Engineer, La Serena, Chile 🇨🇱

## Licencia

MIT — ver [LICENSE](LICENSE).
