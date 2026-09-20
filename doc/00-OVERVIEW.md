# 00 — Visión general

## Qué es elqui-eye

`elqui-eye` es un **copiloto de riego** open source para el pequeño agricultor del semiárido chileno (Región de Coquimbo). Entra con una **dirección o coordenadas de la parcela** y sale una **recomendación de riego en mm/día** más una **explicación en lenguaje natural**.

La pregunta que responde: *¿cuánta agua reponer hoy en esta parcela, y por qué?*

*Nombre de trabajo*: `elqui-eye`. El repositorio y el módulo Go conservan el nombre `elqui-sensor-report` hasta que se decida el renombre completo.

## Entrada y salida

- **Entrada**: dirección o coordenadas del predio (punto o polígono simple).
- **Salida**: recomendación en **mm/día** + explicación en **lenguaje natural** (es-CL).
- **Sin hardware propio, sin tokens y sin registro**: el dato satelital es público y anónimo.

## Stack

- **Earth Search** — catálogo STAC público de Element 84 sobre AWS: Sentinel-2 L2A como COG, acceso anónimo.
- **Python** (`pystac-client` + `rasterio`) — pipeline satelital: NDVI por ventana sobre el AOI.
- **Go** (`internal/analysis`) — motor agronómico: ET0 FAO-56, Kcb y recomendación de riego.
- **DeepSeek** — traducción de los números a lenguaje natural.
- **Web mínima** — Go + HTMX + Leaflet.

## Principios

- **Open source y gratis**: código auditable, sin licencias ni hardware.
- **Sin tokens ni credenciales**: se descartó Copernicus Data Space (OAuth/OData y un redirect que pierde el header de autorización) por fricción operativa; Earth Search es anónimo.
- **Transparencia**: cada recomendación es trazable paso a paso (NDVI → Kcb → ET0 × Kcb → mm/día).
- **Lenguaje natural**: la salida se explica, no se entrega una tabla de coeficientes.

## Relación con otras herramientas

`elqui-eye` **no compite** con las herramientas públicas chilenas; busca ser un complemento open source:

- **PLAS (INIA)** — plataforma del Instituto de Investigaciones Agropecuarias.
- **RiegaBien (Pontificia Universidad Católica de Chile)** — herramienta de apoyo a la decisión de riego.
- TODO: documentar el alcance, las fuentes y la cobertura de ambas antes de publicar comparaciones.

El aporte que `elqui-eye` busca agregar es la combinación de **imagen satelital por parcela** + **explicación en lenguaje natural** + **código abierto**, sobre fuentes gratuitas.

## Usuario objetivo

- **Pequeño agricultor de Coquimbo**: saber cuánto y cuándo regar, sin contratar un estudio.
- **Asesor técnico o cooperativa**: recomendación trazable, parcela por parcela.
- **Comunidad open source**: replicar el método en otros valles semiáridos.

## Alcance

- **Copiloto de riego**: recomendación por parcela en mm/día + explicación.
- **Pipeline satelital**: Sentinel-2 L2A vía Earth Search → NDVI → Kcb (método INIA).
- **Motor agronómico**: ET0 FAO-56 × Kcb (Go, ya construido y validado).
- **Lenguaje natural**: capa DeepSeek.
- **Web mínima**: Go + HTMX + Leaflet.

## Fuera de alcance

- Hardware y sensores propios; drones.
- App móvil nativa.
- Recomendaciones agronómicas sin validación de campo (ver Fase 5 de `doc/06-ROADMAP.md`).

## Estado

Fase 0 (reset y limpieza) en curso: `pkg/copernicus` quedó archivado y Earth Search reemplaza a Copernicus Data Space. Ver `doc/06-ROADMAP.md`.
