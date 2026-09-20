# 04 — Convenciones

## Estilo de código Python

- **PEP 8** siempre; formato con **`black`**.
- **Type hints** en toda función pública (`def ndvi_medio(aoi: str, dias: int = 90) -> float:`).
- Un módulo = una responsabilidad; argumentos por CLI (`argparse`) y salida legible.
- Docstrings cortos en español, en modo imperativo.

## Estructura

- `scripts/` — pipeline satelital y agronómico: un archivo por etapa (`ndvi_probe.py`, `et0.py`, `kcb.py`, `irrigation.py`, `deepseek.py`, `timeseries.py`).
- `web/` — FastAPI + HTMX + Leaflet.
- `archive/` e `internal-go-archive/` — código retirado que se conserva como referencia (sin mantenimiento).
- Sin código suelto en la raíz.

## Dependencias

- Declaradas en `requirements.txt`, con versiones fijadas.
- Resolver y congelar con **pip-tools** o **poetry**.
- Entorno virtual obligatorio (`.venv/`, ya en `.gitignore`); nunca instalar en el Python del sistema.

## Errores

- Excepciones propias del dominio cuando aporten contexto; nunca un `except:` desnudo.
- No silenciar errores: registrar o re-lanzar con contexto.

## Tests

- **pytest**; `@pytest.mark.parametrize` para tablas de casos.
- Nombres `test_<comportamiento>_<escenario>`.
- Los tests viven junto al código (`test_*.py`).

## Datos de prueba

- Los fixtures de test (como `internal-go-archive/testdata/`) son **inmutables**: representan un caso fijo y versionado.
- Si necesitas otro caso, crea un archivo nuevo; nunca modifiques uno existente.

## Commits

- Formato: `tipo(scope): mensaje en español, imperativo`, con `tipo` ∈ `feat|fix|refactor|docs|test|chore`.
- Ejemplo: `feat(ndvi): agrega lectura por ventana del AOI`.

## Branches

- `main` es la rama por defecto del repositorio publicado; features en `feat/<nombre-corto>`, fixes en `fix/<nombre-corto>`.

## Idioma

- Código (módulos, funciones, variables): inglés.
- Comentarios, docstrings, documentación y mensajes de commit: español.
