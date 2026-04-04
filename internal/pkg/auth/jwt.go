package auth

import (
	"time"

	"rentals-go/internal/pkg/tiempo"

	"github.com/golang-jwt/jwt/v5"
)

// Claims comunes para admin y usuario
type Claims struct {
	UsuarioID int    `json:"usuario_id,omitempty"`
	AdminID   int    `json:"admin_id,omitempty"`
	EmpresaID int    `json:"empresa_id,omitempty"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

// GenerarToken crea un JWT con expiración dada.
func GenerarToken(secret []byte, exp time.Duration, claims Claims) (string, error) {
	now := tiempo.AhoraUTC()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidarToken parsea y valida un JWT.
func ValidarToken(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
