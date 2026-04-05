# Documentación de API: Gestión de Empresas

Documentación completa de los endpoints de empresas para consumo del frontend.

---

## 1. Registro de Empresa + Usuario Principal

Crea una nueva empresa con su suscripción y registra al administrador inicial.

- **Endpoint:** `POST /admin/empresas`
- **Autenticación:** Bearer Token de Administrador del Sistema.

### Request Body

| Campo | Tipo | Requerido | Descripción |
| :--- | :--- | :--- | :--- |
| `empresa.nombre` | string | **Sí** | Nombre comercial (max 150 caracteres). |
| `empresa.pais` | string | No | Código ISO 2 letras (ej: `"PE"`, `"CO"`). |
| `empresa.moneda` | string | No | Código ISO 3 letras (ej: `"PEN"`, `"USD"`). Se deduce por país si no se envía. |
| `empresa.suscripcion_dias` | number | No | Días de suscripción (ej: `30`, `365`). Si es `0`, no se asigna vencimiento. |
| `usuario.usuario` | string | **Sí** | Nombre de usuario del administrador de la empresa. |
| `usuario.password` | string | **Sí** | Contraseña del administrador. |

### Ejemplo Request
```json
{
  "empresa": {
    "nombre": "Inmobiliaria Global S.A.",
    "pais": "PE",
    "moneda": "PEN",
    "suscripcion_dias": 365
  },
  "usuario": {
    "usuario": "admin_global",
    "password": "Password123!"
  }
}
```

### Response (201 Created)
```json
{
  "empresa_id": 5,
  "usuario_id": 12
}
```

---

## 2. Listado Paginado de Empresas

Retorna las empresas de forma paginada (máx 10 por página) con campos resumidos. Permite buscar por nombre.

- **Endpoint:** `GET /admin/empresas`
- **Autenticación:** Bearer Token.

### Query Params

| Parámetro | Tipo | Requerido | Descripción |
| :--- | :--- | :--- | :--- |
| `pagina` | number | No | Número de página (default: `1`). |
| `busqueda` | string | No | Texto parcial para filtrar por nombre. |

### Ejemplo Request
```
GET /admin/empresas?pagina=1&busqueda=inmobiliaria
```

### Response (200 OK)
```json
{
  "datos": [
    {
      "id": 1,
      "nombre": "Inmobiliaria Global S.A.",
      "pais": "PE",
      "estado": true,
      "vencimiento": "2027-04-05T15:00:34Z",
      "creado_en": "2026-04-05T15:00:34Z"
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

### Campos de `paginacion`

| Campo | Descripción |
| :--- | :--- |
| `total` | Total de empresas que coinciden con la búsqueda. |
| `paginas` | Cantidad total de páginas disponibles. |
| `pagina` | Página actual. |
| `por_pagina` | Registros por página (fijo: 10). |

---

## 3. Detalle Completo de Empresa

Devuelve **todos** los datos de una empresa: información general, moneda con configuración de render para el frontend, suscripción, estado y fechas.

- **Endpoint:** `GET /admin/empresas/{id}/detalle`
- **Autenticación:** Bearer Token.

### Ejemplo Request
```
GET /admin/empresas/1/detalle
```

### Response (200 OK)
```json
{
  "id": 1,
  "nombre": "Inmobiliaria Global S.A.",
  "pais": "PE",
  "moneda": "PEN",
  "moneda_info": {
    "codigo": "PEN",
    "decimales": 2,
    "incremento": 1,
    "render": {
      "metodo": "Intl.NumberFormat",
      "currency": "PEN",
      "minimum_fraction_digits": 2,
      "maximum_fraction_digits": 2
    }
  },
  "maximo_usuarios": 1,
  "estado": true,
  "vencimiento": "2027-04-05T15:00:34Z",
  "creado_en": "2026-04-05T15:00:34Z"
}
```

### Campos del Response

| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID único de la empresa. |
| `nombre` | string | Nombre comercial. |
| `pais` | string | Código ISO 2 letras del país. |
| `moneda` | string | Código ISO 3 letras de la moneda. |
| `moneda_info` | object | Configuración completa de la moneda (ver abajo). |
| `maximo_usuarios` | number | Límite de usuarios permitidos. |
| `estado` | boolean | `true` = activa, `false` = inactiva. |
| `vencimiento` | string (ISO 8601) | Fecha de expiración de la suscripción en UTC. Puede estar vacío si no se asignó. |
| `creado_en` | string (ISO 8601) | Fecha de creación en UTC. |

### Campos de `moneda_info`

| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `codigo` | string | Código de la moneda (`"PEN"`, `"USD"`, etc). |
| `decimales` | number | Cantidad de decimales (ej: `2`). |
| `incremento` | number | Incremento mínimo. |
| `render.metodo` | string | Método de renderizado sugerido (`"Intl.NumberFormat"`). |
| `render.currency` | string | Código de moneda para `Intl.NumberFormat`. |
| `render.minimum_fraction_digits` | number | Decimales mínimos a mostrar. |
| `render.maximum_fraction_digits` | number | Decimales máximos a mostrar. |

### Ejemplo de uso en Frontend (JavaScript)
```javascript
// Formatear montos usando moneda_info
const { moneda_info } = empresa;
const formatter = new Intl.NumberFormat('es-PE', {
  style: 'currency',
  currency: moneda_info.render.currency,
  minimumFractionDigits: moneda_info.render.minimum_fraction_digits,
  maximumFractionDigits: moneda_info.render.maximum_fraction_digits,
});

console.log(formatter.format(1500.50)); // S/ 1,500.50
```

---

## Notas Generales para el Frontend

1. **Estado Booleano:** `estado` es `true` (activa) o `false` (inactiva). No es un string.
2. **Fechas UTC:** Todas las fechas vienen en formato ISO 8601 con zona horaria UTC (`Z`). Conviértelas a la zona local del usuario para mostrarlas.
3. **Suscripción:** Si `suscripcion_dias` fue `0` al registrar, `vencimiento` estará vacío (acceso sin límite de tiempo).
4. **Paginación:** La lista de empresas siempre viene envuelta en `datos` + `paginacion`. Nunca es un array suelto.
