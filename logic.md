REQUEST
↓
VALIDACIONES
↓
UPDATE DB (TRANSACCIÓN)
↓
SI OK → UPDATE VALKEY
↓
RESPONSE

-
-
-
-
-

HTTP(JSON)
↓
Controller (DTO JSON)
↓
UseCase
↓
Domain (structs puros)
↓
Repository
