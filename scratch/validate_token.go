package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3VhcmlvX2lkIjoxLCJlbXByZXNhX2lkIjoxLCJyb2wiOiJ1c3VhcmlvIiwiZXhwIjoxNzc4MzgxMTQ1LCJpYXQiOjE3NzgyOTQ3NDV9.cz6BBXTbiYZb_pH_FR9nh9M7-5l4opR0twJqjQ0A13Q"
	secret := []byte("PAGA_CAUSA_SUPER_SECRET_KEY_2026")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if token.Valid {
		fmt.Println("Token is valid")
	} else {
		fmt.Println("Token is invalid")
	}
}
