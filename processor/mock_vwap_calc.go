package processor

import "github.com/shopspring/decimal"

func NewMockVWAPCalc() MockVWAPCalc {
	return MockVWAPCalc{}
}

type MockVWAPCalc struct {
}

func (m MockVWAPCalc) VWAP() decimal.Decimal {
	return decimal.NewFromFloat(1.1)
}

func (m MockVWAPCalc) AddDataPoint(price, size decimal.Decimal) error {
	return nil
}
