# API de Clientes

Módulo completo para frontend: catálogo de tipos de identificación, listado con búsqueda y paginación, detalle, registro, edición y eliminación.

**Base URL:** `http://localhost:5000/api/user/clientes`

## Reglas clave para frontend

- Todos los endpoints requieren `Authorization: Bearer <token>`.
- El token del usuario contiene `empresa_id`; el backend valida que coincida con la empresa enviada.
- El listado devuelve máximo `10` registros por página.
- Antes de crear o editar un cliente, frontend debe consultar el catálogo de tipos de identificación y usar un `tipo_identificacion_id` existente.
- Todos los errores ahora deben responder con formato JSON:

```json
{
  "message": "detalle del error"
}
```

## Flujo recomendado de consumo

1. Llamar `GET /api/user/clientes/tipos-identificacion`.
2. Guardar `id`, `codigo` y `nombre` para poblar el select del formulario.
3. Crear o editar enviando uno de esos `id` en `tipo_identificacion_id`.
4. Usar `GET /api/user/clientes` para la tabla con búsqueda y paginación.
5. Usar `GET /api/user/clientes/{id}` para precargar el formulario de edición.

## Headers

```http
Authorization: Bearer TU_TOKEN
Content-Type: application/json
Accept: application/json
```

## 1. Catálogo de tipos de identificación

- **Endpoint:** `GET /api/user/clientes/tipos-identificacion`
- **Objetivo:** obtener IDs válidos para `tipo_identificacion_id`.

### Respuesta `200`

```json
[
  {
    "id": 2,
    "codigo": "CE",
    "nombre": "Carnet de Extranjeria",
    "pais": "PE",
    "activo": true
  },
  {
    "id": 1,
    "codigo": "DNI",
    "nombre": "Documento Nacional de Identidad",
    "pais": "PE",
    "activo": true
  },
  {
    "id": 3,
    "codigo": "PAS",
    "nombre": "Pasaporte",
    "pais": null,
    "activo": true
  },
  {
    "id": 4,
    "codigo": "RUC",
    "nombre": "Registro Unico de Contribuyentes",
    "pais": "PE",
    "activo": true
  }
]
```

### Regla importante

Si frontend envía un `tipo_identificacion_id` que no aparece aquí, el backend responderá:

```json
{
  "message": "tipo_identificacion_id no existe o está inactivo"
}
```

## 2. Listar clientes

- **Endpoint:** `GET /api/user/clientes?empresa_id=1&pag=1&buscar=juan`
- **Uso típico:** tabla principal.

### Query params

| Parámetro | Tipo | Obligatorio | Descripción |
|---|---|---|---|
| `empresa_id` | int | Sí | Empresa propietaria de los clientes. |
| `pag` | int | No | Página actual. Default `1`. |
| `buscar` | string | No | Busca por `nombres`, `apellidos` o `documento_numero`. |

### Respuesta `200`

```json
{
  "datos": [
    {
      "id": 3,
      "empresa_id": 1,
      "tipo_identificacion_id": 1,
      "documento_numero": "90000002",
      "nombres": "Cliente Prueba",
      "apellidos": "Frontend Demo",
      "correo": "cliente.prueba@example.com",
      "fecha_nacimiento": "1995-06-18T00:00:00Z",
      "nacionalidad": "Peruana",
      "direccion": "Av. Lima 123",
      "contacto_emergencia": "Maria Demo",
      "telefono_emergencia": "999888777",
      "notas": "Creado en auditoria",
      "estado": "activo",
      "creado_en": "2026-04-10T03:38:55Z"
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

### Respuesta vacía

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

## 3. Obtener cliente por ID

- **Endpoint:** `GET /api/user/clientes/{id}?empresa_id=1`
- **Uso típico:** precargar edición.

### Respuesta `200`

```json
{
  "id": 3,
  "empresa_id": 1,
  "tipo_identificacion_id": 1,
  "documento_numero": "90000002",
  "nombres": "Cliente Prueba",
  "apellidos": "Frontend Demo",
  "correo": "cliente.prueba@example.com",
  "fecha_nacimiento": "1995-06-18T00:00:00Z",
  "nacionalidad": "Peruana",
  "direccion": "Av. Lima 123",
  "contacto_emergencia": "Maria Demo",
  "telefono_emergencia": "999888777",
  "notas": "Creado en auditoria",
  "estado": "activo",
  "creado_en": "2026-04-10T03:38:55Z"
}
```

## 4. Crear cliente

- **Endpoint:** `POST /api/user/clientes`

### Body mínimo válido

```json
{
  "empresa_id": 1,
  "tipo_identificacion_id": 1,
  "documento_numero": "90000002",
  "nombres": "Cliente Prueba",
  "estado": "activo"
}
```

### Body completo válido

```json
{
  "empresa_id": 1,
  "tipo_identificacion_id": 1,
  "documento_numero": "90000002",
  "nombres": "Cliente Prueba",
  "apellidos": "Frontend Demo",
  "correo": "cliente.prueba@example.com",
  "fecha_nacimiento": "1995-06-18",
  "nacionalidad": "Peruana",
  "direccion": "Av. Lima 123",
  "contacto_emergencia": "Maria Demo",
  "telefono_emergencia": "999888777",
  "notas": "Creado en auditoria",
  "estado": "activo"
}
```

### Respuesta `201`

Retorna el cliente creado:

```json
{
  "id": 3,
  "empresa_id": 1,
  "tipo_identificacion_id": 1,
  "documento_numero": "90000002",
  "nombres": "Cliente Prueba",
  "apellidos": "Frontend Demo",
  "correo": "cliente.prueba@example.com",
  "fecha_nacimiento": "1995-06-18T00:00:00Z",
  "nacionalidad": "Peruana",
  "direccion": "Av. Lima 123",
  "contacto_emergencia": "Maria Demo",
  "telefono_emergencia": "999888777",
  "notas": "Creado en auditoria",
  "estado": "activo",
  "creado_en": "2026-04-10T03:38:55.133895833Z"
}
```

## 5. Actualizar cliente

- **Endpoint:** `PUT /api/user/clientes/{id}`

### Body

```json
{
  "empresa_id": 1,
  "tipo_identificacion_id": 2,
  "documento_numero": "90000002",
  "nombres": "Cliente Prueba Editado",
  "apellidos": "Frontend Demo Editado",
  "correo": "cliente.editado@example.com",
  "fecha_nacimiento": "1995-06-19",
  "nacionalidad": "Peruana",
  "direccion": "Av. Lima 456",
  "contacto_emergencia": "Maria Editada",
  "telefono_emergencia": "999888776",
  "notas": "Actualizado en auditoria",
  "estado": "inactivo"
}
```

### Respuesta `200`

```json
{
  "id": 3,
  "empresa_id": 1,
  "tipo_identificacion_id": 2,
  "documento_numero": "90000002",
  "nombres": "Cliente Prueba Editado",
  "apellidos": "Frontend Demo Editado",
  "correo": "cliente.editado@example.com",
  "fecha_nacimiento": "1995-06-19T00:00:00Z",
  "nacionalidad": "Peruana",
  "direccion": "Av. Lima 456",
  "contacto_emergencia": "Maria Editada",
  "telefono_emergencia": "999888776",
  "notas": "Actualizado en auditoria",
  "estado": "inactivo",
  "creado_en": "2026-04-10T03:38:55Z"
}
```

## 6. Eliminar cliente

- **Endpoint:** `DELETE /api/user/clientes/{id}?empresa_id=1`

### Respuesta `200`

```json
{
  "message": "cliente eliminado"
}
```

## Errores que frontend debe manejar

### `400 Bad Request`

```json
{
  "message": "empresa_id es obligatorio"
}
```

```json
{
  "message": "fecha_nacimiento debe tener formato YYYY-MM-DD"
}
```

```json
{
  "message": "tipo_identificacion_id no existe o está inactivo"
}
```

### `403 Forbidden`

```json
{
  "message": "empresa_id no coincide con la sesión"
}
```

```json
{
  "message": "no autorizado para actualizar este cliente"
}
```

### `404 Not Found`

```json
{
  "message": "cliente no encontrado"
}
```

## Instrucciones directas para frontend

- No hardcodear `tipo_identificacion_id`.
- Cargar el select desde `GET /api/user/clientes/tipos-identificacion`.
- En tabla, usar `GET /api/user/clientes` con `empresa_id`, `pag` y `buscar`.
- Para editar, primero pedir `GET /api/user/clientes/{id}?empresa_id=...`.
- En create y update, enviar `fecha_nacimiento` en formato `YYYY-MM-DD`.
- En delete y detail, siempre enviar `empresa_id` por query string.
- Usar siempre `response.message` para mostrar errores del backend.

## Ejemplos para frontend

### Obtener catálogo

```js
async function listarTiposIdentificacion(token) {
  const res = await fetch("http://localhost:5000/api/user/clientes/tipos-identificacion", {
    headers: {
      "Authorization": `Bearer ${token}`,
      "Accept": "application/json"
    }
  });

  if (!res.ok) throw await res.json();
  return res.json();
}
```

### Listar clientes

```js
async function listarClientes({ token, empresaId, pag = 1, buscar = "" }) {
  const params = new URLSearchParams({
    empresa_id: String(empresaId),
    pag: String(pag)
  });

  if (buscar.trim()) params.set("buscar", buscar.trim());

  const res = await fetch(`http://localhost:5000/api/user/clientes?${params.toString()}`, {
    headers: {
      "Authorization": `Bearer ${token}`,
      "Accept": "application/json"
    }
  });

  if (!res.ok) throw await res.json();
  return res.json();
}
```

### Crear cliente

```js
async function crearCliente({ token, payload }) {
  const res = await fetch("http://localhost:5000/api/user/clientes", {
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
