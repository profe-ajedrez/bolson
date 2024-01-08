package discount

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

// func TestDec(t *testing.T) {
// 	f1 := 0.1
// 	f2 := 0.2

// 	d1 := decimal.NewFromFloat(f1)
// 	fmt.Println(d1)

// 	d2 := decimal.NewFromFloat(f2)
// 	fmt.Println(d2)

// 	fmt.Printf("f1 + f2 = %f + %f = %.18f\n", f1, f2, f1+f2)
// 	fmt.Printf("d1 + d2 = %v + %v = %v\n", d1, d2, d1.Add(d2))

// 	f3 := 0.3

// 	fmt.Printf("f3 = %f,  real f3 = %.24f3", f3, f3)

// }

func TestNewComputedDiscount(t *testing.T) {
	var _ DiscountComputer = NewComputedDiscount()
}

func TestAddDiscounts(t *testing.T) {
	_ = discounterTest(t)
}

func TestDiscounter(t *testing.T) {
	// We will register the next discounts:
	//
	// percentual 12%
	// percentual 12.3%
	// amount by unit 2.3
	// amount by line 2.3
	discounter := discounterTest(t)

	// we will apply the registred discounts in discounter to the value 100 over 10 units
	disocunted, discount, err := discounter.Compute(decimal.NewFromInt32(100), decimal.NewFromInt(10), decimal.NewFromInt(100))

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// ...so the result will be the value of the total discount to apply:
	//
	// = (unitary * Sum(percentual) + Sum(amount_unit)) * qty + Sum(amount_line)
	//
	// = (100 * (24.3 / 100) + 2.3) * 10 + 2.3
	//
	// = 268.3

	expected, _ := decimal.NewFromString("268.3")

	if !expected.Equal(disocunted) {
		t.Logf("could'nt calculate correct discount value. Expected %v, got %v", expected, disocunted)
		t.FailNow()
	}

	fmt.Println(discount)
}

func BenchmarkDiscounter(b *testing.B) {
	discounter := discounterTest(b)

	uv := decimal.NewFromInt32(100)
	qty := decimal.NewFromInt(10)
	md := decimal.NewFromInt(100)

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		_, _, _ = discounter.Compute(uv, qty, md)
	}
}

func discounterTest(t testing.TB) *ComputedDiscount {
	discounter := NewComputedDiscount()

	err := discounter.AddDiscount(decimal.NewFromInt(12), Percentual)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = discounter.AddDiscountFromFloat(12.3, Percentual)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = discounter.AddDiscountFromFloat32(2.3, AmountUnit)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = discounter.AddDiscountFromString("this should fail", Percentual)

	if err == nil {
		t.Log("this should be failed because an invalid decimal value was passed to AddDiscountFromString")
		t.FailNow()
	}

	err = discounter.AddDiscountFromString("2.3", AmountLine)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return discounter
}
