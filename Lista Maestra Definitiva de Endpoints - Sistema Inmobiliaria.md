# Lista Maestra Definitiva de Endpoints - Sistema Inmobiliaria

Esta es la lista completa y actualizada de los endpoints para el backend inmobiliario conforme a Clean Architecture.

## 1. Módulo: Autenticación y Perfil
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| POST | `/auth/login` | Iniciar sesión | ✅ Existe |
| POST | `/auth/logout` | Cerrar sesión | ✅ Existe |
| GET | `/api/me` | Obtener mi perfil | ✅ Existe |
| PATCH | `/api/me/password` | Cambiar mi contraseña | ✅ Existe |

## 2. Módulo: Gestión de Staff (Personal Interno)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/staff` | Listar empleados | ✅ Existe |
| GET | `/api/user/staff/:id` | Ver detalle de empleado | ✅ Existe |
| POST | `/api/user/staff` | Crear empleado | ✅ Existe |
| PUT | `/api/user/staff/:id` | Editar empleado | ✅ Existe |
| DELETE | `/api/user/staff/:id` | Eliminar empleado | ✅ Existe |

## 3. Módulo: Inmuebles (Propiedades y Unidades)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/inmuebles` | Listar inmuebles | ✅ Existe |
| GET | `/api/user/inmuebles/:id` | Detalle + Unidades | ✅ Existe |
| POST | `/api/user/inmuebles` | Crear inmueble | ✅ Existe |
| PUT | `/api/user/inmuebles/:id` | Editar inmueble | ✅ Existe |
| DELETE | `/api/user/inmuebles/:id` | Eliminar inmueble | ✅ Existe |
| GET | `/api/user/inmuebles/:id/unidades` | Listar unidades | ✅ Existe |
| POST | `/api/user/inmuebles/:id/unidades` | Crear unidad | ✅ Existe |
| PUT | `/api/user/inmuebles/:propiedadId/unidades/:id` | Editar unidad | ✅ Existe |
| DELETE | `/api/user/inmuebles/:propiedadId/unidades/:id` | Eliminar unidad | ✅ Existe |

## 4. Módulo: Clientes (Inquilinos)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/clientes` | Listar clientes | ✅ Existe |
| GET | `/api/user/clientes/:id` | Detalle del cliente | ✅ Existe |
| POST | `/api/user/clientes` | Crear cliente | ✅ Existe |
| PUT | `/api/user/clientes/:id` | Editar cliente | ✅ Existe |
| DELETE | `/api/user/clientes/:id` | Eliminar cliente | ✅ Existe |

## 5. Módulo: Alquileres (Contratos)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/alquileres` | Listar contratos | ✅ Existe |
| GET | `/api/user/alquileres/:id` | Detalle del contrato | ✅ Existe |
| POST | `/api/user/alquileres` | Crear contrato | ✅ Existe |
| PUT | `/api/user/alquileres/:id` | Editar términos | ✅ Existe |
| DELETE | `/api/user/alquileres/:id` | Anular contrato | ✅ Existe |
| POST | `/api/user/alquileres/:id/terminar` | Finalizar contrato | ✅ Existe |

## 6. Módulo: Finanzas - Deudas (Cargos)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/cargos` | Listar cargos | ✅ Existe |
| GET | `/api/user/cargos/:id` | Detalle de cargo | ✅ Existe |
| POST | `/api/user/cargos` | Crear cargo manual | ✅ Existe |
| PUT | `/api/user/cargos/:id` | Editar cargo | ✅ Existe |
| DELETE | `/api/user/cargos/:id` | Eliminar cargo | ✅ Existe |

## 7. Módulo: Finanzas - Pagos (Cobros)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/pagos` | Historial de cobros | ✅ Existe |
| GET | `/api/user/pagos/:id` | Detalle de pago | ✅ Existe |
| POST | `/api/user/pagos` | Registrar cobro | ✅ Existe |
| PUT | `/api/user/pagos/:id` | Editar notas/método | ✅ Existe |
| DELETE | `/api/user/pagos/:id` | Anular pago | ✅ Existe |
| GET | `/api/user/pagos/pendientes` | Pagos pendientes mes | ✅ Existe |

## 8. Módulo: Servicios (Mediciones)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/servicios` | Listar consumos | ✅ Existe |
| GET | `/api/user/servicios/:id` | Detalle medición | ✅ Existe |
| POST | `/api/user/servicios` | Registrar lectura | ✅ Existe |
| PUT | `/api/user/servicios/:id` | Editar lectura | ✅ Existe |
| DELETE | `/api/user/servicios/:id` | Eliminar medición | ✅ Existe |

## 9. Módulo: Gastos (Egresos)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/gastos` | Listar egresos | ✅ Existe |
| POST | `/api/user/gastos` | Registrar gasto | ✅ Existe |
| PUT | `/api/user/gastos/:id` | Editar gasto | ✅ Existe |
| DELETE | `/api/user/gastos/:id` | Eliminar gasto | ✅ Existe |

## 10. Módulo: Mantenimiento (Tickets)
| Método | Endpoint | Descripción | Estado |
| :--- | :--- | :--- | :--- |
| GET | `/api/user/tickets` | Listar tickets | ✅ Existe |
| GET | `/api/user/tickets/:id` | Detalle ticket | ✅ Existe |
| POST | `/api/user/tickets` | Abrir incidencia | ✅ Existe |
| PUT | `/api/user/tickets/:id` | Actualizar estado | ✅ Existe |
| DELETE | `/api/user/tickets/:id` | Eliminar ticket | ✅ Existe |
