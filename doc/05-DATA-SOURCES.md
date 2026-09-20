# 05 — Fuentes de datos

Fuente principal: **Copernicus Data Space** (Sentinel-2 L2A). Complementarias: DGA, SMAP, SMOS y CIREN.

## Copernicus Data Space — Sentinel-2 L2A (fuente principal)

- Portal: https://dataspace.copernicus.eu/ — registro gratis; **ya hecho por el autor**.
- Producto: **Sentinel-2 L2A** (reflectancia de superficie, corregida atmosféricamente).
- Uso previsto (Fase 1):
  - búsqueda por **AOI** (polígono en **WKT**) — la parcela específica del agricultor;
  - filtro por **nubosidad** (porcentaje máximo de nubes);
  - descarga de las bandas **B04 (rojo)**, **B08 (NIR)** y **SCL** (máscara de nubes).
- Cálculo: `NDVI = (B08 - B04) / (B08 + B04)` → conversión a **Kcb**.

### Endpoints validados en vivo

- **Token (OAuth2)**: `https://identity.dataspace.copernicus.eu/auth/realms/CDSE/protocol/openid-connect/token` — `POST` `application/x-www-form-urlencoded` con `grant_type=client_credentials`, `client_id` y `client_secret`; devuelve `access_token` (JWT).
- **Catálogo STAC**: `https://stac.dataspace.copernicus.eu/v1/search` — `POST` `application/json` con `collections: ["sentinel-2-l2a"]`, `intersects` (GeoJSON), `datetime: "<inicio>/<fin>"` y `limit`.
- **Descarga de bandas**: los assets del STAC apuntan a `s3://eodata/Sentinel-2/MSI/L2A/...` (`auth:refs: ["s3"]`), pero cada asset incluye una **alternativa HTTPS** en `assets.<banda>.alternate.https.href` → `https://download.dataspace.copernicus.eu/odata/v1/Products(<uuid>)/Nodes(...)/$value`, con `auth:refs: ["oidc"]`: la descarga usa el **mismo token Bearer**.
  - Claves de asset verificadas: `B04_10m` (rojo), `B08_10m` (NIR) y `SCL_20m` (máscara de nubes); también existen variantes `_20m` y `_60m`.
  - Cada `.jp2` de 10 m pesa del orden de **96 MB** (`file:size`) → la descarga debe ir a disco, no a memoria.
  - Nuestro struct `copernicus.Asset` hoy expone solo `Href` (el `s3://...`): **T1.3 debe capturar `alternate.https.href`** para descargar por HTTPS.
- Notas: el catálogo STAC también responde **en anónimo** (la búsqueda no exige token), pero el token sí es necesario para la descarga (`oidc`). La respuesta incluye un enlace `next` con `token` para paginar.
- Alternativas si la API principal resulta engorrosa: `sentinel-images-downloader`, `georeader`, `phidown`.

## DGA — Sistema Hidrométrico

- URL: https://dga.mop.gob.cl/
- Qué aporta: información hidrométrica oficial de la Dirección General de Aguas (caudales, niveles y registros asociados).
- TODO: documentar endpoint o formato de descarga (CSV, Excel, API) y condiciones de uso.
- Rol: complementaria (la ruta principal ya es satelital).

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

## Estrategia

La ruta principal es satelital: Sentinel-2 L2A → NDVI → Kcb sobre el polígono de la parcela, combinado con ET0 FAO-56. El motor de sensores + ET0 queda como **fallback** cuando no hay imagen utilizable (p. ej. nubosidad persistente). DGA, SMAP, SMOS y CIREN quedan para v2+.

## Fuente de validación ET0

Fuente primaria: **FAO-56** — Allen, R.G., Pereira, L.S., Raes, D. & Smith, M. (1998). *Crop evapotranspiration: Guidelines for computing crop water requirements*. FAO Irrigation and Drainage Paper 56. Roma: FAO.

- Ecuación 52 (Hargreaves-Samani): `ET0 = 0.0023 (Tmedia + 17.8) (Tmax - Tmin)^0.5 Ra`
- Radiación extraterrestre Ra: ecuaciones 21–25 (Cap. 3), con `Gsc = 0.0820 MJ m-2 min-1` y `λ = 2.45 MJ kg-1` (Ra en mm/día = Ra[MJ m-2 d-1] / 2.45).

### Casos de referencia

**Caso 1 — FAO-56, Ejemplo 20 (Lyon, Francia).** Ancla externa publicada.
- Latitud 45.7167°N · 15 de julio (día del año J = 196)
- Tmax = 26.6 °C · Tmin = 14.8 °C · Tmedia = 20.7 °C
- **ET0 esperado = 5.0 mm/día** (FAO-56, ec. 52; el mismo ejemplo da 4.56 mm/día con Penman-Monteith)

**Caso 2 — Valle del Elqui, verano.** Valor calculado aplicando la ec. 52.
- Latitud 30.0°S · 15 de enero (J = 15)
- Tmax = 30.0 °C · Tmin = 15.0 °C · Tmedia = 22.5 °C
- **ET0 esperado = 6.32 mm/día**

**Caso 3 — Valle del Elqui, invierno.** Valor calculado aplicando la ec. 52.
- Latitud 30.0°S · 15 de julio (J = 196)
- Tmax = 18.0 °C · Tmin = 6.0 °C · Tmedia = 12.0 °C
- **ET0 esperado = 1.90 mm/día**

### Validación del término radiativo (Ra)

**FAO-56, Ejemplo 8**: 3 de septiembre (J = 246) a 20°S → Ra = 32.2 MJ m⁻² día⁻¹ = **13.1 mm/día** en evaporación equivalente.

### Pendiente

- TODO: validación cruzada con **INIA** (estaciones del Valle del Elqui) para v2. Los casos 2 y 3 son valores calculados con la ec. 52, no mediciones publicadas; el único ancla externa publicada es el Ejemplo 20 de FAO-56.
