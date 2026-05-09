# Guía de Supervisión: Panel de Administración de Asistencia (Actualizado)

Esta documentación detalla los endpoints para supervisar a todos los trabajadores, filtrar por estado y gestionar registros.

---

## 1. Listado Maestro de Asistencia (Supervisión General)

Este endpoint permite ver la actividad de **toda la empresa** en una sola vista, incluyendo el nombre del trabajador y su estado.

*   **URL:** `/api/user/asistencia/registros`
*   **Método:** `GET`
*   **Query Params (Filtros):**
    *   `empresa_id` (Requerido): ID de la empresa.
    *   `estado` (Opcional): Filtra por `tarde`, `puntual`, `falta` o `justificado`. **Ideal para ver quién llegó tarde hoy.**
    *   `desde` / `hasta` (Opcional): Formato `YYYY-MM-DD`.
    *   `usuario_id` (Opcional): Si quieres ver solo a uno.
    *   `pag` / `limite`: Para paginación.

**Respuesta Exitosa (JSON):**
```json
[
  {
    "id": 15,
    "usuario_id": 4,
    "usuario_nombre": "juan_perez",
    "fecha": "2026-05-08T00:00:00Z",
    "hora_entrada": "2026-05-08T08:45:00Z",
    "hora_salida": null,
    "estado": "tarde",
    "horas_trabajadas": null
  }
]
```

---

## 2. Gestión de Registros y Errores

### 2.1 Eliminar Registro Incorrecto
Si un trabajador cometió un error al marcar, puedes borrar su registro del día para que vuelva a marcar correctamente.

*   **URL:** `/api/user/asistencia/registros/:id?empresa_id=1`
*   **Método:** `DELETE`

---

## 3. Gestión de Horarios y Turnos

### 3.1 Consultar Horario Configurado
Permite ver qué horario tiene asignado un usuario antes de hacer cambios.

*   **URL:** `/api/user/asistencia/horarios/detalle?empresa_id=1&usuario_id=4`
*   **Método:** `GET`

### 3.2 Asignar / Cambiar Horario
*   **URL:** `/api/user/asistencia/horarios?empresa_id=1`
*   **Método:** `POST`
*   **Body:**
```json
{
  "usuario_id": 4,
  "hora_entrada": "08:00",
  "hora_salida": "17:00",
  "tolerancia_minutos": 15,
  "dias_laborables": "1,2,3,4,5"
}
```

---

## 4. Control de Permisos

### 4.1 Listado de Solicitudes Pendientes
*   **URL:** `/api/user/asistencia/permisos?empresa_id=1&estado=pendiente`
*   **Método:** `GET`

### 4.2 Decisión Administrativa
*   **URL:** `/api/user/asistencia/permisos/:id/estado?empresa_id=1`
*   **Método:** `PUT`
*   **Body:** `{"estado": "aprobado", "respuesta": "OK"}`
