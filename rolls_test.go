package main

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOperators(t *testing.T) {
	var tests = []struct {
		name     string
		input    []string
		expected []*RollInfo
		wantErr  bool
	}{
		{
			name: "simple d6",
			input: []string{
				"1d6",
			},
			expected: []*RollInfo{
				{
					Operation: AddOperation,
					Number:    1,
					Sides:     6,
				},
			},
		},
		{
			name: "simple d6 + modifier",
			input: []string{
				"1d6",
				"+2",
			},
			expected: []*RollInfo{
				{
					Operation: AddOperation,
					Number:    1,
					Sides:     6,
				},
				{
					Operation: AddOperation,
					Flat:      2,
				},
			},
		},
		{
			name: "complex expression",
			input: []string{
				"3d6+2d8-5",
			},
			expected: []*RollInfo{
				{
					Operation: AddOperation,
					Number:    3,
					Sides:     6,
				},
				{
					Operation: AddOperation,
					Number:    2,
					Sides:     8,
				},
				{
					Operation: SubtractOperation,
					Flat:      5,
				},
			},
		},
		{
			name: "subtraction first",
			input: []string{
				"-2d4",
			},
			expected: []*RollInfo{
				{
					Operation: SubtractOperation,
					Number:    2,
					Sides:     4,
				},
			},
		},
		{
			name:    "invalid dice sides",
			input:   []string{"2d0"},
			wantErr: true,
		},
		{
			name:    "invalid dice count zero",
			input:   []string{"0d6"},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   []string{""},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseOperators(test.input)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, len(test.expected), len(actual))
			for index, exp := range test.expected {
				assert.Equal(t, *exp, *actual[index])
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple dice",
			input:    "2d6",
			expected: []string{"2d6"},
		},
		{
			name:     "dice with addition",
			input:    "2d6+3",
			expected: []string{"2d6", "+3"},
		},
		{
			name:     "dice with subtraction",
			input:    "2d6-3",
			expected: []string{"2d6", "-3"},
		},
		{
			name:     "complex expression",
			input:    "3d6+2d8-5+1d4",
			expected: []string{"3d6", "+2d8", "-5", "+1d4"},
		},
		{
			name:     "starting with operator",
			input:    "-2d4+3",
			expected: []string{"-2d4", "+3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := split(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestGetNextDieStr(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedStr       string
		expectedRemainder string
	}{
		{
			name:              "no operators",
			input:             "2d6",
			expectedStr:       "2d6",
			expectedRemainder: "",
		},
		{
			name:              "with addition",
			input:             "2d6+3",
			expectedStr:       "2d6",
			expectedRemainder: "+3",
		},
		{
			name:              "starting with plus",
			input:             "+3",
			expectedStr:       "+3",
			expectedRemainder: "",
		},
		{
			name:              "starting with minus",
			input:             "-2d4",
			expectedStr:       "-2d4",
			expectedRemainder: "",
		},
		{
			name:              "starting with minus and more",
			input:             "-2d4+3",
			expectedStr:       "-2d4",
			expectedRemainder: "+3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str, remainder := getNextDieStr(test.input)
			assert.Equal(t, test.expectedStr, str)
			assert.Equal(t, test.expectedRemainder, remainder)
		})
	}
}

func TestToRollInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *RollInfo
		wantErr  bool
	}{
		{
			name:  "simple dice",
			input: "2d6",
			expected: &RollInfo{
				Operation: AddOperation,
				Number:    2,
				Sides:     6,
			},
		},
		{
			name:  "dice with plus",
			input: "+3d8",
			expected: &RollInfo{
				Operation: AddOperation,
				Number:    3,
				Sides:     8,
			},
		},
		{
			name:  "dice with minus",
			input: "-2d4",
			expected: &RollInfo{
				Operation: SubtractOperation,
				Number:    2,
				Sides:     4,
			},
		},
		{
			name:  "flat modifier",
			input: "+5",
			expected: &RollInfo{
				Operation: AddOperation,
				Flat:      5,
			},
		},
		{
			name:  "flat modifier negative",
			input: "-3",
			expected: &RollInfo{
				Operation: SubtractOperation,
				Flat:      3,
			},
		},
		{
			name:    "invalid - zero sides",
			input:   "2d0",
			wantErr: true,
		},
		{
			name:    "invalid - negative sides",
			input:   "2d-6",
			wantErr: true,
		},
		{
			name:    "invalid - empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid - just operator",
			input:   "+",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := toRollInfo(test.input)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, *test.expected, *actual)
		})
	}
}

func TestParseModifier(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedVal  int
		expectedNeg  bool
		wantErr      bool
	}{
		{
			name:        "zero",
			input:       "0",
			expectedVal: 0,
			expectedNeg: false,
		},
		{
			name:        "positive with plus",
			input:       "+5",
			expectedVal: 5,
			expectedNeg: false,
		},
		{
			name:        "negative with minus",
			input:       "-3",
			expectedVal: 3,
			expectedNeg: true,
		},
		{
			name:        "positive without operator",
			input:       "7",
			expectedVal: 7,
			expectedNeg: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectedVal: 0,
			expectedNeg: false,
		},
		{
			name:    "just operator",
			input:   "+",
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   "+abc",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, neg, err := parseModifier(test.input)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedVal, val)
			assert.Equal(t, test.expectedNeg, neg)
		})
	}
}

func TestResult(t *testing.T) {
	// Set a predictable seed for testing
	SetRandomSource(rand.NewSource(42))
	
	// Test that results are within bounds
	for sides := 1; sides <= 20; sides++ {
		testName := "d" + strconv.Itoa(sides)
		t.Run(testName, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				roll := result(sides)
				assert.GreaterOrEqual(t, roll, 1, "roll should be at least 1")
				assert.LessOrEqual(t, roll, sides, "roll should be at most %d", sides)
			}
		})
	}
}
