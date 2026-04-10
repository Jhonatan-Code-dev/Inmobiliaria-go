# API de Gastos (Módulo de Egresos Simplificado)

Este módulo permite gestionar los gastos de la empresa, visualizar registros históricos con filtros avanzados y paginación.

**Base URL:** `http://localhost:5000/api/user/gastos`

## Listar Gastos
Obtiene una lista paginada de gastos de una empresa específica.

- **Endpoint completo:** `GET /api/user/gastos`
- **Auth:** Requiere Bearer Token o Cookie `token_usuario`.
- **Requiere:** enviar `empresa_id` por query string.
- **Paginación fija del backend:** `10` registros por página.
- **Orden de resultados:** más recientes primero por `fecha`, y luego por `id` descendente.

### Headers recomendados
```http
Authorization: Bearer TU_TOKEN
Content-Type: application/json
```

### Query Params soportados
| Parámetro | Tipo | Obligatorio | Descripción |
|---|---|---|---|
| `empresa_id` | int | Sí | ID de la empresa a consultar. |
| `pag` | int | No | Número de página. Valor por defecto: `1`. |
| `anio` | int | No | Filtra por año. Ejemplo: `2026`. |
| `mes` | int | No | Filtra por mes del `1` al `12`. Normalmente se usa junto con `anio`. |
| `desde` | string | No | Fecha inicial en formato `YYYY-MM-DD`. |
| `hasta` | string | No | Fecha final en formato `YYYY-MM-DD`. |
| `fecha` | string | No | Fecha exacta en formato `YYYY-MM-DD`. Si se envía, tiene prioridad sobre los demás filtros de fecha. |

### Reglas de filtrado
- `empresa_id` es obligatorio.
- Si envías `fecha`, el backend busca solo ese día.
- Si envías `desde` y `hasta`, el backend filtra por rango.
- Si envías `anio`, el backend filtra todo ese año.
- Si envías `anio` y `mes`, el backend filtra solo ese mes.
- Si no envías más filtros además de `empresa_id`, devuelve todos los gastos de esa empresa paginados.
- Si el token del usuario ya está asociado a una empresa, el backend valida que coincida con el `empresa_id` enviado.

### Importante sobre búsqueda
Actualmente este endpoint **no soporta búsqueda por texto** como `search`, `q`, `descripcion` o similares.

Si el frontend necesita búsqueda por descripción o por método de pago, hay que agregarlo en backend primero.

### Ejemplos de consumo

#### 1. Listar primera página
```http
GET /api/user/gastos
```

#### 2. Listar página 2
```http
GET /api/user/gastos?empresa_id=1&pag=2
```

#### 3. Filtrar por año y mes
```http
GET /api/user/gastos?empresa_id=1&anio=2026&mes=4
```

#### 4. Filtrar por rango de fechas
```http
GET /api/user/gastos?empresa_id=1&desde=2026-04-01&hasta=2026-04-30
```

#### 5. Filtrar por fecha exacta
```http
GET /api/user/gastos?empresa_id=1&fecha=2026-04-08
```

### Respuesta (200 OK)
```json
{
  "datos": [
    {
      "id": 1,
      "empresa_id": 1,
      "monto": 50.00,
      "fecha": "2024-03-10T00:00:00Z",
      "tipo_pago_id": 1,
      "descripcion": "Reparación de tubería"
    }
  ],
  "paginacion": {
    "total": 1,
    "paginas": 1,
    "pagina": 1,
    "por_pagina": 10
  }
}
```

### Estructura de la respuesta
| Campo | Tipo | Descripción |
|---|---|---|
| `datos` | array | Lista de gastos de la página solicitada. |
| `datos[].id` | int | ID del gasto. |
| `datos[].empresa_id` | int | ID de la empresa consultada. |
| `datos[].monto` | number | Monto decimal del gasto. |
| `datos[].fecha` | string | Fecha/hora del gasto en formato ISO 8601. |
| `datos[].tipo_pago_id` | int | ID del tipo de pago. |
| `datos[].descripcion` | string | Descripción del gasto. |
| `paginacion.total` | int | Total de registros encontrados con los filtros aplicados. |
| `paginacion.paginas` | int | Total de páginas disponibles. |
| `paginacion.pagina` | int | Página actual devuelta por la API. |
| `paginacion.por_pagina` | int | Cantidad de registros por página. Actualmente siempre `10`. |

### Ejemplo de respuesta vacía
```json
{
  "datos": [],
  "paginacion": {
    "total": 0,
    "paginas": 0,
    "pagina": 1,
    "por_pagina": 10
  }
}
```

### Errores posibles

#### 401 Unauthorized
Cuando no envías token válido o la cookie de sesión no existe.

```json
{
  "message": "Unauthorized"
}
```

#### 400 Bad Request
Cuando no envías `empresa_id` o envías un valor inválido.

```json
{
  "message": "empresa_id es obligatorio"
}
```

#### 403 Forbidden
Cuando el `empresa_id` enviado no coincide con la empresa asociada al token.

```json
{
  "message": "empresa_id no coincide con la sesión"
}
```

#### 500 Internal Server Error
Cuando ocurre un error interno al consultar los gastos.

```json
{
  "message": "detalle del error"
}
```

### Ejemplo para frontend en JavaScript
```js
async function listarGastos({ token, empresaId, pag = 1, anio, mes, desde, hasta, fecha }) {
  const params = new URLSearchParams();

  params.set("empresa_id", String(empresaId));
  if (pag) params.set("pag", String(pag));
  if (anio) params.set("anio", String(anio));
  if (mes) params.set("mes", String(mes));
  if (desde) params.set("desde", desde);
  if (hasta) params.set("hasta", hasta);
  if (fecha) params.set("fecha", fecha);

  const res = await fetch(`http://localhost:5000/api/user/gastos?${params.toString()}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    }
  });

  if (!res.ok) {
    throw await res.json();
  }

  return res.json();
}
```

### Contrato recomendado para frontend
- Siempre enviar `empresa_id`.
- Usar `paginacion.total` para saber si hay resultados.
- Usar `paginacion.paginas` para construir el paginador.
- Usar `paginacion.pagina` como fuente de verdad de la página actual.
- No asumir más de `10` elementos por respuesta.
- No enviar parámetros de búsqueda por texto porque hoy no serán aplicados por el backend.

---

## Registrar Gasto
Crea un nuevo registro de gasto y genera directamente el registro en caja.

- **URL:** `POST /`
- **Body:**
```json
{
  "empresa_id": 1,
  "monto": 50.50,
  "fecha": "2024-03-01",
  "tipo_pago_id": 3,
  "descripcion": "Pago mensual Starlink"
}
```
*Nota: El monto se envía y se recibe como número decimal (ej. 50.50).*

### Respuesta (201 Created)
Retorna el objeto creado con su ID.

---

## Listar Tipos de Pago
Obtiene la lista de los métodos de pago disponibles en la base de datos.

- **URL:** `GET /tipos-pago`
- **Respuesta:**
```json
[
  {"id": 1, "nombre": "efectivo"},
  {"id": 2, "nombre": "transferencia"},
  {"id": 3, "nombre": "yape"}
]
```

---

## Actualizar Gasto
Actualiza los datos de un gasto existente.

- **URL:** `PUT /{id}`
- **Body:** Similar al de creación.

---

## Eliminar Gasto
Elimina físicamente un gasto.

- **Endpoint completo:** `DELETE /api/user/gastos/{id}?empresa_id=1`
- **Auth:** Requiere Bearer Token o Cookie `token_usuario`.
- **Requiere:** `id` en la ruta y `empresa_id` en query string.

### Parámetros
| Parámetro | Ubicación | Tipo | Obligatorio | Descripción |
|---|---|---|---|---|
| `id` | path | int | Sí | ID del gasto a eliminar. |
| `empresa_id` | query | int | Sí | ID de la empresa propietaria del gasto. |

### Ejemplo de request
```http
DELETE /api/user/gastos/17?empresa_id=1
Authorization: Bearer TU_TOKEN
```

### Respuesta exitosa
```json
{
  "message": "gasto eliminado"
}
```

### Errores posibles

#### 400 Bad Request
Cuando falta `empresa_id` o el `id` es inválido.

```json
{
  "message": "empresa_id es obligatorio"
}
```

#### 401 Unauthorized
Cuando no envías token válido o cookie de sesión.

```json
{
  "message": "Unauthorized"
}
```

#### 403 Forbidden
Cuando el gasto no pertenece a la empresa indicada o el `empresa_id` no coincide con la sesión.

```json
{
  "message": "no autorizado para eliminar este gasto"
}
```

### Regla para frontend
- Debes enviar el mismo `empresa_id` que corresponde a la sesión del usuario.
- Si intentas borrar un gasto de otra empresa, el backend responderá `403`.
