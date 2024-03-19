// Package bolson exposes utils to calculate sales values.
//
// shopspring/decimal is used as backend
package bolson

import (
	"fmt"

	"github.com/profe-ajedrez/bolson/discount"
	"github.com/profe-ajedrez/bolson/numbers"
	"github.com/profe-ajedrez/bolson/tax"
	"github.com/shopspring/decimal"
)

// WithDiscountValues represents the result of operations over sales values with applied discounts
type WithDiscountValues struct {
	// Net is the operation subtotal value without taxes
	Net decimal.Decimal `json:"net"`

	// Brute is the operation subtotal values with taxes
	Brute decimal.Decimal `json:"brute"`

	// Tax is the operation value of the cummulated taxes registered by [bolson]
	Tax decimal.Decimal `json:"tax"`

	// Discount is the cummulated discount percentage
	Discount decimal.Decimal `json:"discount"`

	// DiscountValue is the value of the cummulated discounts registered by [bolson] without taxes
	DiscountedValue decimal.Decimal `json:"discountedValue"`

	// DiscountedValueBrute is the value of the cummulated discounts registered by [bolson] with taxes
	DiscountedValueBrute decimal.Decimal `json:"discountedValueBrute"`

	// UnitValue is the raw unit value recalculated from the subtotals
	UnitValue decimal.Decimal `json:"unitValue"`
}

func (c WithDiscountValues) String() string {
	return fmt.Sprintf(`{
		"Net": %v,
		"Brute": %v,
		"Tax": %v,
		"Discount": %v,
		"DiscountedValue": %v,
		"DiscountedValueBrute": %v,
		"unitValue": %v,
	}`, c.Net, c.Brute, c.Tax, c.Discount, c.DiscountedValue, c.DiscountedValueBrute, c.UnitValue)
}

func (c WithDiscountValues) Round(scale int32) WithDiscountValues {
	return WithDiscountValues{
		Net:                  c.Net.Round(scale),
		Brute:                c.Brute.Round(scale),
		Tax:                  c.Tax.Round(scale),
		Discount:             c.Discount.Round(scale),
		DiscountedValue:      c.DiscountedValue.Round(scale),
		DiscountedValueBrute: c.DiscountedValueBrute.Round(scale),
		UnitValue:            c.UnitValue.Round(scale),
	}
}

// WithoutDiscountValues represents the result of operations over sales values without applied discounts
// as if no discounts were registered.
type WithoutDiscountValues struct {
	// Net is the operation subtotal value without taxes. This time without discount applied
	Net decimal.Decimal `json:"net"`

	// Brute is the operation subtotal values with taxes. This time without discount applied
	Brute decimal.Decimal `json:"brute"`

	// Tax is the operation value of the cummulated taxes registered by [bolson]. This time without discount applied
	Tax decimal.Decimal `json:"tax"`

	// UnitValue is the raw unit value recalculated from the subtotals. This time without discount applied
	UnitValue decimal.Decimal `json:"unitValue"`
}

func (c WithoutDiscountValues) String() string {
	return fmt.Sprintf(`{
		"Net": %v,
		"Brute": %v,
		"Tax": %v,	
		"unitValue": %v,
	}`, c.Net, c.Brute, c.Tax, c.UnitValue)
}

func (c WithoutDiscountValues) Round(scale int32) WithoutDiscountValues {
	return WithoutDiscountValues{
		Net:       c.Net.Round(scale),
		Brute:     c.Brute.Round(scale),
		Tax:       c.Tax.Round(scale),
		UnitValue: c.UnitValue.Round(scale),
	}
}

// Bag is used to contain the result of calculations
type Bag struct {
	// WithDiscount contains the obtained values with discount
	WithDiscount WithDiscountValues `json:"withDiscount"`

	// WithoutDiscount contains the obtained values without discount
	WithoutDiscount WithoutDiscountValues `json:"withoutDiscount"`
}

func (b Bag) String() string {
	return fmt.Sprintf(`{
    	"withDiscount" %v,
	"withoutDiscount": %v
}`, b.WithDiscount.String(), b.WithoutDiscount.String())
}

func (b Bag) Round(scale int32) Bag {
	return Bag{
		WithDiscount:    b.WithDiscount.Round(scale),
		WithoutDiscount: b.WithoutDiscount.Round(scale),
	}
}

// Bolson is the handler provided to perform the sales operations over sales values
//
// Internally Bolson has a handler for taxes and a handler for discounts which
// performs operations and calculations over these concepts.
//
// Bolson can register different types of taxes and discount and is able to
// calculate them correctly.
//
// Bolson uses the concept of stages to the taxes registry and calculations,
// where  a tax can be registered in a particular stage which determines when is calculated.
//
// The taxes stages are:
//
// * OverTaxableStage   represents taxes calculated over its value.
//
// * OverTaxesStage represents taxes calculated over its value plus the cummulated amount of the taxes calculated in the OvertaxableStage
//
// * OverTaxesIgnorableStage represents taxes which are calculated like the taxes of the OverTaxableStage, but are not included in the OVerTaxesStage
//
//	 b := Bolson.New()
//
//	 // adds a percentual tax to the Overtaxable stage
//	 err  := b.AddTax(decimal.NewFromInt(10), tax.PercentualMode, tax.OverTaxableStage)
//
//	 if err != nil {
//		    panic(err) // Remember! Dont Panic!
//	 }
type Bolson struct {
	taxHandler      *tax.Handler
	discountHandler *discount.ComputedDiscount
}

func New() Bolson {
	return Bolson{
		taxHandler:      tax.NewHandler(),
		discountHandler: discount.NewComputedDiscount(),
	}
}

func (b Bolson) OverTaxables() *tax.TaxStage {
	return b.taxHandler.OverTaxables
}

func (b Bolson) OverTaxes() *tax.TaxStage {
	return b.taxHandler.OverTaxes
}

func (b Bolson) OverTaxIgnorables() *tax.TaxStage {
	return b.taxHandler.OverTaxIgnorables
}

func (b Bolson) AddTax(value decimal.Decimal, mode tax.Mode, stage tax.Stage) error {
	return b.taxHandler.AddTax(value, mode, stage)
}

func (b Bolson) AddDiscount(value decimal.Decimal, mode discount.Mode) error {
	return b.discountHandler.AddDiscount(value, mode)
}

func (b Bolson) Untax(taxed decimal.Decimal, qty decimal.Decimal, flow int8) (decimal.Decimal, error) {
	return b.taxHandler.Untax(taxed, qty, flow)
}

func (b Bolson) Tax(taxable decimal.Decimal, qty decimal.Decimal) (decimal.Decimal, error) {
	return b.taxHandler.Tax(taxable, qty)
}

func (b Bolson) Discount(unitValue decimal.Decimal, qty decimal.Decimal, maxDiscount decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	return b.discountHandler.Compute(unitValue, qty, maxDiscount)
}

func (b Bolson) Calculate(unitValue decimal.Decimal, qty decimal.Decimal, maxDiscount decimal.Decimal) (calc Bag, err error) {
	return b.subCalculate(unitValue, qty, maxDiscount, tax.FromUv)
}

func (b Bolson) CalculateFromBruteWD(bruteWD decimal.Decimal, qty decimal.Decimal, maxDiscount decimal.Decimal) (calc Bag, err error) {

	discounted, _, err := b.discountHandler.Compute(bruteWD.Div(qty), qty, maxDiscount)

	if err != nil {
		return
	}

	brute := bruteWD.Sub(discounted)

	calc, err = b.CalculateFromBrute(brute, qty, numbers.Hundred)

	return
}

func (b Bolson) CalculateFromBrute(brute decimal.Decimal, qty decimal.Decimal, maxDiscount decimal.Decimal) (calc Bag, err error) {

	//fmt.Printf("brute: %s\n", brute)

	undiscounted, err := b.discountHandler.UnDiscount(brute, qty)

	if err != nil {
		return
	}

	//fmt.Printf("undiscounted: %s\n", undiscounted)

	//fmt.Printf("sub: %s\n", sub)

	untaxedUnitary, err := b.taxHandler.Untax(undiscounted, qty, tax.FromBrute)

	if err != nil {
		return
	}

	//fmt.Printf("untaxedUnitary: %s\n", untaxedUnitary)

	calc, err = b.subCalculate(untaxedUnitary.Div(qty), qty, numbers.Hundred, tax.FromBrute)

	return
}

func (b Bolson) subCalculate(unitValue decimal.Decimal, qty decimal.Decimal, maxDiscount decimal.Decimal, flow int8) (calc Bag, err error) {
	discounted, discount, err := b.discountHandler.Compute(unitValue, qty, maxDiscount)

	if err != nil {
		return
	}

	tax, err := b.taxHandler.Tax(unitValue.Mul(numbers.Hundred.Sub(discount).Div(numbers.Hundred)), qty)

	if err != nil {
		return
	}

	taxWD, err := b.taxHandler.Tax(unitValue, qty)

	if err != nil {
		return
	}

	calc = calculate(unitValue, qty, discounted, tax, discount, taxWD)

	calc.WithoutDiscount.UnitValue, err = b.taxHandler.Untax(calc.WithoutDiscount.Brute, qty, flow)

	if err != nil {
		err = fmt.Errorf("after try to untax brute to recalculate uv %v", err)
		return
	}

	calc.WithDiscount.UnitValue = calc.WithDiscount.Net.Div(qty)
	calc.WithoutDiscount.UnitValue = calc.WithoutDiscount.Net.Div(qty)

	return
}

func (b Bolson) Reset() {
	b.discountHandler.Reset()
	b.taxHandler.Reset()
}

func calculate(unitValue decimal.Decimal, qty decimal.Decimal, discounted decimal.Decimal, tax decimal.Decimal, discount decimal.Decimal, taxWD decimal.Decimal) (calc Bag) {
	netWD := unitValue.Mul(qty)
	net := netWD.Sub(discounted)

	calc.WithDiscount.Net = net
	calc.WithDiscount.Brute = calc.WithDiscount.Net.Add(tax)
	calc.WithDiscount.Tax = tax
	calc.WithDiscount.Discount = discount
	calc.WithDiscount.DiscountedValue = discounted

	calc.WithoutDiscount.Net = netWD
	calc.WithoutDiscount.Brute = calc.WithoutDiscount.Net.Add(taxWD)
	calc.WithoutDiscount.Tax = taxWD

	calc.WithDiscount.DiscountedValueBrute = calc.WithoutDiscount.Brute.Sub(calc.WithDiscount.Brute)

	return calc
}
