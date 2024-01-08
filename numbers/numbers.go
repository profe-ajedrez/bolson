package numbers

import "github.com/shopspring/decimal"

var (
	Hundred = decimal.NewFromInt(100)
	One     = decimal.NewFromInt(1)
	Inverse = decimal.NewFromInt(-1)
	Zero    = decimal.Zero.Copy()
)
