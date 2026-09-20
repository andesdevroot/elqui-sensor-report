# 00 — Visión general

## Qué es elqui-eye

`elqui-eye` es un **copiloto de riego** open source para el pequeño agricultor del semiárido chileno (Región de Coquimbo) y, a la vez, un **proyecto de investigación** sobre detección temprana de estrés hídrico con deep learning. Entra una **dirección o coordenadas de la parcela** y sale una **recomendación de riego en mm/día** con su **explicación en lenguaje natural**.

La pregunta que responde: *¿cuánta agua reponer hoy en esta parcela, y por qué?*

*Nombre de trabajo*: `elqui-eye`. El repositorio conserva el nombre `elqui-sensor-report`.

## Objetivo de investigación

- **Paper**: detección temprana de estrés hídrico con deep learning en el Valle del Coquimbo.
- **Hipótesis de trabajo**: una serie temporal multitemporal de Sentinel-2 permite anticipar el estrés hídrico antes de que sea visible en el NDVI puntual.
- **Aporte**: método abierto, reproducible y con un baseline agronómico (FAO-56) contra el cual comparar el modelo.

## Entrada y salida

- **Entrada**: dirección o coordenadas del predio (punto o polígono simple).
- **Salida**: recomendación en **mm/día** + explicación en **lenguaje natural** (es-CL).
- **Sin hardware propio, sin tokens y sin registro**: el dato satelital es público y anónimo.

## Stack definitivo: Python puro

- `numpy` — cálculo numérico.
- `rasterio` — lectura por ventana de COG (Sentinel-2 L2A).
- `pystac-client` — búsqueda STAC en Earth Search.
- `xgboost` — baseline de machine learning tabular.
- `torch` — modelo secuencial (LSTM o Transformer temporal).
- `FastAPI` + `HTMX` + `Leaflet` — web mínima.

El **motor ET0** (Hargreaves-Samani) se **migra a Python** en la Fase 2. El motor Go queda archivado como referencia histórica (tag `go-motor-archive`).

## Principios

- **Python puro**: el expertise del autor es Python senior y ahí vive el ecosistema científico (rasters, series temporales, deep learning).
- **Open source, gratis, sin hardware y sin tokens**: el dato viene de Earth Search (anónimo).
- **Transparencia**: cada recomendación es trazable paso a paso (NDVI → Kcb → ET0 × Kcb → mm/día).
- **Baseline antes que modelo**: primero FAO-56 + NDVI, después deep learning; sin baseline no hay paper.

## Relación con otras herramientas

`elqui-eye` **no compite** con las herramientas públicas chilenas; busca ser un complemento open source:

- **PLAS (INIA)** — plataforma del Instituto de Investigaciones Agropecuarias; su enfoque (baseline FAO-56 + NDVI) se adopta como referencia en la Fase 5.
- **RiegaBien (Pontificia Universidad Católica de Chile)** — herramienta de apoyo a la decisión de riego.
- TODO: documentar el alcance, las fuentes y la cobertura de ambas antes de publicar comparaciones.

## Usuario objetivo

- **Pequeño agricultor de Coquimbo**: saber cuánto y cuándo regar, sin contratar un estudio.
- **Comunidad científica y open source**: método reproducible con baseline público.
- **Asesor técnico o cooperativa**: recomendación trazable, parcela por parcela.

## Alcance y fuera de alcance

- **En alcance**: pipeline satelital NDVI, motor ET0 en Python, recomendación en mm/día, serie temporal 2019–2026, baseline FAO-56 + NDVI, modelo deep learning, web mínima y paper.
- **Fuera de alcance**: hardware y sensores propios; app móvil nativa; despliegue productivo multi-tenant.

## Estado

Fase 0 (congelamiento y migración) **completa**: el motor Go quedó archivado y el proyecto pasa a Python puro. Ver `doc/06-ROADMAP.md`.
