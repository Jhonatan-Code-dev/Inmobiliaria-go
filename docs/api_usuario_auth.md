# Documentación de API: Autenticación de Usuarios

Documentación de los endpoints de autenticación y perfil para usuarios de empresas (tenants).

---

## 1. Inicio de Sesión (Login)

Permite a un usuario autenticarse en el sistema. Al tener éxito, el servidor devuelve un token JWT y establece una cookie de sesión `token_usuario` (HTTP-only).

- **Endpoint:** `POST /auth/login`
- **Autenticación:** Ninguna (Público).

### Request Body

| Campo | Tipo | Requerido | Descripción |
| :--- | :--- | :--- | :--- |
| `usuario` | string | **Sí** | Nombre de usuario. |
| `contrasena` | string | **Sí** | Contraseña del usuario. |

### Ejemplo Request
```json
{
  "usuario": "yona_admin",
  "contrasena": "Password123!"
}
```

### Response (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 12,
    "usuario": "yona_admin",
    "empresa_id": 5
  },
  "empresa": {
    "id": 5,
    "nombre": "Mi Empresa S.A.C.",
    "pais": "PE",
    "moneda": "PEN",
    "maximo_usuarios": 10,
    "estado": true,
    "vencimiento": "2027-04-05T15:00:00Z",
    "creado_en": "2026-04-01T10:00:00Z"
  }
}
```

### Notas del Login
1. **Cookie:** El servidor envía una cookie `token_usuario`. El frontend puede usarla automáticamente si el navegador está configurado para enviar cookies en las peticiones (CORS con `withCredentials: true`).
2. **Token:** Se devuelve el token explícitamente para aplicaciones que prefieran guardarlo en `localStorage` o enviarlo en el header `Authorization: Bearer <token>`.

---

## 2. Cerrar Sesión (Logout)

Invalida la sesión del usuario eliminando la cookie `token_usuario`.

- **Endpoint:** `POST /auth/logout`
- **Autenticación:** Ninguna (Se recomienda estar autenticado, pero limpia la cookie de todos modos).

### Response (200 OK)
```json
{
  "message": "sesión cerrada"
}
```

---

## 3. Obtener Perfil (Me)

Retorna la información del usuario actualmente autenticado basado en el token o cookie enviada.

- **Endpoint:** `GET /me`
- **Autenticación:** Requerida (Bearer Token o Cookie `token_usuario`).

### Ejemplo Request
```
GET /me
Authorization: Bearer <token>
```

### Response (200 OK)
```json
{
  "user": {
    "id": 12,
    "usuario": "yona_admin",
    "empresa_id": 5
  },
  "empresa": {
    "id": 5,
    "nombre": "Mi Empresa S.A.C.",
    "pais": "PE",
    "moneda": "PEN",
    "maximo_usuarios": 10,
    "estado": true,
    "vencimiento": "2027-04-05T15:00:00Z",
    "creado_en": "2026-04-01T10:00:00Z"
  }
}
```

### Response (401 Unauthorized)
```json
{
  "message": "Unauthorized"
}
```

---

## Resumen de Estructuras

### Objeto `user`
| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID único del usuario. |
| `usuario` | string | Nombre de usuario (login). |
| `empresa_id` | number | ID de la empresa a la que pertenece. |

### Objeto `empresa`
| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID único de la empresa. |
| `nombre` | string | Nombre comercial de la empresa. |
| `pais` | string | Código ISO del país. |
| `moneda` | string | Código ISO de la moneda (PEN, USD, etc). |
| `estado` | boolean | Estado de la cuenta de la empresa. |
| `vencimiento` | string | Fecha de expiración de la suscripción (ISO 8601 UTC). |
