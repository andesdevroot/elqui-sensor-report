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

![Python](https://img.shields.io/badge/Python-3.12-3776AB.svg?logo=python&logoColor=white)
![Satellite](https://img.shields.io/badge/Sentinel--2-L2A-2f6f4e.svg)
![AI](https://img.shields.io/badge/DeepSeek-lenguaje_natural-4b6bfb.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-WIP-orange.svg)

`elqui-eye` es un copiloto de riego open source para el pequeño agricultor del semiárido chileno (Región de Coquimbo) y, a la vez, un proyecto de investigación sobre detección temprana de estrés hídrico con deep learning. Entra una dirección o coordenadas de la parcela y devuelve una recomendación de riego en **mm/día** con su explicación en **lenguaje natural**.

La misión es el acceso gratuito a una recomendación específica del predio: no un promedio regional, sino la imagen de *tu* parcela, su índice de vegetación (NDVI → Kcb) y el agua que corresponde reponer. Herramientas como **RiegaBien (UC)** y **PLAS (INIA)** resuelven partes del mismo problema; `elqui-eye` no compite con ellas: busca ser la capa open source que aporta transparencia, imagen satelital por parcela y explicación en lenguaje natural.

El proyecto es **Python puro**: `pystac-client` + `rasterio` para el pipeline satelital, ET0 FAO-56 × Kcb para la recomendación, `torch` para el modelo temporal y FastAPI + HTMX + Leaflet para la web. El motor Go anterior quedó archivado en [`internal-go-archive/`](internal-go-archive/) (tag `go-motor-archive`).

## Quick Start

```bash
# 1. Clonar el repositorio
git clone https://github.com/andesdevroot/elqui-sensor-report.git
cd elqui-sensor-report

# 2. Entorno virtual y dependencias
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt          # llega en la Fase 1

# 3. Pipeline satelital: NDVI de una parcela (Fase 1)
python scripts/ndvi_probe.py --lat -29.90453 --lon -71.24894 --days 90 --cloud 20
```

> **Estado WIP**: `scripts/ndvi_probe.py` existe (sin commitear) y `requirements.txt` todavía no. El motor agronómico en Python (Fase 2), la serie temporal (Fase 4), el modelo (Fase 6) y la web (Fase 7) están en construcción; ver [doc/06-ROADMAP.md](doc/06-ROADMAP.md).
>
> **Código histórico**: el motor Go está en [`internal-go-archive/`](internal-go-archive/) (tag `go-motor-archive`) y el cliente Copernicus en `archive/copernicus/`; ninguno se mantiene.

## Estructura del proyecto

Estructura objetivo (se puebla por fases):

```text
elqui-eye/                      # repo andesdevroot/elqui-sensor-report
├── scripts/                    # pipeline Python
│   └── ndvi_probe.py           # NDVI por ventana sobre el AOI            [Fase 1]
│   # et0.py · kcb.py · irrigation.py · deepseek.py · timeseries.py        [Fases 2-4]
├── web/                        # FastAPI + HTMX + Leaflet                 [Fase 7]
├── archive/
│   └── copernicus/             # cliente Copernicus archivado (histórico)
├── internal-go-archive/        # motor Go archivado (tag go-motor-archive)
├── doc/                        # documentación — fuente de verdad
├── testdata/                   # CSV de ejemplo (histórico Go)
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
