# Documentación de API: Usuarios, Roles y Accesos (RBAC)

Esta documentación detalla los endpoints necesarios para la gestión de usuarios, autenticación modular y control de acceso basado en roles (RBAC) para el sistema de gestión inmobiliaria.

## 🔑 Autenticación y Seguridad

El sistema utiliza un esquema híbrido de autenticación para máxima flexibilidad:
- **Cookies HttpOnly:** Recomendado para aplicaciones Web (manejo automático por el navegador).
- **JWT (Bearer Token):** Recomendado para aplicaciones Mobile o integraciones externas.

---

## 👥 Roles del Sistema (RBAC)

El acceso está restringido a **4 roles predefinidos**, cada uno con responsabilidades específicas:

| ID | Rol | Descripción |
| :-- | :--- | :--- |
| **1** | **Administrador** | Acceso total al sistema, configuración de la empresa y gestión de staff. |
| **2** | **Supervisor** | Supervisión operativa, acceso a reportes críticos y métricas de desempeño. |
| **3** | **Vendedor** | Gestión de clientes, registro de inmuebles y creación de contratos. |
| **4** | **Inventario** | Control de activos, equipamiento de unidades y estados físicos. |

---

## 🛠️ Resumen de Endpoints

### Autenticación
| Método | Endpoint | Descripción | Acceso |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/login` | Iniciar sesión y obtener token/cookie. | Público |
| `POST` | `/auth/logout` | Finalizar sesión y limpiar cookies. | Autenticado |
| `GET` | `/me` | Obtener perfil del usuario actual. | Autenticado |
| `PATCH` | `/me/password` | Cambiar contraseña del usuario actual. | Autenticado |

### Gestión de Staff (Personal)
| Método | Endpoint | Descripción | Acceso |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/user/staff/roles` | Obtener catálogo de roles disponibles. | Público / Autenticado |
| `GET` | `/api/user/staff` | Listar todo el personal de la empresa. | Admin / Super |
| `GET` | `/api/user/staff/:id` | Ver detalle de un empleado. | Admin / Super |
| `POST` | `/api/user/staff` | Registrar un nuevo usuario de staff. | Admin |
| `PUT` | `/api/user/staff/:id` | Editar rol o estado de un empleado. | Admin |
| `DELETE` | `/api/user/staff/:id`| Eliminar o dar de baja a un empleado. | Admin |

---

## 1. Autenticación

### 1.1 Login
`POST /auth/login`

**Request Body:**
```json
{
  "usuario": "admin_demo",
  "contrasena": "Password123!"
}
```

**Respuesta Exitosa (200 OK):**
Establece la cookie `token_usuario` (HttpOnly, Secure).
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "usuario": "admin_demo",
    "empresa_id": 10
  },
  "empresa": {
    "id": 10,
    "nombre": "Inmobiliaria El Sol",
    "pais": "PE",
    "moneda": "PEN",
    "estado": true
  }
}
```

---

### 1.2 Obtener Perfil (Me)
`GET /me`

Retorna los datos del usuario en sesión. Sirve para validar si el token sigue vigente.

**Respuesta Exitosa (200 OK):**
```json
{
  "user": {
    "id": 1,
    "usuario": "admin_demo",
    "rol_id": 1,
    "rol_nombre": "administrador",
    "empresa_id": 10
  },
  "empresa": {
    "id": 10,
    "nombre": "Inmobiliaria El Sol"
  }
}
```

---

## 2. Gestión de Staff (Usuarios de Empresa)

### 2.1 Listar Roles Disponibles
`GET /api/user/staff/roles`

Obtiene el catálogo de roles definidos en el sistema para ser usados en el registro o edición de staff.

**Respuesta Exitosa (200 OK):**
```json
[
  {
    "id": 1,
    "nombre": "administrador",
    "descripcion": "Rol con acceso administrativo total a la empresa"
  },
  {
    "id": 2,
    "nombre": "supervisor",
    "descripcion": "Rol para supervisión operativa y reportes"
  },
  {
    "id": 3,
    "nombre": "vendedor",
    "descripcion": "Rol para gestión comercial, clientes y contratos"
  },
  {
    "id": 4,
    "nombre": "inventario",
    "descripcion": "Rol para control de activos y estado de unidades"
  }
]
```

---

### 2.2 Listar Staff
`GET /api/user/staff`

**Parámetros de Query:**
- `pag`: (int) Número de página (defecto: 1).
- `por_pagina`: (int) Cantidad de registros (defecto: 10).
- `buscar`: (string) Filtrar por nombre de usuario.

**Respuesta Exitosa (200 OK):**
```json
{
  "datos": [
    {
      "id": 1,
      "usuario_id": 5,
      "usuario": "juan.perez",
      "rol_id": 3,
      "rol_nombre": "vendedor",
      "principal": false,
      "estado": "activo"
    }
  ],
  "paginacion": {
    "total": 1,
    "pagina_actual": 1,
    "por_pagina": 10,
    "paginas": 1
  }
}
```

---

### 2.2 Crear Nuevo Usuario
`POST /api/user/staff`

**Request Body:**
```json
{
  "usuario": "nuevo.vendedor",
  "contrasena": "Temporal2024",
  "rol_id": 3
}
```
*Nota: El `empresa_id` se toma automáticamente del token de quien crea.*

---

### 2.3 Actualizar Usuario
`PUT /api/user/staff/:id`

Permite cambiar el rol o el estado (activo/inactivo) de un usuario.

**Request Body:**
```json
{
  "rol_id": 2,
  "estado": "inactivo"
}
```

---

### 2.4 Eliminar Usuario
`DELETE /api/user/staff/:id`

**Restricción:** No se puede eliminar al usuario marcado como `principal: true` (dueño de la cuenta).

---

---

## 🚀 Ejemplos de Consumo (JavaScript/Frontend)

### A. Obtener Roles para un Selector (Select)
Útil para cargar las opciones en el formulario de creación de staff.

```javascript
async function cargarRoles() {
  const response = await fetch('/api/user/staff/roles', {
    method: 'GET',
    headers: {
      'Accept': 'application/json'
    }
  });

  if (!response.ok) {
    const error = await response.json();
    console.error('Error al cargar roles:', error.message);
    return [];
  }

  const roles = await response.json();
  // El resultado es un array: [{id: 1, nombre: 'administrador'}, ...]
  return roles;
}
```

### B. Listar Staff de la Empresa
Consumo recomendado para la tabla de administración de personal.

```javascript
async function listarPersonal(pagina = 1, busqueda = '') {
  const url = `/api/user/staff?pag=${pagina}&buscar=${busqueda}`;
  
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Accept': 'application/json'
    }
  });

  const staffData = await response.json();
  /* 
    Respuesta: 
    { 
       datos: [...], 
       paginacion: { total: 10, paginas: 1, pagina_actual: 1 } 
    } 
  */
  return staffData;
}
```

### C. Crear un Nuevo Empleado
Ejemplo de envío de datos del formulario.

```javascript
async function registrarEmpleado(datos) {
  // datos = { usuario: 'pedro.perez', contrasena: 'Pass123', rol_id: 3 }
  const response = await fetch('/api/user/staff', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    body: JSON.stringify(datos)
  });

  if (response.status === 201) {
    alert('Empleado registrado con éxito');
  } else {
    const error = await response.json();
    alert('Error: ' + error.message);
  }
}
```

---
*Documentación generada para el equipo de desarrollo de Rentals Go.*
