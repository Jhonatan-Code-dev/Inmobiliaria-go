# Documentación de API: Autenticación de Usuarios (Tenant)

Esta API permite el acceso y gestión de perfil para los usuarios pertenecientes a una empresa (Business Users). El sistema utiliza autenticación basada en JWT, que puede ser manejado mediante una **cookie HTTP-only** o el encabezado **Authorization**.

## Resumen de Endpoints

| Método | Endpoint | Descripción | Acceso |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/login` | Iniciar sesión y obtener token/cookie. | Público |
| `POST` | `/auth/logout` | Cerrar sesión y limpiar cookie. | Autenticado |
| `GET` | `/me` | Obtener datos del usuario y de su empresa. | Autenticado |

---

## 1. Inicio de Sesión (Login)

Autentica a un usuario y establece una sesión.

- **URL:** `/auth/login`
- **Método:** `POST`
- **Headers:** `Content-Type: application/json`

### Body (JSON)

| Campo | Tipo | Requerido | Descripción |
| :--- | :--- | :--- | :--- |
| `usuario` | string | Sí | Nombre de usuario (login). |
| `contrasena` | string | Sí | Contraseña de acceso. |

**Ejemplo de Request:**
```json
{
  "usuario": "yona_admin",
  "contrasena": "Password123!"
}
```

### Respuestas

#### Success (200 OK)
Devuelve el token JWT y los datos básicos del usuario y su empresa. Se establece la cookie `token_usuario` automáticamente.

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

#### Unauthorized (401 Unauthorized)
Credenciales incorrectas o usuario inexistente.

```json
{
  "message": "credenciales inválidas"
}
```

---

## 2. Cerrar Sesión (Logout)

Invalida la sesión del usuario eliminando la cookie del navegador.

- **URL:** `/auth/logout`
- **Método:** `POST`
- **Acceso:** Se recomienda estar autenticado.

### Respuestas

#### Success (200 OK)
```json
{
  "message": "sesión cerrada"
}
```

---

## 3. Obtener Perfil (Me)

Retorna la información del usuario actualmente autenticado y los datos de su empresa.

- **URL:** `/me`
- **Método:** `GET`
- **Headers:** `Authorization: Bearer <token>` (opcional si se usa cookie)

### Respuestas

#### Success (200 OK)
```json
{
  "token": "",
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
*Nota: El campo `token` se devuelve vacío en este endpoint.*

#### Unauthorized (401 Unauthorized)
Token inválido, expirado o falta de sesión.

```json
{
  "message": "Unauthorized"
}
```

---

## Detalles de Estructuras

### Objeto `user`
| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID único del usuario. |
| `usuario` | string | Nombre de usuario. |
| `empresa_id` | number | ID de la empresa vinculada. |

### Objeto `empresa`
| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | number | ID único de la empresa. |
| `nombre` | string | Nombre de la empresa. |
| `pais` | string | Código ISO de país (ej. PE). |
| `moneda` | string | Código ISO de la moneda principal (ej. PEN). |
| `maximo_usuarios` | number | Límite de usuarios permitidos para esta empresa. |
| `estado` | boolean | `true` si la empresa está activa. |
| `vencimiento` | string | Fecha de expiración de la suscripción (ISO 8601). |
| `creado_en` | string | Fecha de registro de la empresa. |

---

## Notas de Implementación (Frontend)

1. **Manejo de Autenticación:** 
   - El sistema soporta **Cookies** (recomendado para Web) y **Bearer Token** (recomendado para Mobile).
   - Si usas cookies, asegúrate de configurar `withCredentials: true` en tus peticiones de Axios o Fetch.
2. **Timezones:** Todas las fechas se entregan en formato UTC (Z). Se recomienda convertirlas a la zona horaria local del dispositivo para visualización.
3. **CORS:** Las peticiones deben provenir de orígenes permitidos configurados en el servidor.
