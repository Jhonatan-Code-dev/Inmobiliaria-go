package money

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Amount representa un monto decimal exacto con 2 decimales.
type Amount struct {
	cents int64
	valid bool
}

func NewAmountFromCents(cents int64) Amount {
	return Amount{cents: cents, valid: true}
}

func ParseAmount(raw string) (Amount, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return Amount{}, fmt.Errorf("monto es obligatorio")
	}

	sign := int64(1)
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = strings.TrimPrefix(s, "-")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return Amount{}, fmt.Errorf("monto inválido")
	}

	wholePart := parts[0]
	if wholePart == "" {
		wholePart = "0"
	}

	whole, err := strconv.ParseInt(wholePart, 10, 64)
	if err != nil {
		return Amount{}, fmt.Errorf("monto inválido")
	}

	frac := "00"
	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return Amount{}, fmt.Errorf("monto debe tener máximo 2 decimales")
		}
		frac = parts[1] + strings.Repeat("0", 2-len(parts[1]))
	}

	decimals, err := strconv.ParseInt(frac, 10, 64)
	if err != nil {
		return Amount{}, fmt.Errorf("monto inválido")
	}

	return Amount{cents: sign * (whole*100 + decimals), valid: true}, nil
}

func MustParseAmount(raw string) Amount {
	amt, err := ParseAmount(raw)
	if err != nil {
		panic(err)
	}
	return amt
}

func AmountFromFloat64(v float64) Amount {
	text := strconv.FormatFloat(v, 'f', 2, 64)
	amt, err := ParseAmount(text)
	if err != nil {
		return Amount{}
	}
	return amt
}

func (a Amount) Cents() int64 {
	return a.cents
}

func (a Amount) Float64() float64 {
	return float64(a.cents) / 100
}

func (a Amount) String() string {
	sign := ""
	cents := a.cents
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}

func (a Amount) MarshalJSON() ([]byte, error) {
	if !a.valid {
		return []byte("null"), nil
	}
	return []byte(a.String()), nil
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) {
		*a = Amount{}
		return nil
	}

	raw := string(trimmed)
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		unquoted, err := strconv.Unquote(raw)
		if err != nil {
			return fmt.Errorf("monto inválido")
		}
		raw = unquoted
	}

	amt, err := ParseAmount(raw)
	if err != nil {
		return err
	}
	*a = amt
	return nil
}
