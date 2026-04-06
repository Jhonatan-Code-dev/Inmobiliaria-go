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
| `maximo_usuarios` | number | Cantidad máxima de usuarios permitidos en la empresa. |
| `estado` | boolean | `true` = activa, `false` = inactiva. |
| `vencimiento` | string | Fecha de expiración de la suscripción en formato ISO 8601 UTC. Vacío si no tiene límite. |
| `creado_en` | string | Fecha de creación en formato ISO 8601 UTC. |



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

// Verificar si la suscripción está vencida
if (empresa.vencimiento) {
  const vence = new Date(empresa.vencimiento);
  const hoy = new Date();
  if (vence < hoy) {
    console.log('Suscripción vencida');
  }
}
```
