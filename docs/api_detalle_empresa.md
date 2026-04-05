# Endpoint: Detalle Completo de Empresa

Devuelve todos los datos de una empresa por su ID.

---

## Request

```
GET /admin/empresas/{id}/detalle
```

**Autenticación:** Bearer Token (Header `Authorization: Bearer <token>`)

### Parámetros URL

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID de la empresa |

### Ejemplo
```
GET /admin/empresas/1/detalle
```

---

## Response (200 OK)

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

---

## Campos del Response

| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID único de la empresa. |
| `nombre` | string | Nombre comercial de la empresa. |
| `pais` | string | Código ISO 2 letras del país (ej: `"PE"`, `"CO"`). |
| `moneda` | string | Código ISO 3 letras de la moneda (ej: `"PEN"`, `"USD"`). |
| `moneda_info` | object | Configuración completa de la moneda (ver tabla abajo). |
| `maximo_usuarios` | number | Cantidad máxima de usuarios permitidos en la empresa. |
| `estado` | boolean | `true` = activa, `false` = inactiva. |
| `vencimiento` | string | Fecha de expiración de la suscripción en formato ISO 8601 UTC. Vacío si no tiene límite. |
| `creado_en` | string | Fecha de creación en formato ISO 8601 UTC. |

### Campos de `moneda_info`

| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `codigo` | string | Código de la moneda (`"PEN"`, `"USD"`, etc). |
| `decimales` | number | Cantidad de decimales de la moneda (ej: `2`). |
| `incremento` | number | Incremento mínimo. |
| `render.metodo` | string | Método sugerido para formatear (`"Intl.NumberFormat"`). |
| `render.currency` | string | Código para pasar a `Intl.NumberFormat`. |
| `render.minimum_fraction_digits` | number | Decimales mínimos a mostrar. |
| `render.maximum_fraction_digits` | number | Decimales máximos a mostrar. |

---

## Errores

| Código | Mensaje | Causa |
| :--- | :--- | :--- |
| 400 | Bad Request | El `id` no es un número válido. |
| 401 | Unauthorized | Token no enviado o expirado. |
| 500 | Internal Server Error | Error interno del servidor. |

---

## Ejemplo de uso en Frontend

```javascript
const response = await fetch('/admin/empresas/1/detalle', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const empresa = await response.json();

// Formatear montos con la moneda de la empresa
const { moneda_info } = empresa;
const formatter = new Intl.NumberFormat('es-PE', {
  style: 'currency',
  currency: moneda_info.render.currency,
  minimumFractionDigits: moneda_info.render.minimum_fraction_digits,
  maximumFractionDigits: moneda_info.render.maximum_fraction_digits,
});

console.log(formatter.format(1500.50)); // S/ 1,500.50

// Verificar si la suscripción está vencida
if (empresa.vencimiento) {
  const vence = new Date(empresa.vencimiento);
  const hoy = new Date();
  if (vence < hoy) {
    console.log('Suscripción vencida');
  }
}
```
