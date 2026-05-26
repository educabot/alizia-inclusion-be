package config

import (
	"math"
	"testing"
)

func TestBoundedUintToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    uint
		expected int
	}{
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "normal value",
			input:    25,
			expected: 25,
		},
		{
			name:     "at max int32 boundary",
			input:    math.MaxInt32,
			expected: math.MaxInt32,
		},
		{
			name:     "above max int32 clamped",
			input:    math.MaxInt32 + 1,
			expected: math.MaxInt32,
		},
		{
			name:     "very large value clamped",
			input:    math.MaxUint32,
			expected: math.MaxInt32,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := boundedUintToInt(tc.input)
			if got != tc.expected {
				t.Errorf("boundedUintToInt(%d) = %d, want %d", tc.input, got, tc.expected)
			}
		})
	}
}
