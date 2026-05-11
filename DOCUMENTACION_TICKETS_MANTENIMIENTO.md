# Documentación: Módulo de Tickets de Mantenimiento

Este documento detalla los endpoints RESTful para la gestión del ciclo de vida completo de un Ticket de Mantenimiento en el sistema de Rentals.

## Flujo Básico
1. Un cliente (o administrador) **Crea** un ticket indicando la unidad afectada, asunto, descripción y nivel de prioridad.
2. Los supervisores pueden **Listar** y buscar los tickets pendientes.
3. El administrador toma acción sobre el ticket y puede **Cambiar su Estado** (por ejemplo, a `en_progreso`).
4. Una vez solucionado, el estado se cambia a `resuelto` o `cerrado`.

---

## 1. Crear Ticket
Crea un nuevo ticket asociado a una unidad y, opcionalmente, a un cliente específico. El estado inicial siempre será `abierto`.

- **Ruta:** `POST /api/user/tickets`
- **Autenticación:** Requerida (Bearer Token)
- **Body (JSON):**
```json
{
  "unidad_id": 1,
  "cliente_id": 15,
  "asunto": "Fuga de agua en el baño principal",
  "descripcion": "El lavabo gotea constantemente incluso estando cerrado. Ya se reportó antes.",
  "prioridad": "alta"
}
```
*Notas:* 
- `prioridad` puede ser: `baja`, `media`, `alta`.
- `cliente_id` es opcional (solo se envía si se sabe quién reporta).

---

## 2. Listar Tickets Paginados
Obtiene una lista de tickets. Permite realizar búsquedas de texto completo en el asunto/descripción y filtrar por atributos.

- **Ruta:** `GET /api/user/tickets`
- **Query Params:**
  - `pag`: (Opcional) Número de página (default: 1).
  - `por_pagina`: (Opcional) Registros por página (default: 10).
  - `propiedad_id`: (Opcional) ID del Inmueble. Muestra los tickets de todas las unidades de este inmueble.
  - `unidad_id`: (Opcional) ID de la unidad a filtrar.
  - `estado`: (Opcional) `abierto`, `en_progreso`, `resuelto`, `cerrado`.
  - `buscar`: (Opcional) Busca coincidencia de texto en el `asunto` o `descripcion`.

**Ejemplo de Respuesta (JSON):**
```json
{
  "datos": [
    {
      "id": 1,
      "empresa_id": 1,
      "unidad_id": 1,
      "unidad_nombre": "Departamento 101",
      "cliente_id": 15,
      "cliente_nombre": "Juan Pérez",
      "asunto": "Fuga de agua en el baño principal",
      "descripcion": "El lavabo gotea...",
      "prioridad": "alta",
      "estado": "abierto",
      "fecha_apertura": "2026-05-09T18:00:00Z"
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

## 3. Obtener Resumen de Tickets (Dashboard)
Retorna el total de tickets agrupados por su estado, ideal para pintar gráficas o contadores en el Dashboard.

- **Ruta:** `GET /api/user/tickets/resumen`
- **Query Params:**
  - `propiedad_id`: (Opcional) Filtra el resumen para un inmueble específico.

**Ejemplo de Respuesta (JSON):**
```json
{
  "total": 15,
  "abiertos": 5,
  "en_progreso": 3,
  "resueltos": 5,
  "cerrados": 2
}
```

---

## 4. Obtener Detalle de Ticket
Obtiene la información completa de un ticket específico.

- **Ruta:** `GET /api/user/tickets/:id`

---

## 5. Cambiar Estado del Ticket
Este es el endpoint específico recomendado para transicionar un ticket a lo largo de su ciclo de vida sin afectar otros datos (como el título o descripción).

- **Ruta:** `PATCH /api/user/tickets/:id/estado`
- **Body (JSON):**
```json
{
  "estado": "en_progreso"
}
```
*Estados Válidos:* `abierto`, `en_progreso`, `resuelto`, `cerrado`.

---

## 6. Actualizar Ticket Completo
Permite sobreescribir la información base del ticket. Recomendado sólo si el usuario necesita corregir un error tipográfico en la creación.

- **Ruta:** `PUT /api/user/tickets/:id`
- **Body (JSON):** *(Igual a la Creación, agregando el campo opcional de `estado`)*

---

## 7. Eliminar Ticket
Elimina un ticket definitivamente de la base de datos (Borrado físico).

- **Ruta:** `DELETE /api/user/tickets/:id`

---

## 8. Catálogos para el Formulario (Dropdowns)
Este endpoint es vital para el frontend. Retorna las listas de Inmuebles, Clientes, Prioridades y Estados para llenar los selectores (dropdowns) del formulario de "Abrir Nuevo Ticket".

- **Ruta:** `GET /api/user/tickets/config-formulario`
- **Autenticación:** Requerida
- **Respuesta (JSON):**
```json
{
  "inmuebles": [
    { "id": 1, "nombre": "Edificio El Sol" },
    { "id": 2, "nombre": "Residencial Luna" }
  ],
  "clientes": [
    { "id": 15, "nombre": "Juan Pérez" },
    { "id": 22, "nombre": "María García" }
  ],
  "prioridades": ["baja", "media", "alta"],
  "estados": ["abierto", "en_progreso", "resuelto", "cerrado"]
}
```

### Flujo recomendado para el Frontend:
1.  **Cargar Catálogos**: Al abrir el modal de "Abrir Nuevo Ticket", llamar a `/api/user/tickets/config-formulario`.
2.  **Seleccionar Inmueble**: El usuario elige "Edificio El Sol" (ID: 1).
3.  **Filtrar Unidades**: Una vez elegido el inmueble, llamar a `GET /api/user/inmuebles/1/unidades` para obtener la lista de unidades de ese edificio.
4.  **Seleccionar Unidad**: El usuario elige la unidad (ej. "Dpto 101", ID: 5).
5.  **Enviar**: Ahora el frontend ya tiene el `unidad_id` y `cliente_id` correctos para enviar al `POST /api/user/tickets`.

---

## 9. Listar Unidades por Inmueble
(Reutilizado del módulo de Inmuebles para conveniencia del ticket)
- **Ruta:** `GET /api/user/inmuebles/:id/unidades`
- **Uso:** Sirve para llenar el dropdown de unidades después de elegir un inmueble.
