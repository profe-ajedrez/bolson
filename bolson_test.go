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

	bl := New()
	b.ResetTimer()
	for i, tc := range testBolsonCases {
		b.Run(fmt.Sprintf("case %d/%d", i, len(testBolsonCases)), func(b2 *testing.B) {

			for k := 0; k <= b2.N; k++ {
				_, _ = tc.testCase(&bl)
				b.StopTimer()
				bl.Reset()
				b.StartTimer()
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

			qty, _ := decimal.NewFromString("1")
			maxDiscount, _ := decimal.NewFromString("100")
			bruteWD, _ := decimal.NewFromString("640.0")

			_ = b.discountHandler.AddDiscountFromString("100", discount.Percentual)

			calc, err := b.CalculateFromBruteWD(bruteWD, qty, maxDiscount)

			if err != nil {
				return calc, err
			}

			return calc, err
		},
		expected: `{"withDiscount":{"net":"0","brute":"0","tax":"0","discount":"100","discountedValue":"551.7241379310344828","discountedValueBrute":"640.000000000000000048","unitValue":"0"},"withoutDiscount":{"net":"551.7241379310344828","brute":"640.000000000000000048","tax":"88.275862068965517248","unitValue":"551.7241379310344828"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("19", tax.PercentualMode, tax.OverTaxable)

			qty, _ := decimal.NewFromString("1")
			maxDiscount, _ := decimal.NewFromString("100")
			uv, _ := decimal.NewFromString("42008.403361")

			_ = b.discountHandler.AddDiscountFromString("26.005201", discount.Percentual)

			calc, err := b.Calculate(uv, qty, maxDiscount)

			if err != nil {
				return calc, err
			}

			return calc, err
		},
		expected: `{"withDiscount":{"net":"31084.03363008119439","brute":"36990.0000197966213241","tax":"5905.9663897154269341","discount":"26.005201","discountedValue":"10924.36973091880561","discountedValueBrute":"12999.9999797933786759","unitValue":"31084.03363008119439"},"withoutDiscount":{"net":"42008.403361","brute":"49989.99999959","tax":"7981.59663859","unitValue":"42008.403361"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("16", tax.PercentualMode, tax.OverTaxable)

			qty, _ := decimal.NewFromString("2")
			maxDiscount, _ := decimal.NewFromString("100")
			bruteWD, _ := decimal.NewFromString("2119.999998")

			_ = b.discountHandler.AddDiscountFromString("30.0", discount.Percentual)

			calc, err := b.CalculateFromBruteWD(bruteWD, qty, maxDiscount)

			if err != nil {
				return calc, err
			}

			return calc, err
		},
		expected: `{"withDiscount":{"net":"1279.3103436206896552","brute":"1483.9999986000000000384","tax":"204.6896549793103448384","discount":"30","discountedValue":"548.275861551724138","discountedValueBrute":"635.9999994000000000736","unitValue":"639.6551718103448276"},"withoutDiscount":{"net":"1827.5862051724137932","brute":"2119.999998000000000112","tax":"292.413792827586206912","unitValue":"913.7931025862068966"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("16", tax.PercentualMode, tax.OverTaxable)

			qty, _ := decimal.NewFromString("1")
			maxDiscount, _ := decimal.NewFromString("100")
			bruteWD, _ := decimal.NewFromString("1059.999999")

			_ = b.discountHandler.AddDiscountFromString("30.0", discount.Percentual)

			calc, err := b.CalculateFromBruteWD(bruteWD, qty, maxDiscount)

			if err != nil {
				return calc, err
			}

			return calc, err
		},
		expected: `{"withDiscount":{"net":"639.6551718103448276","brute":"741.9999993000000000192","tax":"102.3448274896551724192","discount":"30","discountedValue":"274.137930775862069","discountedValueBrute":"317.9999997000000000368","unitValue":"639.6551718103448276"},"withoutDiscount":{"net":"913.7931025862068966","brute":"1059.999999000000000056","tax":"146.206896413793103456","unitValue":"913.7931025862068966"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("16", tax.PercentualMode, tax.OverTaxable)
			//_ = b.discountHandler.AddDiscountFromString("10", discount.Percentual)

			qty, _ := decimal.NewFromString("4")
			maxDiscount, _ := decimal.NewFromString("100")
			brute, _ := decimal.NewFromString("311.684804")

			calc, err := b.CalculateFromBrute(brute, qty, maxDiscount)

			if err != nil {
				return calc, err
			}

			return calc, err
		},
		expected: `{"withDiscount":{"net":"268.693796551724138","brute":"311.68480400000000008","tax":"42.99100744827586208","discount":"0","discountedValue":"0","discountedValueBrute":"0","unitValue":"67.1734491379310345"},"withoutDiscount":{"net":"268.693796551724138","brute":"311.68480400000000008","tax":"42.99100744827586208","unitValue":"67.1734491379310345"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("10", tax.PercentualMode, tax.OverTaxable)
			//_ = b.discountHandler.AddDiscountFromString("10", discount.Percentual)

			qty, _ := decimal.NewFromString("10")
			maxDiscount, _ := decimal.NewFromString("100")
			brute, _ := decimal.NewFromString("1100")

			calc, err := b.CalculateFromBrute(brute, qty, maxDiscount)

			if err != nil {
				return calc, err
			}

			return calc, err
		},
		expected: `{"withDiscount":{"net":"1000","brute":"1100","tax":"100","discount":"0","discountedValue":"0","discountedValueBrute":"0","unitValue":"100"},"withoutDiscount":{"net":"1000","brute":"1100","tax":"100","unitValue":"100"}}`,
	},
	{
		testCase: func(b *Bolson) (Bag, error) {
			_ = b.taxHandler.AddTaxFromString("20", tax.PercentualMode, tax.OverTaxable)
			_ = b.discountHandler.AddDiscountFromString("10", discount.Percentual)

			qty, _ := decimal.NewFromString("10.0")
			maxDiscount, _ := decimal.NewFromString("100")
			unitValue, _ := decimal.NewFromString("100")

			calc, err := b.Calculate(unitValue, qty, maxDiscount)

			return calc, err
		},
		expected: `{"withDiscount":{"net":"900","brute":"1080","tax":"180","discount":"10","discountedValue":"100","discountedValueBrute":"120","unitValue":"90"},"withoutDiscount":{"net":"1000","brute":"1200","tax":"200","unitValue":"100"}}`,
	},
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
