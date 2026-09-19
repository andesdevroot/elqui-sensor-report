# 00 — Visión general

## Qué es elqui-eye

`elqui-eye` es una herramienta satelital open source para pequeños agricultores del semiárido chileno (Región de Coquimbo). Combina imágenes **Sentinel-2 L2A**, evapotranspiración de referencia **ET0 (FAO-56)** y **DeepSeek** para entregar una recomendación de riego específica de la parcela: cuánto reponer (mm/día) y por qué, en lenguaje natural.

*Nombre de trabajo*: `elqui-eye`. El repositorio y el módulo Go conservan su nombre original `elqui-sensor-report` hasta que se decida el renombre completo.

## Problema que resuelve

En Coquimbo la escasez hídrica es estructural. El pequeño agricultor no dispone de una recomendación de riego ajustada a **su** predio: las guías son regionales, la asesoría técnica es escasa y las herramientas existentes no siempre llegan al agricultor ni explican el fundamento de la recomendación.

`elqui-eye` apunta a tres brechas concretas:

- **Escala del predio**: imagen satelital del polígono específico, no un promedio regional.
- **Transparencia**: cada recomendación es auditable paso a paso (NDVI → Kcb → ET0 → mm/día); el código es open source.
- **Lenguaje natural**: la salida se explica en español de Chile, no como una tabla de coeficientes.

## Relación con otras herramientas

`elqui-eye` **no compite** con las herramientas públicas chilenas; busca ser una capa complementaria y open source:

- **RiegaBien (Pontificia Universidad Católica de Chile)** — herramienta pública de apoyo a la decisión de riego.
- **PLAS (INIA)** — herramienta del Instituto de Investigaciones Agropecuarias.
- TODO: documentar el alcance, las fuentes y la cobertura de ambas antes de publicar comparaciones.

El aporte que `elqui-eye` busca agregar es la combinación de imagen satelital por parcela + explicación en lenguaje natural + código abierto, sobre fuentes gratuitas (Copernicus, FAO-56).

## Usuario objetivo

- **Pequeño agricultor de Coquimbo**: saber cuánto y cuándo regar, sin contratar un estudio.
- **Asesor técnico o cooperativa**: recomendación trazable, parcela por parcela.
- **Comunidad open source**: replicar el método en otros valles semiáridos.

## Alcance

- **Motor de análisis (ya construido)**: parser CSV de sensores + ET0 Hargreaves-Samani validada contra FAO-56. Se mantiene como **fallback** cuando no hay dato satelital.
- **Pipeline satelital**: Sentinel-2 L2A vía Copernicus Data Space → NDVI → Kcb.
- **Recomendación**: ET0 × Kcb → mm/día.
- **Lenguaje natural**: capa DeepSeek.
- **Web mínima**: Go + HTMX + Leaflet.

## Fuera de alcance

- Drones y sensores propios.
- App móvil nativa.
- Recomendaciones agronómicas sin validación de campo (ver Fase 5 de `doc/06-ROADMAP.md`).

## Estado

Fase 0 (refactor de identidad) en curso. Ver `doc/06-ROADMAP.md`.
