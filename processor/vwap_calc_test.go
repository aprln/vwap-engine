package processor

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)

func TestTradeDataPoint_Value(t *testing.T) {
	testCases := []struct {
		price float64
		size  float64
		want  float64
	}{
		{0, 2, 0},
		{5.234, 0, 0},
		{1, 2, 2},
		{1.4, 1.56, 2.184},
		{-1.56, 1.4, -2.184},
		{-1.56, -1.4, 2.184},
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("%v", tc), func(tt *testing.T) {
				got := VWAPCalcDataPoint{
					Price: decimal.NewFromFloat(tc.price),
					Size:  decimal.NewFromFloat(tc.size),
				}.Value()
				assert.Equal(tt, decimal.NewFromFloat(tc.want).String(), got.String())
			},
		)
	}
}

func TestNewVWAPCalc(t *testing.T) {
	tests := []struct {
		name       string
		windowSize int
		want       VWAPCalc
		wantErr    bool
	}{
		{
			name:    "invalid",
			wantErr: true,
		},
		{
			name:       "valid",
			windowSize: 200,
			want: VWAPCalc{
				windowSize:         200,
				dataPoints:         make([]VWAPCalcDataPoint, 200),
				oldestDataPointIdx: 0,
				totalValue:         decimal.Zero,
				totalSize:          decimal.Zero,
				vwap:               decimal.Zero,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c, err := NewVWAPCalc(tt.windowSize)

				if tt.wantErr {
					require.Error(t, err)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, tt.want, *c)
			},
		)
	}
}

func TestVWAPCalc_AddDataPoint(t *testing.T) {
	testCases := []struct {
		name                   string
		windowSize             int
		dataPoints             []VWAPCalcDataPoint
		oldestDataPointIdx     int
		totalValue             float64
		totalSize              float64
		vwap                   float64
		inputPrice             float64
		inputSize              float64
		wantTotalValue         float64
		wantTotalSize          float64
		wantVWAP               float64
		wantOldestDataPointIdx float64
		wantErr                bool
	}{
		{
			name:    "empty data points",
			wantErr: true,
		},
		{
			name:       "invalid window size greater than data points size",
			dataPoints: make([]VWAPCalcDataPoint, 3),
			windowSize: 1,
			wantErr:    true,
		},
		{
			name:       "invalid window size less than data points size",
			windowSize: 1,
			dataPoints: make([]VWAPCalcDataPoint, 3),
			wantErr:    true,
		},
		{
			name:               "invalid oldest data point idx",
			windowSize:         3,
			dataPoints:         make([]VWAPCalcDataPoint, 3),
			oldestDataPointIdx: 3,
			wantErr:            true,
		},
		{
			name:                   "valid with zero values",
			windowSize:             3,
			dataPoints:             make([]VWAPCalcDataPoint, 3),
			oldestDataPointIdx:     1,
			wantOldestDataPointIdx: 2,
		},
		{
			name:       "valid with equal values",
			windowSize: 3,
			dataPoints: []VWAPCalcDataPoint{
				{
					Price: decimal.NewFromFloat(1.4),
					Size:  decimal.NewFromFloat(1.56),
				},
				{
					Price: decimal.NewFromFloat(1.4),
					Size:  decimal.NewFromFloat(1.56),
				},
				{
					Price: decimal.NewFromFloat(1.4),
					Size:  decimal.NewFromFloat(1.56),
				},
			},
			oldestDataPointIdx:     0,
			totalValue:             6.552,
			totalSize:              4.68,
			vwap:                   1.4,
			inputPrice:             1.4,
			inputSize:              1.56,
			wantTotalValue:         6.552,
			wantTotalSize:          4.68,
			wantVWAP:               1.4,
			wantOldestDataPointIdx: 1,
		},
		{
			name:       "valid with 2 values and oldest data idx is 0",
			windowSize: 2,
			dataPoints: []VWAPCalcDataPoint{
				{
					Price: decimal.NewFromFloat(1),
					Size:  decimal.NewFromFloat(1),
				},
				{
					Price: decimal.NewFromFloat(2),
					Size:  decimal.NewFromFloat(2),
				},
			},
			oldestDataPointIdx:     0,
			totalValue:             5,
			totalSize:              3,
			vwap:                   1.666666667,
			inputPrice:             3,
			inputSize:              3,
			wantTotalValue:         13,
			wantTotalSize:          5,
			wantVWAP:               2.6,
			wantOldestDataPointIdx: 1,
		},
		{
			name:       "valid with 2 values and oldest data idx is 1",
			windowSize: 2,
			dataPoints: []VWAPCalcDataPoint{
				{
					Price: decimal.NewFromFloat(1),
					Size:  decimal.NewFromFloat(1),
				},
				{
					Price: decimal.NewFromFloat(2),
					Size:  decimal.NewFromFloat(2),
				},
			},
			oldestDataPointIdx:     1,
			totalValue:             5,
			totalSize:              3,
			vwap:                   1.666666667,
			inputPrice:             3,
			inputSize:              3,
			wantTotalValue:         10,
			wantTotalSize:          4,
			wantVWAP:               2.5,
			wantOldestDataPointIdx: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				c := VWAPCalc{
					windowSize:         tc.windowSize,
					dataPoints:         tc.dataPoints,
					oldestDataPointIdx: tc.oldestDataPointIdx,
					totalValue:         decimal.NewFromFloat(tc.totalValue),
					totalSize:          decimal.NewFromFloat(tc.totalSize),
					vwap:               decimal.NewFromFloat(tc.vwap),
				}

				err := c.AddDataPoint(decimal.NewFromFloat(tc.inputPrice), decimal.NewFromFloat(tc.inputSize))

				if tc.wantErr {
					require.Error(t, err)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, decimal.NewFromFloat(tc.wantTotalValue).String(), c.totalValue.String())
				assert.Equal(t, decimal.NewFromFloat(tc.wantTotalSize).String(), c.totalSize.String())
				assert.Equal(t, decimal.NewFromFloat(tc.wantVWAP).String(), c.vwap.String())
			},
		)
	}
}

func TestVWAPCalc_VWAPFlow(t *testing.T) {
	c, e := NewVWAPCalc(3)
	require.NoError(t, e)
	require.Equal(t, "0", c.VWAP().String(), "failed at init")

	sequence := []struct {
		price    float64
		size     float64
		wantVWAP string
	}{
		{1.1, 1.1, "1.1"},
		{2.2, 2.2, "1.8333333333333333"},
		{3.3, 3.3, "2.5666666666666667"},
		{4.4, 4.4, "3.5444444444444444"},
		{5.5, 5.5, "4.5833333333333333"},
		{6.6, 6.6, "5.6466666666666667"},
		{7.7, 7.7, "6.7222222222222222"},
	}

	for i, s := range sequence {
		e = c.AddDataPoint(decimal.NewFromFloat(s.price), decimal.NewFromFloat(s.size))
		require.NoError(t, e)
		require.Equalf(t, s.wantVWAP, c.VWAP().String(), "failed at step %d", i)
	}
}
