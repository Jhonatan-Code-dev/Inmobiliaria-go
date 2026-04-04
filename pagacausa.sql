
CREATE TABLE roles (
    id_rol TINYINT UNSIGNED PRIMARY KEY,
    rol VARCHAR(20) NOT NULL UNIQUE
)

CREATE TABLE roles_sistema (
    id_rol_sistema TINYINT UNSIGNED PRIMARY KEY,
    rol VARCHAR(20) NOT NULL UNIQUE
)

CREATE TABLE usuarios (
    id_usuario BIGINT PRIMARY KEY,
    usuario VARCHAR(30) NOT NULL UNIQUE,
    pass VARCHAR(255) NOT NULL,
    fecha_registro DATETIME NOT NULL,
    estado TINYINT(1) NOT NULL DEFAULT 1
)

CREATE TABLE empresas (
    id_empresa BIGINT PRIMARY KEY,
    nombre VARCHAR(30) NOT NULL UNIQUE,
    fecha_corte_servicio DATE NOT NULL,
    fecha registro DATETIME NOT NULL,
    estado TINYINT(1) NOT NULL DEFAULT 1
)

CREATE TABLE usuario_roles_sistema (
    id_usuario BIGINT NOT NULL UNIQUE,
    id_rol_sistema TINYINT UNSIGNED NOT NULL,
    PRIMARY KEY (id_usuario, id_rol_sistema),
    CONSTRAINT fk_urs_usuario FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    CONSTRAINT fk_urs_rol_sistema FOREIGN KEY (id_rol_sistema) REFERENCES roles_sistema(id_rol_sistema) ON DELETE CASCADE
)

CREATE TABLE empresa_usuarios (
    id_empresa_usuario BIGINT PRIMARY KEY,
    id_empresa BIGINT NOT NULL,
    id_usuario BIGINT NOT NULL,
    id_rol TINYINT UNSIGNED NOT NULL,
    estado TINYINT(1) NOT NULL,
    principal TINYINT(1) NOT NULL,
    CONSTRAINT fk_eu_empresa FOREIGN KEY (id_empresa) 
        REFERENCES empresas(id_empresa) ON DELETE CASCADE,
    CONSTRAINT fk_eu_usuario FOREIGN KEY (id_usuario) 
        REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    CONSTRAINT fk_eu_rol FOREIGN KEY (id_rol) 
        REFERENCES roles(id_rol)
)

CREATE TABLE control_ingreso_empresa (
    id_empresa BIGINT NOT NULL UNIQUE,
    hora_apertura TIME,
    hora_cierre TIME,
    modo_manual TINYINT(1),
    empresa_abierta TINYINT(1),
    CONSTRAINT fk_cie_empresa FOREIGN KEY (id_empresa) 
        REFERENCES empresas(id_empresa) ON DELETE CASCADE
)