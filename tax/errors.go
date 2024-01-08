package tax

import "fmt"

// ErrNegativeTaxable the value over which calculate tax is negative
func ErrNegativeTaxable(info any) error {
	return fmt.Errorf("[ErrNegativeTaxable] the specified values is negative. a taxable cannot be negative %v", info)
}

// ErrNegativeTax the calculated tax is negative
func ErrNegativeTax(info any) error {
	return fmt.Errorf("[ErrNegativeTax] the specified values is negative. a tax cannot be negative %v", info)
}

// ErrNegativeQty the quatity being sold is negative
func ErrNegativeQty(info any) error {
	return fmt.Errorf("[ErrNegativeQty tax] the specified values is negative. a quantity cannot be negative %v", info)
}

// ErrNegativePercent the percentage of the tax is negative
func ErrNegativePercent(info any) error {
	return fmt.Errorf("[ErrNegativePercent tax] the specified values is negative. a percentual tax cannot be negative %v", info)
}

// ErrNegativeAmountByUnit the amoun by unit of the tax is negative
func ErrNegativeAmountByUnit(info any) error {
	return fmt.Errorf("[ErrNegativeAmountByUnit tax] the specified values is negative. an amount tax cannot be negative %v", info)
}

// ErrNegativeAmountByLine the amoun by line of the tax is negative
func ErrNegativeAmountByLine(info any) error {
	return fmt.Errorf("[ErrNegativeAmountByLine tax] the specified values is negative. an amount line tax cannot be negative %v", info)
}

// ErrInvalidDecimal a calculation o convertion produced an invalid  Bigdecimal value
func ErrInvalidDecimal(info any) error {
	return fmt.Errorf("[ErrInvalidDecimal tax] the specified values is not a valid decimal value. %v", info)
}

// ErrInvalidTaxStage the tax stage not exists
func ErrInvalidTaxStage(info any) error {
	return fmt.Errorf("[ErrInvalidTaxStage] the specified tax stage doesnt exists. %v", info)
}

// ErrInvalidTaxMode the tax stage not exists
func ErrInvalidTaxMode(info any) error {
	return fmt.Errorf("[ErrInvalidTaxMode] the specified tax mode doesnt exists. %v", info)
}

// ErrOther other error
func ErrOther(info any) error {
	return fmt.Errorf("[ErrOther Tax] there was an error. %v", info)
}
