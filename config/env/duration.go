package env

import (
	"strconv"
	"strings"
	"time"
)

// CustomDuration permite usar:
// - 15m, 2h, 30s (time.ParseDuration)
// - 15d (días)
// - 15 (número crudo = días)
type CustomDuration struct {
	time.Duration
}

func (d *CustomDuration) Decode(value string) error {

	// 1️⃣ Sufijo días: "15d"
	if strings.HasSuffix(value, "d") {
		daysStr := strings.TrimSuffix(value, "d")
		if days, err := strconv.Atoi(daysStr); err == nil {
			d.Duration = time.Duration(days) * 24 * time.Hour
			return nil
		}
	}

	// 2️⃣ Entero puro: "15" → 15 días
	if val, err := strconv.Atoi(value); err == nil {
		d.Duration = time.Duration(val) * 24 * time.Hour
		return nil
	}

	// 3️⃣ Fallback estándar: "15m", "2h"
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	d.Duration = parsed
	return nil
}
