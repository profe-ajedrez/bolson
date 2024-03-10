package tax

import (
	"fmt"
	"strconv"

	"github.com/profe-ajedrez/bolson/numbers"
	"github.com/shopspring/decimal"
)

// Different types of taxes are represented here
type Mode uint8

const (
	// PercentualMode it's a discount applied as a tasa over a value as when someone says *a discount of 10%*
	PercentualMode = Mode(0)
	// AmountLineMode it's a discount applied as an amount over the entirety of the line without consider quantity, as when someone says *a discount of $10 over the total $100*
	AmountLineMode = Mode(1)
	// AmountUnitMode it's a discount applied as an amount over the value of the unit. considers quantity, as when someone says *a discount of $1 by each of the ten oranges*
	AmountUnitMode = Mode(2)

	// InvalidMode sometimes a way to define an invalid Node could be necessary
	InvalidMode = Mode(99)
)

// String converts Mode to string
func (m Mode) String() string {
	return fmt.Sprintf("%d", m)
}

// NewModeFromInt returns a Mode from int64
func NewModeFromInt(v int64) (Mode, error) {
	if v < 0 || v > 2 {
		return InvalidMode, ErrInvalidTaxMode(v)
	}

	return Mode(v), nil
}

// NewModeFromInt32 returns a Mode from int32
func NewModeFromInt32(v int32) (Mode, error) {
	if v < 0 || v > 2 {
		return InvalidMode, ErrInvalidTaxMode(v)
	}

	return Mode(v), nil
}

// NewModeFromInt16 returns a Mode from int16
func NewModeFromInt16(v int16) (Mode, error) {
	if v < 0 || v > 2 {
		return InvalidMode, ErrInvalidTaxMode(v)
	}

	return Mode(v), nil
}

// NewModeFromInt8 returns a Mode from int8
func NewModeFromInt8(v int8) (Mode, error) {
	if v < 0 || v > 2 {
		return InvalidMode, ErrInvalidTaxMode(v)
	}

	return Mode(v), nil
}

// NewModeFromString returns a Mode from string
func NewModeFromString(v string) (Mode, error) {
	n, err := strconv.Atoi(v)

	if err != nil {
		return InvalidMode, ErrInvalidTaxMode(err)
	}

	if n < 0 || n > 2 {
		return InvalidMode, ErrInvalidTaxMode(n)
	}

	return Mode(n), nil
}

// Represents when a tax should be calculated.
// There are 3 stages in which a tax could be calculated
//
// 1 directly on the values of the products being sold, we call these taxes
// over taxables
//
// 2 on the value obtained from applying overtaxable taxes, we call these
// overtaxes and they are the typical case of tax on tax
//
// 3 are calculated the same as overtaxable taxes, but they are not considered
// for the calculation of overtax taxes, we call these ignorable overtaxes
type Stage uint8

const (
	// OverTaxable Taxes that are calculated directly on the value of the products
	OverTaxable = Stage(0)

	// OverTax Taxes that are calculated on the values of the products plus the
	// over taxable tax calculated for them
	OverTax = Stage(1)

	// OverTaxIgnorable taxes calculated the same as overtaxables, but are not considered
	// for the calculation of overtaxes.
	OverTaxIgnorable = Stage(2)

	InvalidStage = Stage(99)
)

func NewStageFromInt(st int) (Stage, error) {
	switch st {
	case 0:
		return OverTaxable, nil
	case 1:
		return OverTax, nil
	case 2:
		return OverTaxIgnorable, nil
	default:
		return InvalidStage, ErrInvalidTaxStage(st)
	}
}

// NewStageFromInt32 returns a Stage from int32
func NewStageFromInt32(v int32) (Stage, error) {
	if v < 0 || v > 2 {
		return InvalidStage, ErrInvalidTaxStage(v)
	}

	return Stage(v), nil
}

// NewStageFromInt16 returns a Stage from int16
func NewStageFromInt16(v int16) (Stage, error) {
	if v < 0 || v > 2 {
		return InvalidStage, ErrInvalidTaxStage(v)
	}

	return Stage(v), nil
}

// NewStageFromInt8 returns a Stage from int16
func NewStageFromInt8(v int8) (Stage, error) {
	if v < 0 || v > 2 {
		return InvalidStage, ErrInvalidTaxStage(v)
	}

	return Stage(v), nil
}

// NewStageFromString returns a Stage from string
func NewStageFromString(v string) (Stage, error) {
	n, err := strconv.Atoi(v)

	if err != nil {
		return InvalidStage, ErrInvalidTaxStage(err)
	}

	if n < 0 || n > 2 {
		return InvalidStage, ErrInvalidTaxMode(n)
	}

	return Stage(n), nil
}

type Stager interface {
	AddPercentual(decimal.Decimal) error
	AddAmountUnit(decimal.Decimal) error
	AddAmountLine(decimal.Decimal) error

	Tax(decimal.Decimal, decimal.Decimal) (decimal.Decimal, error)

	AddPercentualFromFloat32(float32) error
	AddAmountUnitFromFloat32(float32) error
	AddAmountLineFromFloat32(float32) error

	TaxFromFloat32(float32, float32) (decimal.Decimal, error)

	AddPercentualFromFloat64(float64) error
	AddAmountUnitFromFloat64(float64) error
	AddAmountLineFromFloat64(float64) error

	TaxFromFloat64(float64, float64) (decimal.Decimal, error)

	AddPercentualFromString(string) error
	AddAmountUnitFromString(string) error
	AddAmountLineFromString(string) error

	TaxFromString(string, string) (decimal.Decimal, error)
}

type Storer interface {
	Percent() decimal.Decimal
	AmountUnit() decimal.Decimal
	AmountLine() decimal.Decimal
}

type Untaxer interface {
	Untax(decimal.Decimal, decimal.Decimal) decimal.Decimal
}

var _ Stager = &TaxStage{}
var _ Storer = &TaxStage{}
var _ Untaxer = &TaxStage{}

type TaxStage struct {
	percentuals decimal.Decimal
	amountUnit  decimal.Decimal
	amountLine  decimal.Decimal
	taxable     decimal.Decimal
}

func NewTaxStage() *TaxStage {
	return &TaxStage{
		percentuals: numbers.Zero.Copy(),
		amountUnit:  numbers.Zero.Copy(),
		amountLine:  numbers.Zero.Copy(),
		taxable:     numbers.Zero.Copy(),
	}
}

// AddAmountLine adds a new decimal value as tax to the tax registry
func (ts *TaxStage) AddAmountLine(tax decimal.Decimal) error {
	if tax.IsNegative() {
		return ErrNegativeAmountByLine(tax)
	}

	ts.amountLine = ts.amountLine.Add(tax)
	return nil
}

// AddAmountLineFromFloat32 adds a new float32 value as tax to the tax registry
func (ts *TaxStage) AddAmountLineFromFloat32(tax float32) error {
	return ts.AddAmountLine(decimal.NewFromFloat32(tax))
}

// AddAmountLineFromFloat64 adds a new float64 value as tax to the tax registry
func (ts *TaxStage) AddAmountLineFromFloat64(tax float64) error {
	return ts.AddAmountLine(decimal.NewFromFloat(tax))
}

// AddAmountLineFromString adds a new string value as tax to the tax registry
func (ts *TaxStage) AddAmountLineFromString(tax string) error {
	tx, err := decimal.NewFromString(tax)

	if err != nil {
		return ErrInvalidDecimal(tax)
	}

	return ts.AddAmountLine(tx)
}

// AddAmountUnit adds a new decimal value as tax to the tax registry
func (ts *TaxStage) AddAmountUnit(tax decimal.Decimal) error {
	if tax.IsNegative() {
		return ErrNegativeAmountByUnit(tax)
	}

	ts.amountUnit = ts.amountUnit.Add(tax)
	return nil
}

// AddAmountUnitFromFloat32 adds a new float32 value as tax to the tax registry
func (ts *TaxStage) AddAmountUnitFromFloat32(tax float32) error {
	return ts.AddAmountUnit(decimal.NewFromFloat32(tax))
}

// AddAmountUnitFromFloat64 adds a new float64 value as tax to the tax registry
func (ts *TaxStage) AddAmountUnitFromFloat64(tax float64) error {
	return ts.AddAmountUnit(decimal.NewFromFloat(tax))
}

// AddAmountUnitFromString adds a new string value as tax to the tax registry
func (ts *TaxStage) AddAmountUnitFromString(tax string) error {
	tx, err := decimal.NewFromString(tax)

	if err != nil {
		return ErrInvalidDecimal(tax)
	}

	return ts.AddAmountUnit(tx)
}

// AddPercentual adds a new decimal value as tax to the tax registry
func (ts *TaxStage) AddPercentual(tax decimal.Decimal) error {
	if tax.IsNegative() {
		return ErrNegativePercent(tax)
	}

	ts.percentuals = ts.percentuals.Add(tax)
	return nil
}

// AddPercentualFromFloat32 adds a new float32 value as tax to the tax registry
func (ts *TaxStage) AddPercentualFromFloat32(tax float32) error {
	return ts.AddPercentual(decimal.NewFromFloat32(tax))
}

// AddPercentualFromFloat64 adds a new float64 value as tax to the tax registry
func (ts *TaxStage) AddPercentualFromFloat64(tax float64) error {
	return ts.AddPercentual(decimal.NewFromFloat(tax))
}

// AddPercentualFromString adds a new string value as tax to the tax registry
func (ts *TaxStage) AddPercentualFromString(tax string) error {
	tx, err := decimal.NewFromString(tax)

	if err != nil {
		return ErrInvalidDecimal(tax)
	}

	return ts.AddPercentual(tx)
}

// Tax calculates the recorded taxes of the stage over the received taxable
func (ts *TaxStage) Tax(taxable decimal.Decimal, qty decimal.Decimal) (decimal.Decimal, error) {
	if taxable.IsNegative() {
		return numbers.Zero.Copy(), ErrNegativeTaxable(taxable)
	}

	if qty.IsNegative() {
		return numbers.Zero.Copy(), ErrNegativeTaxable(qty)
	}

	ts.taxable = taxable.Copy()

	return (taxable.Mul(ts.percentuals.Div(numbers.Hundred)).Add(ts.amountUnit)).Mul(qty).Add(ts.amountLine), nil
}

// TaxFromFloat32 implements Stager.
func (ts *TaxStage) TaxFromFloat32(taxable float32, qty float32) (decimal.Decimal, error) {
	return ts.Tax(decimal.NewFromFloat32(taxable), decimal.NewFromFloat32(qty))
}

// TaxFromFloat64 implements Stager.
func (ts *TaxStage) TaxFromFloat64(taxable float64, qty float64) (decimal.Decimal, error) {
	return ts.Tax(decimal.NewFromFloat(taxable), decimal.NewFromFloat(qty))
}

// TaxFromString implements Stager.
func (ts *TaxStage) TaxFromString(taxable string, qty string) (decimal.Decimal, error) {
	tx, err := decimal.NewFromString(taxable)

	if err != nil {
		return numbers.Zero.Copy(), ErrInvalidDecimal(taxable)
	}

	qt, err := decimal.NewFromString(qty)

	if err != nil {
		return numbers.Zero.Copy(), ErrInvalidDecimal(qty)
	}

	return ts.Tax(tx, qt)
}

// Untax implements Untaxer.
func (ts *TaxStage) Untax(taxed decimal.Decimal, qty decimal.Decimal) decimal.Decimal {
	return taxed.Sub(ts.AmountLine()).Sub(ts.AmountUnit().Mul(qty)).Div(numbers.One.Add(ts.Percent().Div(numbers.Hundred)))
}

// AmountLine implements Storer.
func (ts *TaxStage) AmountLine() decimal.Decimal {
	return ts.amountLine.Copy()
}

// AmountUnit implements Storer.
func (ts *TaxStage) AmountUnit() decimal.Decimal {
	return ts.amountUnit.Copy()
}

// Percent implements Storer.
func (ts *TaxStage) Percent() decimal.Decimal {
	return ts.percentuals.Copy()
}

type Handler struct {
	OverTaxables      *TaxStage
	OverTaxes         *TaxStage
	OverTaxIgnorables *TaxStage
}

func NewHandler() *Handler {
	return &Handler{
		OverTaxables:      NewTaxStage(),
		OverTaxes:         NewTaxStage(),
		OverTaxIgnorables: NewTaxStage(),
	}
}

func (h *Handler) AddTax(value decimal.Decimal, mode Mode, stage Stage) error {
	switch stage {
	case OverTaxable:
		switch mode {
		case PercentualMode:
			return h.OverTaxables.AddPercentual(value)
		case AmountLineMode:
			return h.OverTaxables.AddAmountLine(value)
		case AmountUnitMode:
			return h.OverTaxables.AddAmountUnit(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTax:
		switch mode {
		case PercentualMode:
			return h.OverTaxes.AddPercentual(value)
		case AmountLineMode:
			return h.OverTaxes.AddAmountLine(value)
		case AmountUnitMode:
			return h.OverTaxes.AddAmountUnit(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTaxIgnorable:
		switch mode {
		case PercentualMode:
			return h.OverTaxIgnorables.AddPercentual(value)
		case AmountLineMode:
			return h.OverTaxIgnorables.AddAmountLine(value)
		case AmountUnitMode:
			return h.OverTaxIgnorables.AddAmountUnit(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	}

	return ErrInvalidTaxStage(stage)
}

func (h *Handler) AddTaxFromFloat32(value float32, mode Mode, stage Stage) error {
	switch stage {
	case OverTaxable:
		switch mode {
		case PercentualMode:
			return h.OverTaxables.AddPercentualFromFloat32(value)
		case AmountLineMode:
			return h.OverTaxables.AddAmountLineFromFloat32(value)
		case AmountUnitMode:
			return h.OverTaxables.AddAmountUnitFromFloat32(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTax:
		switch mode {
		case PercentualMode:
			return h.OverTaxes.AddPercentualFromFloat32(value)
		case AmountLineMode:
			return h.OverTaxes.AddAmountLineFromFloat32(value)
		case AmountUnitMode:
			return h.OverTaxes.AddAmountUnitFromFloat32(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTaxIgnorable:
		switch mode {
		case PercentualMode:
			return h.OverTaxIgnorables.AddPercentualFromFloat32(value)
		case AmountLineMode:
			return h.OverTaxIgnorables.AddAmountLineFromFloat32(value)
		case AmountUnitMode:
			return h.OverTaxIgnorables.AddAmountUnitFromFloat32(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	}

	return ErrInvalidTaxStage(stage)
}

func (h *Handler) AddTaxFromFloat(value float64, mode Mode, stage Stage) error {
	switch stage {
	case OverTaxable:
		switch mode {
		case PercentualMode:
			return h.OverTaxables.AddPercentualFromFloat64(value)
		case AmountLineMode:
			return h.OverTaxables.AddAmountLineFromFloat64(value)
		case AmountUnitMode:
			return h.OverTaxables.AddAmountUnitFromFloat64(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTax:
		switch mode {
		case PercentualMode:
			return h.OverTaxes.AddPercentualFromFloat64(value)
		case AmountLineMode:
			return h.OverTaxes.AddAmountLineFromFloat64(value)
		case AmountUnitMode:
			return h.OverTaxes.AddAmountUnitFromFloat64(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTaxIgnorable:
		switch mode {
		case PercentualMode:
			return h.OverTaxIgnorables.AddPercentualFromFloat64(value)
		case AmountLineMode:
			return h.OverTaxIgnorables.AddAmountLineFromFloat64(value)
		case AmountUnitMode:
			return h.OverTaxIgnorables.AddAmountUnitFromFloat64(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	}

	return ErrInvalidTaxStage(stage)
}

func (h *Handler) AddTaxFromString(value string, mode Mode, stage Stage) error {
	switch stage {
	case OverTaxable:
		switch mode {
		case PercentualMode:
			return h.OverTaxables.AddPercentualFromString(value)
		case AmountLineMode:
			return h.OverTaxables.AddAmountLineFromString(value)
		case AmountUnitMode:
			return h.OverTaxables.AddAmountUnitFromString(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTax:
		switch mode {
		case PercentualMode:
			return h.OverTaxes.AddPercentualFromString(value)
		case AmountLineMode:
			return h.OverTaxes.AddAmountLineFromString(value)
		case AmountUnitMode:
			return h.OverTaxes.AddAmountUnitFromString(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	case OverTaxIgnorable:
		switch mode {
		case PercentualMode:
			return h.OverTaxIgnorables.AddPercentualFromString(value)
		case AmountLineMode:
			return h.OverTaxIgnorables.AddAmountLineFromString(value)
		case AmountUnitMode:
			return h.OverTaxIgnorables.AddAmountUnitFromString(value)
		default:
			return ErrInvalidTaxMode(mode)
		}
	}

	return ErrInvalidTaxStage(stage)
}

func (h *Handler) Tax(unit_taxable decimal.Decimal, qty decimal.Decimal) (decimal.Decimal, error) {
	overTaxables, err := h.OverTaxables.Tax(unit_taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxes, err := h.OverTaxes.Tax(unit_taxable.Add(overTaxables.Div(qty)), qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxIgnorable, err := h.OverTaxIgnorables.Tax(unit_taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	return overTaxables.Add(overTaxes).Add(overTaxIgnorable), nil
}

func (h *Handler) TaxFromFloat32(taxable float32, qty float32) (decimal.Decimal, error) {
	overTaxables, err := h.OverTaxables.TaxFromFloat32(taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxes, err := h.OverTaxes.TaxFromFloat32(taxable+float32(overTaxables.InexactFloat64()), qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxIgnorable, err := h.OverTaxIgnorables.TaxFromFloat32(taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	return overTaxables.Add(overTaxes).Add(overTaxIgnorable), nil
}

func (h *Handler) TaxFromFloat(taxable float64, qty float64) (decimal.Decimal, error) {
	overTaxables, err := h.OverTaxables.TaxFromFloat64(taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxes, err := h.OverTaxes.TaxFromFloat64(taxable+overTaxables.InexactFloat64(), qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxIgnorable, err := h.OverTaxIgnorables.TaxFromFloat64(taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	return overTaxables.Add(overTaxes).Add(overTaxIgnorable), nil
}

func (h *Handler) TaxFromString(taxable string, qty string) (decimal.Decimal, error) {
	overTaxables, err := h.OverTaxables.TaxFromString(taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	tx, _ := decimal.NewFromString(taxable)

	overTaxes, err := h.OverTaxes.TaxFromString(tx.Add(overTaxables).String(), qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	overTaxIgnorable, err := h.OverTaxIgnorables.TaxFromString(taxable, qty)

	if err != nil {
		return numbers.Zero.Copy(), err
	}

	return overTaxables.Add(overTaxes).Add(overTaxIgnorable), nil
}

func (h *Handler) Untax(brute decimal.Decimal, q decimal.Decimal) (decimal.Decimal, error) {

	if q.LessThanOrEqual(numbers.Zero) {
		return numbers.Zero.Copy(), ErrNegativeQty(fmt.Sprintf("untaxing %v with qty %v", brute, q))
	}

	u1 := h.OverTaxIgnorables.Untax(brute, q)
	//fmt.Printf("u1: %s", u1)
	u2 := h.OverTaxes.Untax(u1, q)
	//fmt.Printf("u2: %s", u2)
	return h.OverTaxables.Untax(u2, q), nil
}

func (h *Handler) LineTax(taxable decimal.Decimal, qty decimal.Decimal, value decimal.Decimal, mode Mode) (decimal.Decimal, error) {

	switch mode {
	case PercentualMode:
		return taxable.Mul(value.Div(numbers.Hundred)), nil
	case AmountLineMode:
		return value.Copy(), nil
	case AmountUnitMode:
		return value.Mul(qty), nil
	}

	return numbers.Zero.Copy(), ErrInvalidTaxMode(mode)
}
