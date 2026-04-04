// Package security proporciona utilidades para manejo de JWT, cookies y hashing
package security

import "golang.org/x/crypto/bcrypt"

// CostoBcrypt fijo para balancear seguridad y rendimiento
const costoBcrypt = 10

// ServicioHash define las operaciones para manejar contraseñas
type ServicioHash interface {
	Encriptar(password string) (string, error)
	Comparar(hash, passwordPlano string) bool
}

type servicioBcrypt struct{}

// NewServicioHash crea una instancia lista para usar con costo 10
func NewServicioHash() ServicioHash {
	return &servicioBcrypt{}
}

// Encriptar genera el hash de la contraseña de forma eficiente
func (s *servicioBcrypt) Encriptar(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), costoBcrypt)
	return string(bytes), err
}

// Comparar valida si la clave es correcta
func (s *servicioBcrypt) Comparar(hash, passwordPlano string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordPlano))
	return err == nil
}
