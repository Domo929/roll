package roll

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected *Dice
		wantErr  bool
	}{
		{
			name: "simple d20",
			expr: "d20",
			expected: &Dice{
				NumDice: 1,
				Sides:   20,
			},
		},
		{
			name: "2d6",
			expr: "2d6",
			expected: &Dice{
				NumDice: 2,
				Sides:   6,
			},
		},
		{
			name: "1d20+5",
			expr: "1d20+5",
			expected: &Dice{
				NumDice:  1,
				Sides:    20,
				Modifier: 5,
			},
		},
		{
			name: "1d20-3",
			expr: "1d20-3",
			expected: &Dice{
				NumDice:  1,
				Sides:    20,
				Modifier: -3,
			},
		},
		{
			name: "4d6dl1 - drop lowest",
			expr: "4d6dl1",
			expected: &Dice{
				NumDice:  4,
				Sides:    6,
				RollMod:  ModDropLowest,
				DropKeep: 1,
			},
		},
		{
			name: "4d6dl - drop lowest default 1",
			expr: "4d6dl",
			expected: &Dice{
				NumDice:  4,
				Sides:    6,
				RollMod:  ModDropLowest,
				DropKeep: 1,
			},
		},
		{
			name: "2d20kh1 - keep highest",
			expr: "2d20kh1",
			expected: &Dice{
				NumDice:  2,
				Sides:    20,
				RollMod:  ModKeepHighest,
				DropKeep: 1,
			},
		},
		{
			name: "2d20kl1 - keep lowest",
			expr: "2d20kl1",
			expected: &Dice{
				NumDice:  2,
				Sides:    20,
				RollMod:  ModKeepLowest,
				DropKeep: 1,
			},
		},
		{
			name: "5d10dh2 - drop highest 2",
			expr: "5d10dh2",
			expected: &Dice{
				NumDice:  5,
				Sides:    10,
				RollMod:  ModDropHighest,
				DropKeep: 2,
			},
		},
		{
			name:    "invalid expression",
			expr:    "not a dice roll",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.expr)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected.NumDice, result.NumDice)
			assert.Equal(t, tt.expected.Sides, result.Sides)
			assert.Equal(t, tt.expected.Modifier, result.Modifier)
			assert.Equal(t, tt.expected.RollMod, result.RollMod)
			assert.Equal(t, tt.expected.DropKeep, result.DropKeep)
		})
	}
}

func TestRoll(t *testing.T) {
	t.Run("basic d20 roll", func(t *testing.T) {
		d := &Dice{NumDice: 1, Sides: 20}
		result, err := d.Roll()
		require.NoError(t, err)
		assert.Len(t, result.Rolls, 1)
		assert.GreaterOrEqual(t, result.Rolls[0], 1)
		assert.LessOrEqual(t, result.Rolls[0], 20)
		assert.Equal(t, result.Rolls[0], result.Total)
	})

	t.Run("2d6 roll", func(t *testing.T) {
		d := &Dice{NumDice: 2, Sides: 6}
		result, err := d.Roll()
		require.NoError(t, err)
		assert.Len(t, result.Rolls, 2)
		for _, r := range result.Rolls {
			assert.GreaterOrEqual(t, r, 1)
			assert.LessOrEqual(t, r, 6)
		}
	})

	t.Run("roll with modifier", func(t *testing.T) {
		d := &Dice{NumDice: 1, Sides: 20, Modifier: 5}
		result, err := d.Roll()
		require.NoError(t, err)
		assert.Equal(t, result.Rolls[0]+5, result.Total)
	})

	t.Run("invalid roll - zero dice", func(t *testing.T) {
		d := &Dice{NumDice: 0, Sides: 20}
		_, err := d.Roll()
		assert.Error(t, err)
	})

	t.Run("invalid roll - zero sides", func(t *testing.T) {
		d := &Dice{NumDice: 1, Sides: 0}
		_, err := d.Roll()
		assert.Error(t, err)
	})
}

func TestApplyRollModifier(t *testing.T) {
	tests := []struct {
		name         string
		rolls        []int
		mod          Modifier
		count        int
		expectedKept []int
		expectedDrop []int
	}{
		{
			name:         "drop lowest 1",
			rolls:        []int{3, 5, 2, 6},
			mod:          ModDropLowest,
			count:        1,
			expectedKept: []int{3, 5, 6},
			expectedDrop: []int{2},
		},
		{
			name:         "drop highest 1",
			rolls:        []int{3, 5, 2, 6},
			mod:          ModDropHighest,
			count:        1,
			expectedKept: []int{2, 3, 5},
			expectedDrop: []int{6},
		},
		{
			name:         "keep highest 1",
			rolls:        []int{3, 5, 2, 6},
			mod:          ModKeepHighest,
			count:        1,
			expectedKept: []int{6},
			expectedDrop: []int{2, 3, 5},
		},
		{
			name:         "keep lowest 1",
			rolls:        []int{3, 5, 2, 6},
			mod:          ModKeepLowest,
			count:        1,
			expectedKept: []int{2},
			expectedDrop: []int{3, 5, 6},
		},
		{
			name:         "no modifier",
			rolls:        []int{3, 5, 2, 6},
			mod:          ModNone,
			count:        0,
			expectedKept: []int{3, 5, 2, 6},
			expectedDrop: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kept, dropped := applyRollModifier(tt.rolls, tt.mod, tt.count)
			assert.ElementsMatch(t, tt.expectedKept, kept)
			assert.ElementsMatch(t, tt.expectedDrop, dropped)
		})
	}
}

func TestRollString(t *testing.T) {
	t.Run("valid expression", func(t *testing.T) {
		result, err := RollString("2d6+3")
		require.NoError(t, err)
		assert.Equal(t, "2d6+3", result.Expression)
		assert.Len(t, result.Rolls, 2)
		assert.Equal(t, 3, result.Modifier)
	})

	t.Run("invalid expression", func(t *testing.T) {
		_, err := RollString("invalid")
		assert.Error(t, err)
	})
}

func TestRollAdvantage(t *testing.T) {
	result, err := RollAdvantage(5)
	require.NoError(t, err)
	assert.Len(t, result.Rolls, 2)
	assert.Len(t, result.Kept, 1)
	assert.Len(t, result.Dropped, 1)
	// The kept value should be the higher of the two
	assert.GreaterOrEqual(t, result.Kept[0], result.Dropped[0])
}

func TestRollDisadvantage(t *testing.T) {
	result, err := RollDisadvantage(5)
	require.NoError(t, err)
	assert.Len(t, result.Rolls, 2)
	assert.Len(t, result.Kept, 1)
	assert.Len(t, result.Dropped, 1)
	// The kept value should be the lower of the two
	assert.LessOrEqual(t, result.Kept[0], result.Dropped[0])
}

func TestRollAbilityScore(t *testing.T) {
	result, err := RollAbilityScore()
	require.NoError(t, err)
	assert.Len(t, result.Rolls, 4)
	assert.Len(t, result.Kept, 3)
	assert.Len(t, result.Dropped, 1)
	// Total should be sum of kept dice
	total := 0
	for _, v := range result.Kept {
		total += v
	}
	assert.Equal(t, total, result.Total)
	// Result should be between 3 and 18
	assert.GreaterOrEqual(t, result.Total, 3)
	assert.LessOrEqual(t, result.Total, 18)
}

func TestResultString(t *testing.T) {
	result := &Result{
		Rolls:    []int{3, 5, 2, 6},
		Kept:     []int{3, 5, 6},
		Dropped:  []int{2},
		Modifier: 2,
		Total:    16,
	}
	str := result.String()
	assert.Contains(t, str, "[3 5 2 6]")
	assert.Contains(t, str, "dropped: [2]")
	assert.Contains(t, str, "+2")
	assert.Contains(t, str, "= 16")
}
