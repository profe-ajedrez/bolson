package baggins

import (
	"encoding/json"
	"testing"

	"github.com/profe-ajedrez/baggins/discount"
	"github.com/profe-ajedrez/baggins/tax"
	"github.com/shopspring/decimal"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestBaggins(t *testing.T) {
	b := New()

	for _, taxCase := range taxCases {
		b.AddTax(taxCase.value, taxCase.mode, taxCase.stage)
	}

	for _, discCase := range discCases {
		b.AddDiscount(discCase.value, discCase.mode)
	}

	unitValue, _ := decimal.NewFromString("100.0")
	qty, _ := decimal.NewFromString("10.0")
	maxDiscount, _ := decimal.NewFromString("100")

	result, err := b.Calculate(unitValue, qty, maxDiscount)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := `{"withDiscount":{"net":"879","brute":"1247.492","tax":"368.492","discount":"12.1","discountedValue":"121","discountedValueBrute":"157.058","unitValue":"87.9"},"withoutDiscount":{"net":"1000","brute":"1404.55","tax":"404.55","unitValue":"100"}}`

	js, _ := json.Marshal(result)

	if string(js) != expected {
		t.Logf("fails! expected %s, got %s", expected, string(js))
		t.FailNow()
	}
}

// func TestBaggins2(t *testing.T) {
// 	b := New()

// 	b.AddTax(decimal.NewFromInt(10), tax.PercentualMode, tax.OverTaxable)

// 	b.AddDiscount(decimal.NewFromInt(10), discount.Percentual)

// 	unitValue, _ := decimal.NewFromString("100.0")
// 	qty, _ := decimal.NewFromString("10.0")
// 	maxDiscount, _ := decimal.NewFromString("100")

// 	result, err := b.Calculate(unitValue, qty, maxDiscount)

// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}

// 	js, _ := json.Marshal(result)
// 	fmt.Println(string(js))
// }

func BenchmarkBaggins(b *testing.B) {
	bg := New()

	for _, taxCase := range taxCases {
		bg.AddTax(taxCase.value, taxCase.mode, taxCase.stage)
	}

	for _, discCase := range discCases {
		bg.AddDiscount(discCase.value, discCase.mode)
	}

	unitValue, _ := decimal.NewFromString("100.0")
	qty, _ := decimal.NewFromString("10.0")
	maxDiscount, _ := decimal.NewFromString("100")

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		_, _ = bg.Calculate(unitValue, qty, maxDiscount)
	}

}

var discCases = []struct {
	value decimal.Decimal
	mode  discount.Mode
}{
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("10.0")
			return d
		}(),
		mode: discount.Percentual,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("2.0")
			return d
		}(),
		mode: discount.AmountUnit,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("1.0")
			return d
		}(),
		mode: discount.AmountLine,
	},
}

var taxCases = []struct {
	value decimal.Decimal
	mode  tax.Mode
	stage tax.Stage
}{
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("16.0")
			return d
		}(),
		mode:  tax.PercentualMode,
		stage: tax.OverTaxable,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("5")
			return d
		}(),
		mode:  tax.AmountUnitMode,
		stage: tax.OverTaxable,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("1")
			return d
		}(),
		mode:  tax.AmountLineMode,
		stage: tax.OverTaxable,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("5")
			return d
		}(),
		mode:  tax.PercentualMode,
		stage: tax.OverTax,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("2")
			return d
		}(),
		mode:  tax.AmountUnitMode,
		stage: tax.OverTax,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("1")
			return d
		}(),
		mode:  tax.AmountLineMode,
		stage: tax.OverTax,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("8")
			return d
		}(),
		mode:  tax.PercentualMode,
		stage: tax.OverTaxIgnorable,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("3")
			return d
		}(),
		mode:  tax.AmountUnitMode,
		stage: tax.OverTaxIgnorable,
	},
	{
		value: func() decimal.Decimal {
			d, _ := decimal.NewFromString("2")
			return d
		}(),
		mode:  tax.AmountLineMode,
		stage: tax.OverTaxIgnorable,
	},
}
