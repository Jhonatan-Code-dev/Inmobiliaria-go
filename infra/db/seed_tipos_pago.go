package db

import (
	"context"
	"log"

	"rentals-go/ent"
	"rentals-go/ent/tipopago"
)

// SeedTiposPago inicializa los métodos de pago peruanos por defecto.
func SeedTiposPago(client *ent.Client) error {
	ctx := context.Background()

	metodos := []string{
		"efectivo",
		"transferencia",
		"yape",
		"plin",
		"tarjeta",
		"deposito",
		"otro",
	}

	for _, nombre := range metodos {
		exists, err := client.TipoPago.
			Query().
			Where(tipopago.NombreEQ(nombre)).
			Exist(ctx)
		
		if err != nil {
			return err
		}
		if !exists {
			if _, err := client.TipoPago.
				Create().
				SetNombre(nombre).
				Save(ctx); err != nil {
				return err
			}
			log.Printf("💳 Tipo de pago '%s' inicializado\n", nombre)
		}
	}

	log.Println("✅ Todos los tipos de pago inicializados correctamente")
	return nil
}
