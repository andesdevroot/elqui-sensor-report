# elqui-eye

```text
         _                   _
   ___  | |   __ _   _   _  (_)           ___   _   _    ___
  / _ \ | |  / _` | | | | | | |  _____   / _ \ | | | |  / _ \
 |  __/ | | | (_| | | |_| | | | |_____| |  __/ | |_| | |  __/
  \___| |_|  \__, |  \__,_| |_|          \___|  \__, |  \___|
                |_|                             |___/
```

**Copiloto de riego open source para pequeños agricultores del semiárido chileno.**

![Python](https://img.shields.io/badge/Python-3.12-3776AB.svg?logo=python&logoColor=white)
![Satellite](https://img.shields.io/badge/Sentinel--2-L2A-2f6f4e.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-WIP-orange.svg)

## ¿Qué es?

Responde una pregunta simple: **¿cuánto debo regar esta semana?**

- **Entrada**: las coordenadas de tu parcela.
- **Salida**: una recomendación de riego en **mm/día** y su explicación en **lenguaje natural**.

Sin hardware, sin tokens y sin costo: el dato satelital es público y anónimo.

## Cómo funciona

```text
Sentinel-2 (Earth Search) → NDVI → Kcb (INIA) → ETc = ET0 × Kcb → DeepSeek → mm/día
```

1. **Sentinel-2 vía Earth Search** — catálogo STAC público y anónimo; se lee solo la ventana de la parcela sobre un COG, sin descargar la escena completa.
2. **NDVI** — estado de la vegetación en el predio.
3. **Kcb (método INIA)** — `Kcb = 1.51 × NDVI − 0.23`.
4. **ETc = ET0 × Kcb** — ET0 FAO-56 (Hargreaves-Samani) por la demanda del cultivo.
5. **DeepSeek** — traduce los números a una recomendación en español de Chile.

## Estado

**Fase 1 completa**: el pipeline satelital está validado en vivo contra dos parcelas reales de Coquimbo, sin tokens, sin OAuth y sin Copernicus.

## Resultados validados

| Parcela | Coordenadas | NDVI (media) |
|---|---|---|
| Calle urbana, La Serena | -29.90453, -71.24894 | 0.08 — suelo desnudo |
| Viña, Vicuña | -30.032, -70.712 | 0.40 — vegetación activa |

## Roadmap

- [x] **Fase 1** — Pipeline satelital mínimo (`ndvi_probe.py`, validado con la parcela La Serena `2W89+VG`)
- [ ] **Fase 2** — Motor agronómico en Python (`et0.py`, `kcb.py`, `irrigation.py`)
- [ ] **Fase 3** — Capa DeepSeek (`deepseek.py`, prompt en español chileno)
- [ ] **Fase 4** — Serie temporal multitemporal (`timeseries.py`, 2019–2026, 10 parcelas de Coquimbo)
- [ ] **Fase 5** — Baseline FAO-56 + NDVI (método PLAS/INIA)
- [ ] **Fase 6** — Modelo deep learning (LSTM o Transformer temporal)
- [ ] **Fase 7** — Web mínima (FastAPI + HTMX + Leaflet)
- [ ] **Fase 8** — Escritura del paper + publicación + release open source

El detalle de cada fase (entregables, criterios y tiempos) está en [doc/06-ROADMAP.md](doc/06-ROADMAP.md).

## Objetivo de investigación

Detección **temprana** de estrés hídrico con deep learning sobre series temporales Sentinel-2, con un baseline agronómico (FAO-56) como punto de comparación. Paper objetivo: **Computers and Electronics in Agriculture**.

## Inspiración

**PLAS (INIA)** y **RiegaBien (UC)** resolvieron partes de este problema antes que nosotros. Este proyecto **no compite** con ellas: busca ser una capa open source que aporte transparencia, imagen satelital por parcela y explicación en lenguaje natural.

## Instalación

```bash
git clone https://github.com/andesdevroot/elqui-sensor-report.git
cd elqui-sensor-report
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## Uso

```bash
python3 scripts/ndvi_probe.py --lat -30.032 --lon -70.712 --days 90 --cloud 20
```

## Código histórico

El motor Go anterior quedó archivado en `internal-go-archive/` (tag `go-motor-archive`) y el cliente Copernicus en `archive/copernicus/`; ninguno se mantiene.

## Autor

Cesar Rivas — Senior Software Engineer, La Serena, Chile 🇨🇱

## Licencia

MIT — ver [LICENSE](LICENSE).
