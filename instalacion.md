go get -u github.com/gofiber/fiber/v2

go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/gofiber/swagger
go get -u github.com/swaggo/files

swag init

go install entgo.io/ent/cmd/ent@latest
go get entgo.io/ent
go get entgo.io/ent/dialect@v0.14.5
go get entgo.io/ent/dialect/sql/schema@v0.14.5

go get github.com/go-sql-driver/mysql
go get -u github.com/gofiber/fiber/v2/middleware/cors

go get -u github.com/joho/godotenv
go get -u github.com/kelseyhightower/envconfig

go get golang.org/x/crypto/bcrypt

instalacion de jwt
go get github.com/golang-jwt/jwt/v5@v5.3.0

go generate ./ent

ent new Role RoleSistema Usuario Empresa UsuarioRolSistema EmpresaUsuario ControlIngresoEmpresa


se instalo esto :

go get -u github.com/h2non/bimg