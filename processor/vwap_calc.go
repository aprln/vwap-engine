package processor

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type VWAPCalcDataPoint struct {
	Price decimal.Decimal
	Size  decimal.Decimal
}

func (t VWAPCalcDataPoint) Value() decimal.Decimal {
	return t.Price.Mul(t.Size)
}

func NewVWAPCalc(windowSize int) (*VWAPCalc, error) {
	c := VWAPCalc{
		windowSize:         windowSize,
		dataPoints:         make([]VWAPCalcDataPoint, windowSize),
		oldestDataPointIdx: 0,
		totalValue:         decimal.Zero,
		totalSize:          decimal.Zero,
		vwap:               decimal.Zero,
	}

	if err := c.checkIntegrity(); err != nil {
		return nil, err
	}

	return &c, nil
}

type VWAPCalc struct {
	windowSize         int
	dataPoints         []VWAPCalcDataPoint
	oldestDataPointIdx int
	totalValue         decimal.Decimal
	totalSize          decimal.Decimal
	vwap               decimal.Decimal
}

func (c *VWAPCalc) VWAP() decimal.Decimal {
	return c.vwap
}

func (c *VWAPCalc) AddDataPoint(price, size decimal.Decimal) error {
	if err := c.checkIntegrity(); err != nil {
		return fmt.Errorf("VWAPCalc.AddDataPoint failed data integrity test: %v", err)
	}

	oldDP, newDP := c.replaceOldestDataPoint(price, size)
	c.adjustTotalValue(oldDP, newDP)
	c.adjustTotalSize(oldDP, newDP)
	c.calcVWAP()
	c.adjustOldestDataPointIdx()

	return nil
}

func (c *VWAPCalc) replaceOldestDataPoint(price, size decimal.Decimal) (oldDP, newDP VWAPCalcDataPoint) {
	oldDP = c.dataPoints[c.oldestDataPointIdx]
	newDP = VWAPCalcDataPoint{Price: price, Size: size}
	c.dataPoints[c.oldestDataPointIdx] = newDP
	return oldDP, newDP
}

func (c *VWAPCalc) adjustTotalValue(oldDP, newDP VWAPCalcDataPoint) {
	c.totalValue = c.totalValue.Sub(oldDP.Value()).Add(newDP.Value())
}

func (c *VWAPCalc) adjustTotalSize(oldDP, newDP VWAPCalcDataPoint) {
	c.totalSize = c.totalSize.Sub(oldDP.Size).Add(newDP.Size)
}

func (c *VWAPCalc) calcVWAP() {
	if c.totalSize.IsZero() {
		c.vwap = decimal.Zero

		return
	}

	c.vwap = c.totalValue.Div(c.totalSize)
}

func (c *VWAPCalc) adjustOldestDataPointIdx() {
	c.oldestDataPointIdx += 1
	if c.oldestDataPointIdx == c.windowSize {
		c.oldestDataPointIdx = 0
	}
}

func (c *VWAPCalc) checkIntegrity() error {
	if len(c.dataPoints) == 0 {
		return fmt.Errorf("data points pool is empty")
	}

	if c.windowSize != len(c.dataPoints) {
		return fmt.Errorf("window size %d is not the same as data points size %d", c.windowSize, len(c.dataPoints))
	}

	if c.oldestDataPointIdx < 0 || c.oldestDataPointIdx >= len(c.dataPoints) {
		return fmt.Errorf("oldestDataPointIdx out of range: %d", c.oldestDataPointIdx)
	}

	return nil
}
