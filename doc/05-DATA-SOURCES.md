# 05 — Fuentes de datos

v1 prioriza dos fuentes: **DGA + Sentinel-2**. El resto queda documentado para no perder el hilo.

## DGA — Sistema Hidrométrico

- URL: https://dga.mop.gob.cl/
- Qué aporta: información hidrométrica oficial de la Dirección General de Aguas (caudales, niveles y registros asociados).
- TODO: documentar endpoint o formato de descarga (CSV, Excel, API) y condiciones de uso.
- Rol en v1: **fuente prioritaria**.

## Copernicus Sentinel-2

- Acceso: [Copernicus Data Space](https://dataspace.copernicus.eu/) (registro gratis).
- Qué aporta: imágenes multiespectrales; índices NDVI/NDWI para estimar estrés hídrico de la vegetación.
- TODO: definir el mecanismo de obtención (descarga manual vs. API) y el preprocesamiento mínimo.
- Rol en v1: **fuente prioritaria**.

## CIREN — Coquimbo

- Plataforma de monitoreo hídrico en desarrollo.
- TODO: investigar si publica API o datos descargables.
- Rol: v2+.

## SMAP (NASA)

- Humedad de suelo con resolución de 9 km.
- Acceso: [NASA Earthdata](https://www.earthdata.nasa.gov/) (registro gratis).
- Rol: v2+.

## SMOS (ESA)

- Humedad de suelo con más de 15 años de registros.
- Acceso: Copernicus Data Space (gratis).
- Rol: v2+.

## Estrategia v1

Priorizar **DGA + Sentinel-2**: son las fuentes más directas para el balance hídrico y el estrés hídrico. CIREN, SMAP y SMOS quedan para v2+: su resolución y formatos agregan complejidad que v1 no necesita.
