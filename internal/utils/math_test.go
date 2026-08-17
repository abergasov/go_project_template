package utils_test

import (
	"errors"
	"math"
	"testing"

	"go_project_template/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestSafeIntToUint32(t *testing.T) {
	testCases := []struct {
		name      string
		input     int
		output    uint32
		expectErr bool
	}{
		{
			name:   "valid input",
			input:  10,
			output: 10,
		},
		{
			name:   "zero",
			input:  0,
			output: 0,
		},
		{
			name:   "max uint32",
			input:  math.MaxUint32,
			output: math.MaxUint32,
		},
		{
			name:      "negative",
			input:     -1,
			expectErr: true,
		},
		{
			name:      "overflow",
			input:     math.MaxUint32 + 1,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := utils.SafeIntToUint32(tc.input)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output, output)
			}
		})
	}

	require.Panics(t, func() {
		utils.MustSafeIntToUint32(-1)
	})
	require.NotPanics(t, func() {
		utils.MustSafeIntToUint32(math.MaxUint32)
	})
}

func TestSafeUint64ToInt64(t *testing.T) {
	testCases := []struct {
		name      string
		input     uint64
		output    int64
		expectErr bool
	}{
		{
			name:   "valid input",
			input:  10,
			output: 10,
		},
		{
			name:   "max int64",
			input:  math.MaxInt64,
			output: math.MaxInt64,
		},
		{
			name:      "overflow",
			input:     math.MaxInt64 + 1,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := utils.SafeUint64ToInt64(tc.input)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output, output)
			}
		})
	}

	require.Panics(t, func() {
		utils.MustSafeUint64ToInt64(math.MaxInt64 + 1)
	})
	require.NotPanics(t, func() {
		utils.MustSafeUint64ToInt64(math.MaxInt64)
	})
}

func TestSafeAddUint64Ptr(t *testing.T) {
	for _, tt := range []struct {
		title    string
		initial  uint64
		keyLen   int
		valueLen int
		err      error
		expect   uint64
	}{
		{"valid input", 0, 5, 10, nil, 15},
		{"negative key length", 0, -1, 10, errors.New("value is negative"), 0},
		{"negative value length", 0, 5, -1, errors.New("value is negative"), 0},
		{"integer overflow", 0, int(^uint(0) >> 1), int(^uint(0) >> 1), errors.New("integer overflow detected"), 0},
		{"uint64 counter overflow", math.MaxUint64 - 5, 10, 5, errors.New("uint64 counter overflow"), math.MaxUint64 - 5},
		{"uint64 counter near max", math.MaxUint64 - 20, 10, 5, nil, math.MaxUint64 - 5},
		{"uint64 counter at max", math.MaxUint64, 1, 0, errors.New("uint64 counter overflow"), math.MaxUint64},
	} {
		t.Run(tt.title, func(t *testing.T) {
			counter := tt.initial
			err := utils.SafeAddUint64Ptr(&counter, tt.keyLen, tt.valueLen)
			if tt.err != nil {
				require.Equal(t, tt.err, err, "%s: expected error %v, got %v", tt.title, tt.err, err)
				require.Equal(t, tt.expect, counter, "%s: counter should remain unchanged on error", tt.title)
				return
			}
			require.NoError(t, err, "%s: expected no error", tt.title)
			require.Equal(t, tt.expect, counter, "%s: expected counter to be %d, got %d", tt.title, tt.expect, counter)
		})
	}
}

func TestClamp(t *testing.T) {
	testCases := []struct {
		name     string
		v        float64
		lower    float64
		upper    float64
		expected float64
	}{
		{
			name:     "value within range",
			v:        5.0,
			lower:    0.0,
			upper:    10.0,
			expected: 5.0,
		},
		{
			name:     "value below lower",
			v:        -5.0,
			lower:    0.0,
			upper:    10.0,
			expected: 0.0,
		},
		{
			name:     "value above upper",
			v:        15.0,
			lower:    0.0,
			upper:    10.0,
			expected: 10.0,
		},
		{
			name:     "value equals lower",
			v:        0.0,
			lower:    0.0,
			upper:    10.0,
			expected: 0.0,
		},
		{
			name:     "value equals upper",
			v:        10.0,
			lower:    0.0,
			upper:    10.0,
			expected: 10.0,
		},
		{
			name:     "negative range",
			v:        -3.0,
			lower:    -5.0,
			upper:    -1.0,
			expected: -3.0,
		},
		{
			name:     "negative value below negative lower",
			v:        -10.0,
			lower:    -5.0,
			upper:    -1.0,
			expected: -5.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.Clamp(tc.v, tc.lower, tc.upper)
			require.Equal(t, tc.expected, result, "Clamp(%f, %f, %f) = %f, want %f", tc.v, tc.lower, tc.upper, result, tc.expected)
		})
	}
}
