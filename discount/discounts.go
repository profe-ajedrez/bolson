package discount

import (
	"fmt"
	"strconv"

	"github.com/profe-ajedrez/bolson/numbers"
	"github.com/shopspring/decimal"
)

// Mode Different types of discounts are represented here
type Mode uint8

const (
	// Percentual it's a discount applied as a tasa over a value as when someone says *a discount of 10%*
	Percentual = Mode(0)
	// AmountLine it's a discount applied as an amount over the entirety of the line without consider quantity, as when someone says *a discount of $10 over the total $100*
	AmountLine = Mode(1)
	// AmountUnit it's a discount applied as an amount over the value of the unit. considers quantity, as when someone says *a discount of $1 by each of the ten oranges*
	AmountUnit = Mode(2)

	// Invalid sometimes a way to define an invalid Node could be necessary
	Invalid = Mode(99)
)

// String converts Mode to string
func (m Mode) String() string {
	return fmt.Sprintf("%d", m)
}

// NewFromInt returns a Mode from int64
func NewFromInt(v int64) (Mode, error) {
	if v < 0 || v > 2 {
		return Invalid, ErrInvalidDiscountMode(v)
	}

	return Mode(v), nil
}

// NewFromInt32 returns a Mode from int32
func NewFromInt32(v int32) (Mode, error) {
	if v < 0 || v > 2 {
		return Invalid, ErrInvalidDiscountMode(v)
	}

	return Mode(v), nil
}

// NewFromInt16 returns a Mode from int16
func NewFromInt16(v int16) (Mode, error) {
	if v < 0 || v > 2 {
		return Invalid, ErrInvalidDiscountMode(v)
	}

	return Mode(v), nil
}

// NewFromInt8 returns a Mode from int8
func NewFromInt8(v int8) (Mode, error) {
	if v < 0 || v > 2 {
		return Invalid, ErrInvalidDiscountMode(v)
	}

	return Mode(v), nil
}

// NewFromString returns a Mode from string
func NewFromString(v string) (Mode, error) {
	n, err := strconv.Atoi(v)

	if err != nil {
		return Invalid, ErrInvalidDiscountMode(err)
	}

	if n < 0 || n > 2 {
		return Invalid, ErrInvalidDiscountMode(n)
	}

	return Mode(n), nil
}

// Something able to calculate discounts
type DiscountComputer interface {
	AddDiscountFromFloat(float64, Mode) error
	AddDiscountFromFloat32(float32, Mode) error
	AddDiscountFromString(string, Mode) error
	AddDiscount(decimal.Decimal, Mode) error

	ComputeFromFloat64(float64, float64, float64) (decimal.Decimal, decimal.Decimal, error)
	ComputeFromFloat32(float32, float32, float32) (decimal.Decimal, decimal.Decimal, error)
	ComputeFromString(string, string, string) (decimal.Decimal, decimal.Decimal, error)
	Compute(decimal.Decimal, decimal.Decimal, decimal.Decimal) (decimal.Decimal, decimal.Decimal, error)

	UnDiscount(decimal.Decimal, decimal.Decimal) (decimal.Decimal, error)
	UnDiscountFromFloat64(float64, float64) (decimal.Decimal, error)
	UnDiscountFromFloat32(float32, float32) (decimal.Decimal, error)
	UnDiscountFromString(string, string) (decimal.Decimal, error)

	Ratio(decimal.Decimal, decimal.Decimal) decimal.Decimal

	Reset()
}

var _ DiscountComputer = &ComputedDiscount{}

// ComputedDiscount implements [DiscountComputer] providing a discount calculator
type ComputedDiscount struct {
	percentual decimal.Decimal
	amountLine decimal.Decimal
	amountUnit decimal.Decimal
}

// NewComputedDiscount returns a new pointer to [ComputedDiscount]
func NewComputedDiscount() *ComputedDiscount {
	return &ComputedDiscount{
		percentual: decimal.Zero.Copy(),
		amountLine: decimal.Zero.Copy(),
		amountUnit: decimal.Zero.Copy(),
	}
}

func (cd *ComputedDiscount) Reset() {
	cd.amountLine = numbers.Zero.Copy()
	cd.amountUnit = numbers.Zero.Copy()
	cd.percentual = numbers.Zero.Copy()
}

// AddDiscount adds a discount to the discounter
func (cd *ComputedDiscount) AddDiscount(d decimal.Decimal, mode Mode) error {
	switch mode {
	case Percentual:
		cd.percentual = cd.percentual.Add(d)
	case AmountLine:
		cd.amountLine = cd.amountLine.Add(d)
	case AmountUnit:
		cd.amountUnit = cd.amountUnit.Add(d)
	default:
		return ErrInvalidDiscountMode(mode)
	}

	return nil
}

// AddDiscountFromFloat adds a discount to the discounter from a float64 value. Some precission may be lost
func (cd *ComputedDiscount) AddDiscountFromFloat(d float64, mode Mode) error {
	return cd.AddDiscount(decimal.NewFromFloat(d), mode)
}

// AddDiscountFromFloat32 adds a discount to the discounter from a float32 value. Some precission may be lost
func (cd *ComputedDiscount) AddDiscountFromFloat32(d float32, mode Mode) error {
	return cd.AddDiscount(decimal.NewFromFloat32(d), mode)
}

// AddDiscountFromString adds a discount to the discounter from a string
func (cd *ComputedDiscount) AddDiscountFromString(d string, mode Mode) error {
	v, err := decimal.NewFromString(d)

	if err != nil {
		return ErrInvalidDecimal(d)
	}

	return cd.AddDiscount(v, mode)
}

// Compute calculates the values of the registered discounts over a discountable value
func (cd *ComputedDiscount) Compute(uv decimal.Decimal, qty decimal.Decimal, maxDiscount decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	if uv.IsNegative() {
		return numbers.Zero.Copy(), numbers.Zero.Copy(), ErrNegativeUnitValue(uv)
	}

	if qty.IsNegative() {
		return numbers.Zero.Copy(), numbers.Zero.Copy(), ErrNegativeQuantity(uv)
	}

	if maxDiscount.IsNegative() || maxDiscount.GreaterThan(numbers.Hundred) {
		maxDiscount = numbers.Hundred.Copy()
	}

	maxDiscountValue := uv.Mul(maxDiscount).Div(numbers.Hundred).Mul(qty)
	discounted := uv.Mul(cd.percentual).Div(numbers.Hundred).Add(cd.amountUnit).Mul(qty).Add(cd.amountLine)

	if discounted.GreaterThan(maxDiscountValue) {
		return numbers.Zero.Copy(), numbers.Zero.Copy(), ErrOverMaxDiscount(fmt.Sprintf("discount: %v  max discount: %v", discounted, maxDiscount))
	}

	var discount decimal.Decimal

	if uv.Equal(numbers.Zero) {
		discount = numbers.Hundred.Copy() // discounted.Mul(numbers.Hundred).Div(uv.Mul(qty))
	} else {
		discount = discounted.Mul(numbers.Hundred).Div(uv.Mul(qty))
	}

	return discounted, discount, nil
}

// ComputeFromFloat32 calculates the values of the registered discounts over a float32 discountable value
func (cd *ComputedDiscount) ComputeFromFloat32(uv float32, qty float32, maxDisocunt float32) (decimal.Decimal, decimal.Decimal, error) {
	return cd.Compute(decimal.NewFromFloat32(uv), decimal.NewFromFloat32(qty), decimal.NewFromFloat32(maxDisocunt))
}

// ComputeFromFloat64 calculates the values of the registered discounts over a float64 discountable value
func (cd *ComputedDiscount) ComputeFromFloat64(uv float64, qty float64, maxDiscount float64) (decimal.Decimal, decimal.Decimal, error) {
	return cd.Compute(decimal.NewFromFloat(uv), decimal.NewFromFloat(qty), decimal.NewFromFloat(maxDiscount))
}

// ComputeFromString calculates the values of the registered discounts over a string discountable value
func (cd *ComputedDiscount) ComputeFromString(uv string, qty string, maxDiscount string) (decimal.Decimal, decimal.Decimal, error) {
	duv, err := decimal.NewFromString(uv)

	if err != nil {
		return numbers.Zero.Copy(), numbers.Zero.Copy(), ErrInvalidDecimal(err)
	}

	dqty, err := decimal.NewFromString(qty)

	if err != nil {
		return numbers.Zero.Copy(), numbers.Zero.Copy(), ErrInvalidDecimal(err)
	}

	dmd, err := decimal.NewFromString(maxDiscount)

	if err != nil {
		return numbers.Zero.Copy(), numbers.Zero.Copy(), ErrInvalidDecimal(err)
	}

	return cd.Compute(duv, dqty, dmd)
}

// Ratio implements DiscountComputer.
func (*ComputedDiscount) Ratio(discounted decimal.Decimal, discount decimal.Decimal) decimal.Decimal {
	return numbers.Hundred.Mul(discount).Div(discounted.Add(discount))
}

// UnDiscount returns the original discountable value. The value to which the registered discounts where applied
//
// If somewhere in the un discount process a negative value is detected, a zero decimal value and the error will be returned
func (cd *ComputedDiscount) UnDiscount(discounted decimal.Decimal, qty decimal.Decimal) (decimal.Decimal, error) {
	original := discounted.Add(cd.amountLine)

	if original.IsNegative() {
		return numbers.Zero.Copy(),
			ErrNegativeDiscountable(
				fmt.Sprintf(
					`when [un_discounting] amount line discounts. 
					discounted: %v   amount_line discount: %v`,
					discounted,
					cd.amountLine,
				),
			)
	}

	original = original.Add(cd.amountUnit.Mul(qty))

	if original.IsNegative() {
		return numbers.Zero.Copy(),
			ErrNegativeDiscountable(
				fmt.Sprintf(
					`when [un_discounting] amount unit discounts. 
					discounted: %v   amount_unit discount: %v   quantity: %v`,
					discounted,
					cd.amountLine,
					qty,
				),
			)
	}

	if !cd.percentual.Equal(numbers.Hundred) {
		original = original.Div((numbers.Hundred.Sub(cd.percentual))).Mul(numbers.Hundred)
	}

	return original, nil
}

// UnDiscountFromFloat32 returns the original float32 discountable value. The value to which the registered discounts where applied
func (cd *ComputedDiscount) UnDiscountFromFloat32(discounted float32, qty float32) (decimal.Decimal, error) {
	discted := decimal.NewFromFloat32(discounted)
	qtydec := decimal.NewFromFloat32(qty)

	return cd.UnDiscount(discted, qtydec)
}

// UnDiscountFromFloat64 returns the original float64 discountable value. The value to which the registered discounts where applied
func (cd *ComputedDiscount) UnDiscountFromFloat64(discounted float64, qty float64) (decimal.Decimal, error) {
	discted := decimal.NewFromFloat(discounted)
	qtydec := decimal.NewFromFloat(qty)

	return cd.UnDiscount(discted, qtydec)
}

// UnDiscountFromString returns the original string discountable value. The value to which the registered discounts where applied
//
// When there were errors at converting strings to decimal values a zero decimal value and the error will be returned
func (cd *ComputedDiscount) UnDiscountFromString(discounted string, qty string) (decimal.Decimal, error) {
	discted, err := decimal.NewFromString(discounted)

	if err != nil {
		return numbers.Zero.Copy(),
			ErrInvalidDecimal(
				fmt.Sprintf(
					`when [un_discounting_from_string] converting string discounted to decimal. 
					discounted: %v quantity: %v`,
					discounted,
					qty,
				),
			)
	}

	qtydec, err := decimal.NewFromString(qty)

	if err != nil {
		return numbers.Zero.Copy(),
			ErrInvalidDecimal(
				fmt.Sprintf(
					`when [un_discounting_from_string] converting string qty to decimal. 
					discounted: %v quantity: %v`,
					discounted,
					qty,
				),
			)
	}

	return cd.UnDiscount(discted, qtydec)
}
