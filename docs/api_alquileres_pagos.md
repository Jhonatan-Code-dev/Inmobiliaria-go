# API de Alquileres y Pagos

Módulo para contratos de alquiler y registro de cobros en Alquilamax.

**Base URL Alquileres:** `http://localhost:5000/api/user/alquileres`  
**Base URL Pagos:** `http://localhost:5000/api/user/pagos`

## Reglas globales

- Todos los endpoints requieren `Authorization: Bearer <token>`.
- El backend valida multiempresa con el `empresa_id` del token.
- Los errores responden como:

```json
{
  "message": "detalle del error"
}
```

- Los montos se envían y reciben como números decimales.
- Los listados responden con:

```json
{
  "datos": [],
  "paginacion": {
    "total": 0,
    "paginas": 0,
    "pagina": 1,
    "pagina_actual": 1,
    "por_pagina": 10
  }
}
```

## 1. Listar alquileres

- **Endpoint:** `GET /api/user/alquileres?empresa_id=1&pag=1&por_pagina=10&buscar=QA-301`

### Query params

| Parámetro | Tipo | Obligatorio | Descripción |
|---|---|---|---|
| `empresa_id` | int | Sí | Empresa dueña del contrato. |
| `buscar` | string | No | Busca por nombre de cliente o código de unidad. |
| `estado` | string | No | `activo`, `vencido`, `finalizado`, etc. |
| `unidad_id` | int | No | Filtra por unidad. |
| `pag` | int | No | Página actual. |
| `por_pagina` | int | No | Tamaño de página. Default `10`. |

### Respuesta `200`

```json
{
  "datos": [
    {
      "id": 1,
      "cliente_id": 5,
      "cliente": "Cliente QA Empresa Uno",
      "unidad_id": 3,
      "unidad": "QA-301",
      "monto": 1250.50,
      "fecha_inicio": "2026-04-11",
      "fecha_fin": "2026-12-31",
      "estado": "activo",
      "moneda": "PEN"
    }
  ],
  "paginacion": {
    "total": 1,
    "paginas": 1,
    "pagina": 1,
    "pagina_actual": 1,
    "por_pagina": 5
  }
}
```

## 2. Obtener alquiler por ID

- **Endpoint:** `GET /api/user/alquileres/{id}?empresa_id=1`

## 3. Crear alquiler

- **Endpoint:** `POST /api/user/alquileres`

### Body

```json
{
  "cliente_id": 5,
  "unidad_id": 3,
  "fecha_inicio": "2026-04-10",
  "fecha_fin": "2026-12-31",
  "vencimiento_dia_pago": 5,
  "monto_renta": 1200.50,
  "deposito_garantia": 800.00,
  "moneda": "PEN"
}
```

### Reglas de negocio

- `cliente_id` debe pertenecer a la misma empresa del usuario.
- `unidad_id` debe pertenecer a la misma empresa del usuario.
- La unidad debe estar en estado `disponible`.
- Al crear el contrato, la unidad pasa a `ocupado`.

### Respuesta `201`

```json
{
  "id": 1,
  "cliente_id": 5,
  "cliente": "Cliente QA Empresa Uno",
  "unidad_id": 3,
  "unidad": "QA-301",
  "monto": 1250.50,
  "fecha_inicio": "2026-04-11",
  "fecha_fin": "2026-12-31",
  "estado": "activo",
  "moneda": "PEN"
}
```

## 4. Registrar pago

- **Endpoint:** `POST /api/user/pagos`

### Body

```json
{
  "alquiler_id": 1,
  "monto_pagado": 1250.50,
  "fecha_pago": "2026-04-10",
  "metodo_pago": "transferencia",
  "nota": "Pago abril",
  "mes_correspondiente": 4
}
```

### Respuesta `201`

```json
{
  "id": 1,
  "alquiler_id": 1,
  "numero_recibo": "PAGO-1-1775838730484567952",
  "fecha_pago": "2026-04-10",
  "moneda": "PEN",
  "monto_pagado": 1250.50,
  "metodo_pago": "transferencia",
  "nota": "Pago abril",
  "mes_correspondiente": 4
}
```

## 5. Listar pagos pendientes del mes actual

- **Endpoint:** `GET /api/user/pagos/pendientes?empresa_id=1`

### Respuesta `200`

```json
[
  {
    "alquiler_id": 2,
    "cliente": "Cliente QA Empresa Uno",
    "unidad": "QA-302",
    "monto": 900.00,
    "fecha_vencimiento": "2026-04-05",
    "estado": "activo"
  }
]
```

## Errores importantes

### Cliente de otra empresa

```json
{
  "message": "cliente no pertenece a la empresa"
}
```

### Unidad ya ocupada

```json
{
  "message": "la unidad no está disponible"
}
```

### Pago inválido

```json
{
  "message": "monto_pagado excede el saldo pendiente del alquiler"
}
```

## Instrucciones directas para frontend

- Para crear contrato, primero cargar clientes e inmuebles/unidades disponibles.
- No intentar crear contratos sobre unidades ya ocupadas.
- Usar `buscar` en alquileres para buscar por cliente o código de unidad.
- Usar `GET /api/user/pagos/pendientes` para dashboard de morosidad.
- En el formulario de pago, enviar `mes_correspondiente` del 1 al 12.
- Mostrar siempre `response.message` en cualquier error.

## Ejemplo frontend: crear alquiler

```js
async function crearAlquiler({ token, payload }) {
  const res = await fetch("http://localhost:5000/api/user/alquileres", {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${token}`,
      "Content-Type": "application/json",
      "Accept": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) throw await res.json();
  return res.json();
}
```

## Ejemplo frontend: registrar pago

```js
async function registrarPago({ token, payload }) {
  const res = await fetch("http://localhost:5000/api/user/pagos", {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${token}`,
      "Content-Type": "application/json",
      "Accept": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) throw await res.json();
  return res.json();
}
```
