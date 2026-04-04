Uso bajo acoplamiento y alta cohesión, aislando los DTOs por controller para evitar impactos colaterales ante cambios.

PROMT DOMAIN:

Define qué debe y qué no debe ir en la capa Domain siguiendo Clean Architecture. El Domain debe contener solo entidades del negocio como structs puros, value objects, reglas de negocio, errores del dominio e interfaces (repositorios/ports), sin dependencias externas. No debe incluir JSON, DTOs, HTTP, frameworks, ORMs, SQL, lógica de infraestructura, configuración ni detalles de persistencia o transporte. El Domain debe ser independiente, testeable y expresar únicamente el lenguaje del negocio.

PROMT CONTROLLER

Define qué debe ir en la capa Controller siguiendo Clean Architecture. El Controller debe recibir la petición (HTTP), validar datos básicos, definir DTOs de request y response en el mismo archivo, mapear JSON ↔ structs, invocar los Use Cases, manejar errores y construir las respuestas. No debe contener reglas de negocio ni acceso a base de datos u ORMs. Se debe usar bajo acoplamiento y alta cohesión, aislando los DTOs por controller para evitar impactos colaterales ante cambios.

PROMT SERVICE
Define qué debe ir en la capa Service siguiendo Clean Architecture. El Service debe orquestar y coordinar reglas de negocio complejas entre entidades del dominio, ejecutar casos de uso que involucren múltiples repositorios o agregados y manejar transacciones o flujos compuestos. No debe manejar HTTP, JSON, DTOs, ni detalles de infraestructura. El Service actúa como una capa de aplicación que encapsula la lógica de negocio reutilizable, manteniendo bajo acoplamiento y alta cohesión.

PROMT REPOSITORY
Define qué debe ir en la capa Repository siguiendo Clean Architecture. El Repository debe implementar las interfaces definidas en el Domain, encargarse exclusivamente del acceso a datos y la persistencia (ORM, SQL, Ent, etc.), mapear modelos de infraestructura ↔ entidades del dominio y ocultar los detalles de la base de datos. No debe contener reglas de negocio ni lógica de aplicación. Su objetivo es aislar la persistencia, manteniendo bajo acoplamiento y alta cohesión, para que el dominio y los casos de uso no dependan de la tecnología de almacenamiento.

# PROMT COMPLETO UNIDO:

Analiza y refactoriza el código existente para que cumpla estrictamente con Clean Architecture, separando correctamente las responsabilidades por capas y documentando adecuadamente las APIs con Swagger/OpenAPI.

En la capa Domain, conserva únicamente entidades del negocio como structs puros, value objects, reglas de negocio, errores del dominio e interfaces (repositorios/ports), sin dependencias externas. Elimina cualquier uso de JSON, DTOs, HTTP, frameworks, ORMs, SQL, configuración o detalles de persistencia o transporte. El Domain debe ser totalmente independiente, testeable y expresar solo el lenguaje del negocio.

En la capa Controller, limita la responsabilidad a recibir peticiones HTTP, realizar validaciones básicas, definir DTOs de request y response en el mismo archivo, mapear JSON ↔ structs, invocar los Use Cases, manejar errores y construir las respuestas. Documenta cada endpoint con Swagger/OpenAPI (resumen, descripción, tags, parámetros, body, respuestas, códigos HTTP y modelos), asegurando una documentación clara, consistente y mantenible. No incluyas reglas de negocio ni acceso a base de datos u ORMs. Aplica bajo acoplamiento y alta cohesión, aislando los DTOs por controller para evitar impactos colaterales ante cambios.

En la capa Service, concentra la orquestación de reglas de negocio complejas, la coordinación entre múltiples entidades o agregados y la ejecución de flujos compuestos o transacciones. No manejes HTTP, JSON, DTOs ni detalles de infraestructura. El Service debe encapsular lógica de negocio reutilizable y mantenerse independiente del transporte y la persistencia.

En la capa Repository, implementa exclusivamente las interfaces definidas en el Domain, encargándote del acceso a datos y la persistencia (ORM, SQL, Ent, etc.). Mapea modelos de infraestructura ↔ entidades del dominio y oculta completamente los detalles de la base de datos. No incluyas reglas de negocio ni lógica de aplicación, garantizando bajo acoplamiento y alta cohesión para que el dominio y los casos de uso no dependan de la tecnología de almacenamiento.
