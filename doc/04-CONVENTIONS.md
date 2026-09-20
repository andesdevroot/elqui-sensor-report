# 04 — Convenciones

## Estilo de código Go

- `gofmt` siempre; `golint` sin hallazgos.
- Nombres cortos y en inglés (`r`, `nm`, `pct`, `et0`), sin abreviaturas inventadas.
- Todo símbolo exportado lleva un comentario que empieza con su nombre.

## Estilo de código Python

- **PEP 8** siempre; formato con **`black`**.
- **Type hints** en toda función pública (`def ndvi_medio(aoi: str, dias: int = 90) -> float:`).
- Un script = una responsabilidad; argumentos por CLI (`argparse`) y salida legible.
- Dependencias fijadas en `scripts/requirements.txt`.

## Paquetes y directorios

- `internal/<dominio>` para lo específico del proyecto en Go; hoy: `analysis` (motor: ET0, Kcb, recomendación), `ingest` (parser CSV) y `models`.
- `scripts/` para el pipeline satelital en Python (`ndvi_probe.py`, `kcb.py`).
- `pkg/<lib>` solo si una librería Go es reutilizable fuera del proyecto; hoy no hay ninguna.
- `archive/` para código retirado que se conserva como referencia histórica (sin mantenimiento).
- Sin código suelto en la raíz.

## Errores

- Wrapping con contexto: `fmt.Errorf("parsear fecha: %w", err)`.
- Nunca `panic` en librerías: los errores se devuelven y `cmd` decide cómo reportarlos.
- Usar `errors.Is` / `errors.As` cuando exista un error centinela.

## Tests

- Tabla de casos: slice de structs + `t.Run` para cada subtest.
- Nombres `TestXxx_Escenario` (p. ej. `TestParseSensorCSV_HumedadFueraDeRango`).
- Los tests viven junto al código, en el mismo paquete.

## Fixtures

- Los archivos de `testdata/` son **inmutables**: no se editan ni se "arreglan"; representan un caso fijo y versionado.
- Si necesitas otro caso, **crea un archivo nuevo** (p. ej. `sensor_sample_2.csv`), nunca modifiques uno existente.
- Un test que falla por un fixture cambiado es drift del working tree, no una razón para editar el fixture.
- Antes de correr la suite, `git status` debe estar limpio en `testdata/`.

## Commits

- Formato: `tipo(scope): mensaje en español, imperativo`, con `tipo` ∈ `feat|fix|refactor|docs|test|chore`.
- Ejemplo: `docs: inicializa estructura doc/ como fuente de verdad del proyecto`.

## Branches

- `main` es la rama por defecto del repositorio publicado; features en `feat/<nombre-corto>`, fixes en `fix/<nombre-corto>`.

## Idioma

- Código (paquetes, símbolos, nombres de archivo): inglés.
- Comentarios, documentación y mensajes de commit: español.
