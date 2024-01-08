package discount

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func ExampleComputedDiscount() {
	discounter := NewComputedDiscount()

	// We will register the next discounts:
	//
	// percentual 12%
	// percentual 12.3%
	// amount by unit 2.3
	// amount by line 2.3

	err := discounter.AddDiscount(decimal.NewFromInt(12), Percentual)

	if err != nil {
		panic(err)

	}

	err = discounter.AddDiscountFromFloat(12.3, Percentual)

	if err != nil {
		panic(err)

	}

	err = discounter.AddDiscountFromFloat32(2.3, AmountUnit)

	if err != nil {
		panic(err)

	}

	err = discounter.AddDiscountFromString("this should fail", Percentual)

	if err == nil {
		panic("this should be failed because an invalid decimal value was passed to AddDiscountFromString")

	}

	err = discounter.AddDiscountFromString("2.3", AmountLine)

	if err != nil {
		panic(err)

	}

	// we will apply the registred discounts in discounter to the value 100 over 10 units
	result, _, err := discounter.Compute(decimal.NewFromInt32(100), decimal.NewFromInt(10), decimal.NewFromInt(100))

	if err != nil {
		panic(err)
	}

	// ...so the result will be the value of the total discount to apply:
	//
	// = (unitary * Sum(percentual) + Sum(amount_unit)) * qty + Sum(amount_line)
	//
	// = (100 * (24.3 / 100) + 2.3) * 10 + 2.3
	//
	// = 268.3

	expected, _ := decimal.NewFromString("268.3")

	if !expected.Equal(result) {
		panic(fmt.Sprintf("could'nt calculate correct discount value. Expected %v, got %v", expected, result))
	}

	fmt.Printf("Success!! expected: %v -- got: %v", expected, result)

	// Output:
	// Success!! expected: 268.3 -- got: 268.3

}
