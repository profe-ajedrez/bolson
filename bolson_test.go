package bolson

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/profe-ajedrez/bolson/discount"
	"github.com/profe-ajedrez/bolson/tax"
	"github.com/shopspring/decimal"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestBolson(t *testing.T) {

	for i, tc := range testBolsonCases {
		b := New()
		calc, err := tc.testCase(&b)

		js, _ := json.Marshal(calc)
		strJs := string(js)
		t.Log(strJs)

		if err != nil {
			t.Logf("Fail test case[%d] --- %v", i, err)
			t.FailNow()
		}

		if strJs != tc.expected {
			t.Logf("Fail test case[%d] --- expected %v --- got %v", i, tc.expected, strJs)
			t.FailNow()
		}

	}
}

func BenchmarkBolson(b *testing.B) {

	b.ResetTimer()
	for i, tc := range testBolsonCases {
		b.Run(fmt.Sprintf("case %d/%d", i, len(testBolsonCases)), func(b2 *testing.B) {

			for k := 0; k <= b2.N; k++ {
				bl := New()
				_, _ = tc.testCase(&bl)

			}
		})
	}

}

var testBolsonCases = []struct {
	testCase func(b *Bolson) (Bag, error)
	expected string
}{
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("16", tax.PercentualMode, tax.OverTaxable)
			_ = b.discountHandler.AddDiscountFromString("30.1885553573578", discount.Percentual)

			unitValue, _ := decimal.NewFromString("913.793103448276")
			qty, _ := decimal.NewFromString("1.0")
			maxDiscount, _ := decimal.NewFromString("100")

			return b.Calculate(unitValue, qty, maxDiscount)
		},
		expected: `{"withDiscount":{"net":"637.9321665620753722","brute":"740.00131321200743174459467975552","tax":"102.06914664993205954459467975552","discount":"30.1885553573578","discountedValue":"275.8609368862006278","discountedValueBrute":"319.99868678799272825540532024448","unitValue":"637.9321665620753722"},"withoutDiscount":{"net":"913.793103448276","brute":"1060.00000000000016","tax":"146.20689655172416","unitValue":"913.793103448276"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("16", tax.PercentualMode, tax.OverTaxable)
			_ = b.discountHandler.AddDiscountFromString("30.1885553573578", discount.Percentual)

			qty, _ := decimal.NewFromString("1.0")
			maxDiscount, _ := decimal.NewFromString("100")
			brute, _ := decimal.NewFromString("740")

			return b.CalculateFromBrute(brute, qty, maxDiscount)
		},
		expected: `{"withDiscount":{"net":"637.9310344827586222","brute":"740.000000000000001756274132676085568","tax":"102.068965517241379556274132676085568","discount":"30.1885553573578","discountedValue":"275.8604473412657312","discountedValueBrute":"319.998118915868248187725867323914432","unitValue":"637.9310344827586222"},"withoutDiscount":{"net":"913.7914818240243534","brute":"1059.998118915868249944","tax":"146.206637091843896544","unitValue":"913.7914818240243534"}}`,
	},
}
