package discount

import "fmt"

func ErrNegativeDiscountable(info interface{}) error {
	return fmt.Errorf("[ErrNegativeDiscountable] the value to be discounted is negative. %v", info)
}

func ErrNegativeDiscount(info interface{}) error {
	return fmt.Errorf("[ErrNegativeDiscount] the value to discount is negative. %v", info)
}

func ErrNegativePercent(info interface{}) error {
	return fmt.Errorf("[ErrNegativePercentDiscount] the percent to discount is negative. %v", info)
}

func ErrNegativeAmountByUnit(info interface{}) error {
	return fmt.Errorf("[ErrNegativeAmountByUnitDiscount] the amount by unut discount is negative. %v", info)
}

func ErrNegativeAmountByLine(info interface{}) error {
	return fmt.Errorf("[ErrNegativeAmountByLineDiscount] the amount by line discount is negative. %v", info)
}

func ErrNegativeQuantity(info interface{}) error {
	return fmt.Errorf("[ErrNegativeQuantityDiscount] the quantity is negative. %v", info)
}

func ErrNegativeUnitValue(info interface{}) error {
	return fmt.Errorf("[ErrNegativeUnitValue] the unit value couldnt be negative. %v", info)
}

func ErrOverMaxDiscount(info interface{}) error {
	return fmt.Errorf("[ErrOverMaxDiscount] the value is over the max discount is negative. %v", info)
}

func ErrInvalidDecimal(info interface{}) error {
	return fmt.Errorf("[ErrInvalidDecimal] the value to discount is invalid as decimal. %v", info)
}

func ErrInvalidDiscountMode(info interface{}) error {
	return fmt.Errorf("[ErrInvalidDiscountMode] the mode of the discount is invalid. %v", info)
}

func ErrDiscountOther(info interface{}) error {
	return fmt.Errorf("[ErrOther Discount] there was an error. %v", info)
}
