package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnvString(t *testing.T) {
	testCases := []struct {
		name   string
		key    string
		defVal string
		envVal string
		want   string
	}{
		{
			name:   "no env var",
			key:    "BANANA",
			defVal: "banana",
			want:   "banana",
		},
		{
			name:   "with env var",
			key:    "BANANA",
			defVal: "banana",
			envVal: "monkey",
			want:   "monkey",
		},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				if tc.envVal != "" {
					t.Setenv(tc.key, tc.envVal)
				}
				got := LoadEnvString(tc.key, tc.defVal)
				assert.Equal(t, tc.want, got)
			},
		)
	}
}

func TestLoadEnvStringSlice(t *testing.T) {
	testCases := []struct {
		name   string
		key    string
		defVal []string
		envVal string
		want   []string
	}{
		{
			name:   "no env var",
			key:    "BANANA",
			defVal: []string{"banana"},
			want:   []string{"banana"},
		},
		{
			name:   "with env var",
			key:    "BANANA",
			defVal: []string{"banana"},
			envVal: "monkey|banana",
			want:   []string{"monkey", "banana"},
		},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				if tc.envVal != "" {
					t.Setenv(tc.key, tc.envVal)
				}
				got := LoadEnvStringSlice(tc.key, tc.defVal)
				assert.Equal(t, tc.want, got)
			},
		)
	}
}

func TestMustLoadEnvPositiveInt(t *testing.T) {
	testCases := []struct {
		name      string
		key       string
		defVal    int
		envVal    string
		want      int
		wantPanic bool
	}{
		{
			name:      "invalid with string env var",
			key:       "BANANA",
			defVal:    10,
			envVal:    "monkey",
			wantPanic: true,
		},
		{
			name:      "invalid with negative env var",
			key:       "BANANA",
			defVal:    10,
			envVal:    "-2",
			wantPanic: true,
		},
		{
			name:      "invalid with zero env var",
			key:       "BANANA",
			defVal:    10,
			envVal:    "0",
			wantPanic: true,
		},
		{
			name:   "valid with positive env var",
			key:    "BANANA",
			defVal: 10,
			envVal: "200",
			want:   200,
		},
		{
			name:   "valid with no env var",
			key:    "BANANA",
			defVal: 10,
			want:   10,
		},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				if tc.envVal != "" {
					t.Setenv(tc.key, tc.envVal)
				}

				if tc.wantPanic {
					require.Panics(
						t, func() {
							MustLoadEnvPositiveInt(tc.key, tc.defVal)
						},
					)

					return
				}

				got := MustLoadEnvPositiveInt(tc.key, tc.defVal)
				assert.Equal(t, tc.want, got)
			},
		)
	}
}
