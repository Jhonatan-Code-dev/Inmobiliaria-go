CREATE DATABASE IF NOT EXISTS rentals_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE rentals_go;

CREATE TABLE empresas (
  id INT NOT NULL AUTO_INCREMENT,
  nombre VARCHAR(150) NOT NULL,
  documento_fiscal VARCHAR(50) NULL,
  correo VARCHAR(150) NULL,
  telefono VARCHAR(30) NULL,
  direccion VARCHAR(255) NULL,
  ciudad VARCHAR(120) NULL,
  pais VARCHAR(2) NULL,
  moneda VARCHAR(10) NOT NULL DEFAULT 'PEN',
  maximo_usuarios INT NOT NULL DEFAULT 1,
  estado VARCHAR(20) NOT NULL DEFAULT 'activa',
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_empresas_documento_fiscal (documento_fiscal)
);

CREATE TABLE roles (
  id INT NOT NULL AUTO_INCREMENT,
  nombre VARCHAR(60) NOT NULL,
  descripcion VARCHAR(255) NULL,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_roles_nombre (nombre)
);

CREATE TABLE usuarios (
  id INT NOT NULL AUTO_INCREMENT,
  nombres VARCHAR(120) NOT NULL,
  apellidos VARCHAR(120) NULL,
  correo VARCHAR(150) NOT NULL,
  telefono VARCHAR(30) NULL,
  hash_contrasena VARCHAR(255) NOT NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'activo',
  ultimo_acceso TIMESTAMP NULL,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_usuarios_correo (correo)
);

CREATE TABLE tipos_identificacion (
  id INT NOT NULL AUTO_INCREMENT,
  codigo VARCHAR(20) NOT NULL,
  nombre VARCHAR(80) NOT NULL,
  pais VARCHAR(2) NULL,
  activo TINYINT(1) NOT NULL DEFAULT 1,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_tipos_identificacion_codigo (codigo)
);

CREATE TABLE empresa_usuarios (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  usuario_id INT NOT NULL,
  rol_id INT NOT NULL,
  principal TINYINT(1) NOT NULL DEFAULT 0,
  estado VARCHAR(20) NOT NULL DEFAULT 'activo',
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_empresa_usuarios_empresa_usuario (empresa_id, usuario_id),
  CONSTRAINT fk_empresa_usuarios_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id),
  CONSTRAINT fk_empresa_usuarios_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
  CONSTRAINT fk_empresa_usuarios_rol FOREIGN KEY (rol_id) REFERENCES roles(id)
);

CREATE TABLE clientes (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  tipo_identificacion_id INT NOT NULL,
  documento_numero VARCHAR(30) NOT NULL,
  nombres VARCHAR(120) NOT NULL,
  apellidos VARCHAR(120) NULL,
  correo VARCHAR(150) NULL,
  fecha_nacimiento DATE NULL,
  nacionalidad VARCHAR(60) NULL,
  direccion VARCHAR(255) NULL,
  contacto_emergencia VARCHAR(150) NULL,
  telefono_emergencia VARCHAR(30) NULL,
  notas VARCHAR(1000) NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'activo',
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_clientes_documento (empresa_id, tipo_identificacion_id, documento_numero),
  CONSTRAINT fk_clientes_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id),
  CONSTRAINT fk_clientes_tipo_identificacion FOREIGN KEY (tipo_identificacion_id) REFERENCES tipos_identificacion(id)
);

CREATE TABLE cliente_telefonos (
  id INT NOT NULL AUTO_INCREMENT,
  cliente_id INT NOT NULL,
  telefono VARCHAR(30) NOT NULL,
  etiqueta VARCHAR(30) NULL,
  principal TINYINT(1) NOT NULL DEFAULT 0,
  whatsapp TINYINT(1) NOT NULL DEFAULT 0,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  KEY idx_cliente_telefonos_cliente (cliente_id),
  CONSTRAINT fk_cliente_telefonos_cliente FOREIGN KEY (cliente_id) REFERENCES clientes(id) ON DELETE CASCADE
);

CREATE TABLE propiedades (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  nombre VARCHAR(150) NOT NULL,
  tipo VARCHAR(20) NOT NULL DEFAULT 'casa',
  descripcion VARCHAR(1000) NULL,
  direccion VARCHAR(255) NOT NULL,
  ciudad VARCHAR(120) NULL,
  region VARCHAR(120) NULL,
  pais VARCHAR(2) NULL,
  codigo_postal VARCHAR(20) NULL,
  total_pisos INT NOT NULL DEFAULT 1,
  total_unidades INT NOT NULL DEFAULT 1,
  estado VARCHAR(20) NOT NULL DEFAULT 'activa',
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  KEY idx_propiedades_empresa (empresa_id),
  CONSTRAINT fk_propiedades_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id)
);

CREATE TABLE unidades (
  id INT NOT NULL AUTO_INCREMENT,
  propiedad_id INT NOT NULL,
  codigo VARCHAR(30) NOT NULL,
  nombre VARCHAR(120) NULL,
  tipo VARCHAR(20) NOT NULL DEFAULT 'cuarto',
  numero_piso INT NULL,
  dormitorios INT NOT NULL DEFAULT 0,
  banos INT NOT NULL DEFAULT 0,
  area_m2 DECIMAL(10,2) NULL,
  capacidad INT NOT NULL DEFAULT 1,
  precio_base DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  deposito_requerido DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  incluye_agua TINYINT(1) NOT NULL DEFAULT 0,
  incluye_luz TINYINT(1) NOT NULL DEFAULT 0,
  incluye_internet TINYINT(1) NOT NULL DEFAULT 0,
  notas VARCHAR(1000) NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'disponible',
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_unidades_propiedad_codigo (propiedad_id, codigo),
  CONSTRAINT fk_unidades_propiedad FOREIGN KEY (propiedad_id) REFERENCES propiedades(id)
);

CREATE TABLE contratos (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  cliente_id INT NOT NULL,
  unidad_id INT NOT NULL,
  codigo VARCHAR(40) NOT NULL,
  tipo VARCHAR(20) NOT NULL DEFAULT 'alquiler',
  fecha_inicio DATE NOT NULL,
  fecha_fin DATE NULL,
  dia_vencimiento INT NOT NULL,
  moneda VARCHAR(10) NOT NULL DEFAULT 'PEN',
  monto_renta DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  monto_deposito DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  mora_diaria DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  servicios_incluidos TINYINT(1) NOT NULL DEFAULT 0,
  activo_para_cobro TINYINT(1) NOT NULL DEFAULT 1,
  estado VARCHAR(20) NOT NULL DEFAULT 'activo',
  observaciones VARCHAR(1500) NULL,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_contratos_empresa_codigo (empresa_id, codigo),
  KEY idx_contratos_cliente (cliente_id),
  KEY idx_contratos_unidad (unidad_id),
  CONSTRAINT fk_contratos_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id),
  CONSTRAINT fk_contratos_cliente FOREIGN KEY (cliente_id) REFERENCES clientes(id),
  CONSTRAINT fk_contratos_unidad FOREIGN KEY (unidad_id) REFERENCES unidades(id)
);

CREATE TABLE cargos (
  id INT NOT NULL AUTO_INCREMENT,
  contrato_id INT NOT NULL,
  concepto VARCHAR(20) NOT NULL DEFAULT 'renta',
  descripcion VARCHAR(255) NULL,
  periodo_inicio DATE NOT NULL,
  periodo_fin DATE NOT NULL,
  fecha_emision DATE NOT NULL,
  fecha_vencimiento DATE NOT NULL,
  monto DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  saldo DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
  generado_automaticamente TINYINT(1) NOT NULL DEFAULT 0,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  KEY idx_cargos_contrato_vencimiento (contrato_id, fecha_vencimiento),
  CONSTRAINT fk_cargos_contrato FOREIGN KEY (contrato_id) REFERENCES contratos(id)
);

CREATE TABLE pagos (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  cliente_id INT NULL,
  contrato_id INT NULL,
  numero_recibo VARCHAR(40) NOT NULL,
  fecha_pago DATETIME NOT NULL,
  moneda VARCHAR(10) NOT NULL DEFAULT 'PEN',
  monto_total DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  metodo VARCHAR(20) NOT NULL DEFAULT 'efectivo',
  referencia VARCHAR(120) NULL,
  notas VARCHAR(1000) NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'confirmado',
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_pagos_empresa_recibo (empresa_id, numero_recibo),
  KEY idx_pagos_cliente (cliente_id),
  KEY idx_pagos_contrato (contrato_id),
  CONSTRAINT fk_pagos_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id),
  CONSTRAINT fk_pagos_cliente FOREIGN KEY (cliente_id) REFERENCES clientes(id),
  CONSTRAINT fk_pagos_contrato FOREIGN KEY (contrato_id) REFERENCES contratos(id)
);

CREATE TABLE pago_aplicaciones (
  id INT NOT NULL AUTO_INCREMENT,
  pago_id INT NOT NULL,
  cargo_id INT NOT NULL,
  monto_aplicado DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_pago_aplicaciones_pago_cargo (pago_id, cargo_id),
  CONSTRAINT fk_pago_aplicaciones_pago FOREIGN KEY (pago_id) REFERENCES pagos(id) ON DELETE CASCADE,
  CONSTRAINT fk_pago_aplicaciones_cargo FOREIGN KEY (cargo_id) REFERENCES cargos(id) ON DELETE CASCADE
);

CREATE TABLE servicio_mediciones (
  id INT NOT NULL AUTO_INCREMENT,
  unidad_id INT NOT NULL,
  tipo_servicio VARCHAR(20) NOT NULL DEFAULT 'agua',
  periodo_inicio DATE NOT NULL,
  periodo_fin DATE NOT NULL,
  lectura_anterior DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  lectura_actual DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  consumo DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  tarifa_unitaria DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  monto_total DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  observaciones VARCHAR(1000) NULL,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_servicio_mediciones_periodo (unidad_id, tipo_servicio, periodo_inicio, periodo_fin),
  CONSTRAINT fk_servicio_mediciones_unidad FOREIGN KEY (unidad_id) REFERENCES unidades(id)
);

CREATE TABLE gastos (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  propiedad_id INT NULL,
  unidad_id INT NULL,
  categoria VARCHAR(20) NOT NULL DEFAULT 'otro',
  descripcion VARCHAR(255) NOT NULL,
  fecha_gasto DATETIME NOT NULL,
  monto DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  metodo_pago VARCHAR(20) NOT NULL DEFAULT 'efectivo',
  referencia VARCHAR(120) NULL,
  pagado_a VARCHAR(150) NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'pagado',
  notas VARCHAR(1000) NULL,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  KEY idx_gastos_empresa (empresa_id),
  KEY idx_gastos_propiedad (propiedad_id),
  KEY idx_gastos_unidad (unidad_id),
  CONSTRAINT fk_gastos_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id),
  CONSTRAINT fk_gastos_propiedad FOREIGN KEY (propiedad_id) REFERENCES propiedades(id),
  CONSTRAINT fk_gastos_unidad FOREIGN KEY (unidad_id) REFERENCES unidades(id)
);

CREATE TABLE movimientos_caja (
  id INT NOT NULL AUTO_INCREMENT,
  empresa_id INT NOT NULL,
  pago_id INT NULL,
  gasto_id INT NULL,
  tipo VARCHAR(20) NOT NULL DEFAULT 'ingreso',
  concepto VARCHAR(150) NOT NULL,
  fecha_movimiento DATETIME NOT NULL,
  monto DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  metodo VARCHAR(20) NOT NULL DEFAULT 'efectivo',
  referencia VARCHAR(120) NULL,
  observaciones VARCHAR(1000) NULL,
  creado_en TIMESTAMP NOT NULL,
  actualizado_en TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  KEY idx_movimientos_caja_empresa (empresa_id),
  KEY idx_movimientos_caja_pago (pago_id),
  KEY idx_movimientos_caja_gasto (gasto_id),
  CONSTRAINT fk_movimientos_caja_empresa FOREIGN KEY (empresa_id) REFERENCES empresas(id),
  CONSTRAINT fk_movimientos_caja_pago FOREIGN KEY (pago_id) REFERENCES pagos(id),
  CONSTRAINT fk_movimientos_caja_gasto FOREIGN KEY (gasto_id) REFERENCES gastos(id)
);

INSERT INTO roles (nombre, descripcion) VALUES
  ('admin', 'Control total de la empresa'),
  ('cobrador', 'Registra pagos y seguimiento de contratos'),
  ('operador', 'Gestion operativa de propiedades y clientes');

INSERT INTO tipos_identificacion (codigo, nombre, pais) VALUES
  ('DNI', 'Documento Nacional de Identidad', 'PE'),
  ('CE', 'Carnet de Extranjeria', 'PE'),
  ('PAS', 'Pasaporte', NULL),
  ('RUC', 'Registro Unico de Contribuyentes', 'PE');
