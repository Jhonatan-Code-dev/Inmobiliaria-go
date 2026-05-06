# 🚀 Guía de Consumo API - Sistema Inmobiliaria

Bienvenido a la documentación técnica del backend. Este sistema está diseñado bajo una arquitectura multi-tenant y un estricto control de acceso basado en roles (RBAC).

---

## 🔐 1. Autenticación y Seguridad Global

El sistema utiliza **JWT (JSON Web Tokens)** con una estrategia dual:
1. **Cookies HTTP-Only:** Ideal para SPAs (React, Vue, etc.) ya que ofrece protección contra ataques XSS.
2. **Bearer Token:** Header `Authorization: Bearer <token>` para consumo desde aplicaciones móviles u otras integraciones.

### Instrucción Crítica para Frontend:
Para que las cookies de sesión funcionen correctamente, todas las peticiones (Fetch/Axios) deben incluir la propiedad `withCredentials: true`.

---

## 🎭 2. Sistema de Roles (RBAC)

El sistema opera únicamente con estos **4 roles**. El `rol_id` define los permisos que el backend validará en cada petición:

| ID | Rol | Función Principal |
| :-- | :--- | :--- |
| **1** | **Administrador** | Gestión total de la empresa y del staff. |
| **2** | **Supervisor** | Monitoreo, revisión operativa y reportes. |
| **3** | **Vendedor** | Operativa comercial (Inmuebles, Clientes, Alquileres). |
| **4** | **Inventario** | Control físico de unidades y equipamiento. |

---

## 📂 3. Módulos de la API

A continuación, los enlaces a la documentación detallada de cada módulo:

| # | Módulo | Enlace a Documentación |
| :--- | :--- | :--- |
| **1.1** | **Usuarios y Autenticación** | [api_usuario_auth.md](api_usuario_auth.md) |
| **1.2** | **Gestión de Propiedades** | [api_inmuebles.md](api_inmuebles.md) |
| **1.3** | **Gestión de Clientes** | [api_clientes.md](api_clientes.md) |
| **1.4** | **Contratos de Alquiler** | [api_alquileres_pagos.md](api_alquileres_pagos.md) |
| **1.5** | **Cobranzas y Pagos** | [api_alquileres_pagos.md](api_alquileres_pagos.md) |
| **1.6** | **Gastos (Egresos)** | [api_gastos.md](api_gastos.md) |
| **1.7** | **Mantenimiento y Servicios** | [api_mantenimiento_servicios.md](api_mantenimiento_servicios.md) |
| **1.8** | **Detalle de Empresa** | [api_detalle_empresa.md](api_detalle_empresa.md) |

---

## 💡 Mejores Prácticas de Consumo

1. **Estructura de Respuesta:** Todas las respuestas exitosas devuelven un JSON. Los listados siempre incluyen un objeto `paginacion`.
2. **Manejo de Errores:** En caso de error, la API responde con un status code `4xx` o `5xx` y un body estándar:
   ```json
   { "message": "Descripción legible del error" }
   ```
3. **Multi-empresa:** El `empresa_id` es obligatorio en casi todas las peticiones. Si no se envía o no coincide con la del token del usuario, el sistema responderá `403 Forbidden`.
4. **Fechas:** Se manejan en formato ISO 8601 (UTC). El frontend debe transformarlas a la zona horaria del usuario para visualización.

---
*Documentación generada para facilitar la integración rápida y segura del Frontend.*
