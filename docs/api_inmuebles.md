# API de Inmuebles

Módulo completo para inmobiliaria: registro de inmuebles, listado con búsqueda y paginación, detalle con unidades y CRUD de unidades.

**Base URL:** `http://localhost:5000/api/user/inmuebles`

## Qué representa este módulo

- `Inmueble` = propiedad principal.
- `Unidad` = ambiente o espacio alquilable dentro del inmueble.

Ejemplos:
- Un edificio es un `inmueble`.
- Cada departamento o cuarto dentro del edificio es una `unidad`.

## Reglas para frontend

- Todos los endpoints requieren `Authorization: Bearer <token>`.
- El backend valida `empresa_id` contra la empresa del token.
- El listado devuelve máximo `10` inmuebles por página.
- El detalle del inmueble ya puede devolver sus `unidades`.
- Los montos `precio_base` y `deposito_requerido` se envían y reciben como número decimal.
- Todos los errores responden así:

```json
{
  "message": "detalle del error"
}
```

## Flujo recomendado para frontend

1. Listar inmuebles con `GET /api/user/inmuebles`.
2. Crear inmueble con `POST /api/user/inmuebles`.
3. Entrar al detalle con `GET /api/user/inmuebles/{id}?empresa_id=...`.
4. Crear unidades dentro del inmueble con `POST /api/user/inmuebles/{id}/unidades`.
5. Editar unidad o inmueble según necesidad.

## 1. Listar inmuebles

- **Endpoint:** `GET /api/user/inmuebles?empresa_id=1&pag=1&buscar=central`

### Query params

| Parámetro | Tipo | Obligatorio | Descripción |
|---|---|---|---|
| `empresa_id` | int | Sí | Empresa dueña de los inmuebles. |
| `pag` | int | No | Página actual. Default `1`. |
| `buscar` | string | No | Busca por `nombre`, `direccion` o `ciudad`. |
| `estado` | string | No | Filtra por estado: `activa`, `mantenimiento`, `inactiva`. |
| `tipo` | string | No | Filtra por tipo: `casa`, `edificio`, `quinta`, `condominio`, `otro`. |

### Respuesta `200`

```json
{
  "datos": [
    {
      "id": 1,
      "empresa_id": 1,
      "nombre": "Edificio Central",
      "tipo": "edificio",
      "descripcion": "Edificio principal",
      "direccion": "Av. Principal 123",
      "ciudad": "Lima",
      "region": "Lima",
      "pais": "PE",
      "codigo_postal": "15001",
      "total_pisos": 5,
      "total_unidades": 12,
      "estado": "activa",
      "creado_en": "2026-04-10T04:00:00Z"
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

## 2. Obtener detalle de inmueble

- **Endpoint:** `GET /api/user/inmuebles/{id}?empresa_id=1`

### Respuesta `200`

```json
{
  "id": 1,
  "empresa_id": 1,
  "nombre": "Edificio Central",
  "tipo": "edificio",
  "descripcion": "Edificio principal",
  "direccion": "Av. Principal 123",
  "ciudad": "Lima",
  "region": "Lima",
  "pais": "PE",
  "codigo_postal": "15001",
  "total_pisos": 5,
  "total_unidades": 12,
  "estado": "activa",
  "creado_en": "2026-04-10T04:00:00Z",
  "unidades": [
    {
      "id": 2,
      "propiedad_id": 1,
      "codigo": "A-101",
      "nombre": "Departamento 101",
      "tipo": "departamento",
      "numero_piso": 1,
      "dormitorios": 2,
      "banos": 1,
      "area_m2": 58.5,
      "capacidad": 3,
      "moneda": "PEN",
      "precio_base": 850.00,
      "deposito_requerido": 500.00,
      "incluye_agua": true,
      "incluye_luz": false,
      "incluye_internet": true,
      "notas": "Vista a la calle",
      "estado": "disponible",
      "creado_en": "2026-04-10T04:10:00Z"
    }
  ]
}
```

## 3. Crear inmueble

- **Endpoint:** `POST /api/user/inmuebles`

### Body

```json
{
  "empresa_id": 1,
  "nombre": "Edificio Central",
  "tipo": "edificio",
  "descripcion": "Edificio principal",
  "direccion": "Av. Principal 123",
  "ciudad": "Lima",
  "region": "Lima",
  "pais": "PE",
  "codigo_postal": "15001",
  "total_pisos": 5,
  "total_unidades": 12,
  "estado": "activa"
}
```

### Respuesta `201`

Retorna el inmueble creado.

## 4. Actualizar inmueble

- **Endpoint:** `PUT /api/user/inmuebles/{id}`

### Body

```json
{
  "empresa_id": 1,
  "nombre": "Edificio Central Remodelado",
  "tipo": "edificio",
  "descripcion": "Edificio actualizado",
  "direccion": "Av. Principal 456",
  "ciudad": "Lima",
  "region": "Lima",
  "pais": "PE",
  "codigo_postal": "15001",
  "total_pisos": 6,
  "total_unidades": 14,
  "estado": "mantenimiento"
}
```

## 5. Eliminar inmueble

- **Endpoint:** `DELETE /api/user/inmuebles/{id}?empresa_id=1`

### Respuesta `200`

```json
{
  "message": "inmueble eliminado"
}
```

## 6. Listar unidades de un inmueble

- **Endpoint:** `GET /api/user/inmuebles/{id}/unidades?empresa_id=1`

### Respuesta `200`

```json
[
  {
    "id": 2,
    "propiedad_id": 1,
    "codigo": "A-101",
    "nombre": "Departamento 101",
    "tipo": "departamento",
    "numero_piso": 1,
    "dormitorios": 2,
    "banos": 1,
    "area_m2": 58.5,
    "capacidad": 3,
    "moneda": "PEN",
    "precio_base": 850.00,
    "deposito_requerido": 500.00,
    "incluye_agua": true,
    "incluye_luz": false,
    "incluye_internet": true,
    "notas": "Vista a la calle",
    "estado": "disponible",
    "creado_en": "2026-04-10T04:10:00Z"
  }
]
```

## 7. Crear unidad

- **Endpoint:** `POST /api/user/inmuebles/{id}/unidades`

### Body

```json
{
  "codigo": "A-101",
  "nombre": "Departamento 101",
  "tipo": "departamento",
  "numero_piso": 1,
  "dormitorios": 2,
  "banos": 1,
  "area_m2": 58.5,
  "capacidad": 3,
  "moneda": "PEN",
  "precio_base": 850.00,
  "deposito_requerido": 500.00,
  "incluye_agua": true,
  "incluye_luz": false,
  "incluye_internet": true,
  "notas": "Vista a la calle",
  "estado": "disponible"
}
```

### Respuesta `201`

Retorna la unidad creada.

## 8. Obtener unidad

- **Endpoint:** `GET /api/user/inmuebles/{id}/unidades/{unidadId}?empresa_id=1`

## 9. Actualizar unidad

- **Endpoint:** `PUT /api/user/inmuebles/{id}/unidades/{unidadId}`

### Body

```json
{
  "codigo": "A-101",
  "nombre": "Departamento 101 Remodelado",
  "tipo": "departamento",
  "numero_piso": 1,
  "dormitorios": 3,
  "banos": 2,
  "area_m2": 62.0,
  "capacidad": 4,
  "moneda": "PEN",
  "precio_base": 920.00,
  "deposito_requerido": 600.00,
  "incluye_agua": true,
  "incluye_luz": true,
  "incluye_internet": true,
  "notas": "Unidad actualizada",
  "estado": "reservado"
}
```

## 10. Eliminar unidad

- **Endpoint:** `DELETE /api/user/inmuebles/{id}/unidades/{unidadId}?empresa_id=1`

### Respuesta `200`

```json
{
  "message": "unidad eliminada"
}
```

## Errores esperados

### `400`

```json
{
  "message": "empresa_id es obligatorio"
}
```

```json
{
  "message": "ID inválido"
}
```

### `403`

```json
{
  "message": "empresa_id no coincide con la sesión"
}
```

### `404`

```json
{
  "message": "inmueble no encontrado"
}
```

```json
{
  "message": "unidad no encontrada"
}
```

## Instrucciones directas para frontend

- Usar el endpoint de listado para tabla principal.
- Usar el detalle del inmueble para pantalla de administración de unidades.
- En delete y detail enviar `empresa_id` por query string.
- Para precios, enviar números decimales como `850.00`.
- Si el formulario maneja moneda, hoy el backend acepta el código como texto, por ejemplo `PEN`.
- Mostrar siempre `response.message` en errores.

## Ejemplos rápidos para frontend

### Listar inmuebles

```js
async function listarInmuebles({ token, empresaId, pag = 1, buscar = "" }) {
  const params = new URLSearchParams({
    empresa_id: String(empresaId),
    pag: String(pag)
  });

  if (buscar.trim()) params.set("buscar", buscar.trim());

  const res = await fetch(`http://localhost:5000/api/user/inmuebles?${params.toString()}`, {
    headers: {
      "Authorization": `Bearer ${token}`,
      "Accept": "application/json"
    }
  });

  if (!res.ok) throw await res.json();
  return res.json();
}
```

### Crear inmueble

```js
async function crearInmueble({ token, payload }) {
  const res = await fetch("http://localhost:5000/api/user/inmuebles", {
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

### Crear unidad

```js
async function crearUnidad({ token, inmuebleId, payload }) {
  const res = await fetch(`http://localhost:5000/api/user/inmuebles/${inmuebleId}/unidades`, {
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
