# 04 — Convenciones

## Estilo de código Go

- `gofmt` siempre; `golint` sin hallazgos.
- Nombres cortos y en inglés (`r`, `nm`, `pct`, `et0`), sin abreviaturas inventadas.
- Todo símbolo exportado lleva un comentario que empieza con su nombre.

## Paquetes

- `internal/<dominio>` para lo específico del proyecto (`ingest`, `analysis`, `report`, `models`).
- `pkg/<lib>` solo si la librería es reutilizable fuera del proyecto (`dga`, `sentinel`).
- Sin código suelto en la raíz; todo archivo pertenece a un paquete.

## Errores

- Wrapping con contexto: `fmt.Errorf("parsear fecha: %w", err)`.
- Nunca `panic` en librerías: los errores se devuelven y `cmd` decide cómo reportarlos.
- Usar `errors.Is` / `errors.As` cuando exista un error centinela.

## Tests

- Tabla de casos: slice de structs + `t.Run` para cada subtest.
- Nombres `TestXxx_Escenario` (p. ej. `TestParseSensorCSV_HumedadFueraDeRango`).
- Los tests viven junto al código, en el mismo paquete.

## Commits

- Formato: `tipo(scope): mensaje en español, imperativo`, con `tipo` ∈ `feat|fix|refactor|docs|test|chore`.
- Ejemplo: `docs: inicializa estructura doc/ como fuente de verdad del proyecto`.

## Branches

- `main` protegida una vez publicado el repositorio; features en `feat/<nombre-corto>`, fixes en `fix/<nombre-corto>`.
- TODO: definir si la rama actual `master` se renombra a `main`.

## Idioma

- Código (paquetes, símbolos, nombres de archivo): inglés.
- Comentarios, documentación y mensajes de commit: español.
