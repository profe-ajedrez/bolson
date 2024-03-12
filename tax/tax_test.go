package tax

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewTaxStage(t *testing.T) {
	var _ Stager = NewTaxStage()
}

func TestTaxStageRegistryTaxes(t *testing.T) {
	taxStager := NewTaxStage()

	err := taxStager.AddAmountLineFromFloat64(10.32)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = taxStager.AddAmountUnitFromFloat64(11.35)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = taxStager.AddPercentual(func() decimal.Decimal {
		d, _ := decimal.NewFromString("16.092732673726362323232")
		return d
	}())

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	taxValue, err := taxStager.Tax(func() decimal.Decimal {
		d, _ := decimal.NewFromString("100")
		return d
	}(),
		func() decimal.Decimal {
			q, _ := decimal.NewFromString("10")
			return q
		}(),
	)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := "284.7473267372636"

	if expected != taxValue.String() {
		t.Logf("Fails! expeted %s, got %v", expected, taxValue)
		t.FailNow()
	}

}

func TestTaxHandler(t *testing.T) {
	h := NewHandler()

	err := h.AddTaxFromString("10", AmountLineMode, OverTaxable)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("10", PercentualMode, OverTaxable)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString(
		"10",
		AmountUnitMode,
		OverTaxable,
	)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("1", AmountLineMode, OverTax)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("1.1", PercentualMode, OverTax)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("1.2", AmountUnitMode, OverTax)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("1", AmountLineMode, OverTaxIgnorable)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("5", PercentualMode, OverTaxIgnorable)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = h.AddTaxFromString("0.2", AmountUnitMode, OverTaxIgnorable)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	taxable, _ := decimal.NewFromString("100.0")
	qty, _ := decimal.NewFromString("10.0")

	taxes, err := h.Tax(taxable, qty)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := "289.31"

	if expected != taxes.String() {
		t.Logf("Fails! expected %s  got %v", expected, taxes)
		t.FailNow()
	}

	//fmt.Println(taxable.Mul(qty).Add(taxes))

	originalTaxable, err := h.Untax(func() decimal.Decimal {
		d, _ := decimal.NewFromString("1289.31")
		return d
	}(), qty, FromUv)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	//fmt.Println(originalTaxable)

	expected = "100"

	if expected != originalTaxable.String() {
		t.Logf("Fails! expected %s  got %v", expected, originalTaxable)
		t.FailNow()
	}

}

func BenchmarkTaxStageRegistryTaxes(b *testing.B) {
	taxStager := NewTaxStage()

	_ = taxStager.AddAmountLineFromFloat64(10.32)

	_ = taxStager.AddAmountUnitFromFloat64(11.35)

	_ = taxStager.AddPercentual(func() decimal.Decimal {
		d, _ := decimal.NewFromString("16.092732673726362323232")
		return d
	}())

	d, _ := decimal.NewFromString("100")
	q, _ := decimal.NewFromString("10")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = taxStager.Tax(d, q)
	}

}

func BenchmarkTaxHandlerTax(b *testing.B) {
	handler := NewHandler()

	_ = handler.AddTaxFromFloat(10.32, AmountLineMode, OverTaxable)

	_ = handler.AddTaxFromFloat(11.35, AmountUnitMode, OverTaxable)

	_ = handler.AddTax(func() decimal.Decimal {
		d, _ := decimal.NewFromString("16.092732673726362323232")
		return d
	}(), PercentualMode, OverTaxable)

	_ = handler.AddTaxFromFloat(10.32, AmountLineMode, OverTax)

	_ = handler.AddTaxFromFloat(11.35, AmountUnitMode, OverTax)

	_ = handler.AddTax(func() decimal.Decimal {
		d, _ := decimal.NewFromString("16.092732673726362323232")
		return d
	}(), PercentualMode, OverTax)

	_ = handler.AddTaxFromFloat(10.32, AmountLineMode, OverTaxIgnorable)

	_ = handler.AddTaxFromFloat(11.35, AmountUnitMode, OverTaxIgnorable)

	_ = handler.AddTax(func() decimal.Decimal {
		d, _ := decimal.NewFromString("16.092732673726362323232")
		return d
	}(), PercentualMode, OverTaxIgnorable)

	d, _ := decimal.NewFromString("100")
	q, _ := decimal.NewFromString("10")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = handler.Tax(d, q)
	}

}
