# elqui-eye

```text
         _                   _
   ___  | |   __ _   _   _  (_)           ___   _   _    ___
  / _ \ | |  / _` | | | | | | |  _____   / _ \ | | | |  / _ \
 |  __/ | | | (_| | | |_| | | | |_____| |  __/ | |_| | |  __/
  \___| |_|  \__, |  \__,_| |_|          \___|  \__, |  \___|
                |_|                             |___/

   riego satelital para pequeños agricultores de Coquimbo
   v0.2.0-dev • by Cesar Rivas
```

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8.svg?logo=go)
![Satellite](https://img.shields.io/badge/Sentinel--2-L2A-2f6f4e.svg)
![AI](https://img.shields.io/badge/DeepSeek-lenguaje_natural-4b6bfb.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-WIP-orange.svg)

`elqui-eye` es una herramienta satelital open source para pequeños agricultores del semiárido chileno (Región de Coquimbo). Combina imágenes **Sentinel-2 L2A**, evapotranspiración de referencia **ET0 (FAO-56)** y **DeepSeek** para traducir el estado hídrico de una parcela a una recomendación de riego en **mm/día** y en lenguaje natural.

La misión es el acceso gratuito a una recomendación específica del predio: no un promedio regional, sino la imagen de *tu* parcela, su índice de vegetación (NDVI → Kcb) y el agua que corresponde reponer. Herramientas como **RiegaBien (UC)** y **PLAS (INIA)** resuelven partes del mismo problema; `elqui-eye` no compite con ellas: busca ser la capa open source que aporta transparencia, imagen satelital por parcela y explicación en lenguaje natural.

El motor de análisis ya construido (parser CSV de sensores + ET0 Hargreaves-Samani validada contra FAO-56) se mantiene como **fallback** cuando no hay dato satelital disponible.

## Quick Start

```bash
# 1. Clonar el repositorio
git clone https://github.com/andesdevroot/elqui-sensor-report.git
cd elqui-sensor-report

# 2. Ejecutar la suite de tests
go test -v ./...
```

> **Estado WIP**: la capa satelital (Fase 1) aún no existe; hoy funciona el motor de análisis. Ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md).
>
> **CLI en construcción (Fase 2)**: el binario `cmd/elqui` todavía no existe, por eso no se documenta un `go build` aquí.

## Estructura del proyecto

Estructura objetivo (se puebla por fases):

```text
elqui-eye/                      # repo andesdevroot/elqui-sensor-report
├── cmd/
│   ├── elqui/                  # CLI motor (análisis)
│   └── elqui-web/              # servidor web
├── internal/
│   ├── ingest/                 # CSV de sensores → models
│   ├── analysis/               # ET0 FAO-56, balance hídrico
│   ├── satellite/              # NDVI, Kcb
│   ├── ai/                     # DeepSeek (lenguaje natural)
│   ├── report/                 # render del reporte
│   └── models/                 # structs compartidos
├── pkg/
│   └── copernicus/             # cliente Sentinel-2 L2A
├── web/                        # HTML + HTMX + Leaflet
├── doc/                        # documentación — fuente de verdad
├── testdata/                   # CSV de ejemplo para tests
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
- [07 — Handoff entre sesiones](doc/07-HANDOFF.md)

## Autor

Cesar Rivas — Senior Software Engineer, La Serena, Chile 🇨🇱

## Licencia

MIT — ver [LICENSE](LICENSE).
