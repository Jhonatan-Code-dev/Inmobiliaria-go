package money

import "testing"

func TestParseAmount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		raw   string
		cents int64
	}{
		{raw: "50", cents: 5000},
		{raw: "50.5", cents: 5050},
		{raw: "50.50", cents: 5050},
		{raw: "\"50.50\"", cents: 5050},
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			t.Parallel()

			var amt Amount
			if err := amt.UnmarshalJSON([]byte(tc.raw)); err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if amt.Cents() != tc.cents {
				t.Fatalf("Cents() = %d, want %d", amt.Cents(), tc.cents)
			}
			if got := amt.String(); got == "" {
				t.Fatalf("String() should not be empty")
			}
		})
	}
}

func TestParseAmountRejectsTooManyDecimals(t *testing.T) {
	t.Parallel()

	if _, err := ParseAmount("50.505"); err == nil {
		t.Fatal("expected error for more than 2 decimals")
	}
}
