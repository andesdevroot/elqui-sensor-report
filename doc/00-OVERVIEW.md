# 00 — Visión general

## ¿Qué es elqui-sensor-report?

`elqui-sensor-report` es un CLI open source escrito en Go que procesa datos de sensores de humedad de suelo y genera reportes de eficiencia hídrica para agricultores del Valle del Elqui, Chile.

Responde una pregunta concreta: **¿está cada parcela regando con el agua que realmente necesita?**

## Problema que resuelve

En el Valle del Elqui la escasez hídrica es estructural y la mayoría de los riegos se programa por calendario o intuición, sin medición de por medio. Sin una línea base de consumo real vs. demanda óptima, es imposible:

- validar la eficiencia hídrica **individual** de cada parcela;
- detectar sobre-riego (desperdicio, lixiviación) o sub-riego (estrés hídrico);
- justificar decisiones de riego con datos.

`elqui-sensor-report` cruza las lecturas del sensor de humedad con estimaciones de evapotranspiración de referencia (ET0) y datos públicos, y entrega un reporte legible con semáforo de eficiencia por período.

## Usuario objetivo

- **Agricultor pequeño/mediano**: saber si su riego es eficiente sin contratar un estudio.
- **Agrónomo**: análisis repetible por parcela para asesorar a sus clientes.
- **Cooperativa**: comparar eficiencia entre parcelas asociadas y priorizar apoyo técnico.

## Alcance v1

- Ingesta de un CSV de sensor de humedad de suelo.
- Integración con datos de la DGA (TODO: documentar formato de descarga en `doc/05-DATA-SOURCES.md`).
- Cálculo de ET0 y balance hídrico por período.
- Reporte en Markdown con semáforo de eficiencia.

## Fuera de alcance v1

- Drones y teledetección embarcada.
- Dashboard web.
- Aplicación móvil.
- Integraciones CIREN, SMAP y SMOS (planificadas para v2+, ver `doc/05-DATA-SOURCES.md`).
