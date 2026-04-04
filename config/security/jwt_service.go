// Package security proporciona utilidades para manejo de JWT, cookies y hashing
package security

import (
	"errors"
	"time"

	"rentals-go/internal/pkg/tiempo"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenInvalido = errors.New("token inválido")

// Claims define la estructura fija del contenido del token
type Claims struct {
	UsuarioID   int64  `json:"usuario_id"`
	EmpresaID   int64  `json:"empresa_id,omitempty"`
	Rol         string `json:"rol"`
	ZonaHoraria string `json:"zona_horaria,omitempty"`
	jwt.RegisteredClaims
}

// ServicioJWT define las operaciones de tokens con nombres en español
type ServicioJWT interface {
	GenerarToken(usuarioID int64, empresaID int64, rol string, zona string) (string, error)
	ValidarToken(token string) (*Claims, error)
}

type servicioJWT struct {
	secreto    string
	expiracion time.Duration
}

// NewServicioJWT inicializa el servicio con la configuración del entorno
func NewServicioJWT(secreto string, duracion time.Duration) ServicioJWT {
	return &servicioJWT{secreto: secreto, expiracion: duracion}
}

// GenerarToken crea un nuevo JWT usando una estructura fija de Claims
func (s *servicioJWT) GenerarToken(usuarioID int64, empresaID int64, rol string, zona string) (string, error) {
	ahora := tiempo.AhoraUTC()
	claims := Claims{
		UsuarioID:   usuarioID,
		EmpresaID:   empresaID,
		Rol:         rol,
		ZonaHoraria: zona,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(ahora.Add(s.expiracion)),
			IssuedAt:  jwt.NewNumericDate(ahora),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secreto))
}

// ValidarToken comprueba la firma y retorna los datos estructurados
func (s *servicioJWT) ValidarToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalido
		}
		return []byte(s.secreto), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrTokenInvalido
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTokenInvalido
	}

	return claims, nil
}
