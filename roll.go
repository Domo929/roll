// Package roll provides dice rolling functionality for tabletop RPGs like D&D.
// It supports standard NdM notation, advantage/disadvantage, and modifiers like "drop lowest".
package roll

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Modifier represents special roll modifiers
type Modifier int

const (
	ModNone Modifier = iota
	ModAdvantage
	ModDisadvantage
	ModDropLowest
	ModDropHighest
	ModKeepHighest
	ModKeepLowest
)

// Result holds the outcome of a dice roll
type Result struct {
	Expression string   // Original expression
	Rolls      []int    // Individual die results
	Dropped    []int    // Dice that were dropped
	Kept       []int    // Dice that were kept for the total
	Modifier   int      // Numeric modifier (+/-X)
	Total      int      // Final total
	Sides      int      // Number of sides on the dice
	NumDice    int      // Number of dice rolled
}

// Dice represents a parsed dice expression
type Dice struct {
	NumDice   int
	Sides     int
	Modifier  int      // +/- modifier
	RollMod   Modifier // advantage, disadvantage, drop lowest, etc.
	DropKeep  int      // number to drop/keep
}

// diceRegex matches expressions like "2d6", "4d6dl1", "1d20+5", "2d20kh1"
var diceRegex = regexp.MustCompile(`(?i)^(\d+)?d(\d+)(dl|dh|kh|kl)?(\d+)?([+-]\d+)?$`)

// Parse parses a dice notation string into a Dice struct
func Parse(expr string) (*Dice, error) {
	expr = strings.TrimSpace(strings.ToLower(expr))
	
	matches := diceRegex.FindStringSubmatch(expr)
	if matches == nil {
		return nil, fmt.Errorf("invalid dice expression: %s", expr)
	}

	d := &Dice{
		NumDice: 1,
		RollMod: ModNone,
	}

	// Number of dice (defaults to 1)
	if matches[1] != "" {
		n, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid number of dice: %s", matches[1])
		}
		d.NumDice = n
	}

	// Number of sides
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid number of sides: %s", matches[2])
	}
	d.Sides = sides

	// Drop/keep modifier
	if matches[3] != "" {
		switch matches[3] {
		case "dl":
			d.RollMod = ModDropLowest
		case "dh":
			d.RollMod = ModDropHighest
		case "kh":
			d.RollMod = ModKeepHighest
		case "kl":
			d.RollMod = ModKeepLowest
		}
		
		if matches[4] != "" {
			dk, err := strconv.Atoi(matches[4])
			if err != nil {
				return nil, fmt.Errorf("invalid drop/keep count: %s", matches[4])
			}
			d.DropKeep = dk
		} else {
			d.DropKeep = 1 // default to 1
		}
	}

	// Numeric modifier (+/- X)
	if matches[5] != "" {
		mod, err := strconv.Atoi(matches[5])
		if err != nil {
			return nil, fmt.Errorf("invalid modifier: %s", matches[5])
		}
		d.Modifier = mod
	}

	return d, nil
}

// Roll performs the dice roll and returns the result
func (d *Dice) Roll() (*Result, error) {
	if d.NumDice < 1 {
		return nil, errors.New("must roll at least 1 die")
	}
	if d.Sides < 1 {
		return nil, errors.New("dice must have at least 1 side")
	}

	result := &Result{
		Sides:    d.Sides,
		NumDice:  d.NumDice,
		Modifier: d.Modifier,
	}

	// Roll all dice
	rolls := make([]int, d.NumDice)
	for i := 0; i < d.NumDice; i++ {
		roll, err := rollDie(d.Sides)
		if err != nil {
			return nil, err
		}
		rolls[i] = roll
	}
	result.Rolls = rolls

	// Apply modifiers
	result.Kept, result.Dropped = applyRollModifier(rolls, d.RollMod, d.DropKeep)

	// Calculate total
	total := 0
	for _, v := range result.Kept {
		total += v
	}
	total += d.Modifier
	result.Total = total

	return result, nil
}

// rollDie rolls a single die with the given number of sides
func rollDie(sides int) (int, error) {
	return rand.Intn(sides) + 1, nil
}

// applyRollModifier applies drop/keep modifiers to the rolls
func applyRollModifier(rolls []int, mod Modifier, count int) (kept, dropped []int) {
	if mod == ModNone || count == 0 || count >= len(rolls) {
		return append([]int{}, rolls...), nil
	}

	sorted := make([]int, len(rolls))
	copy(sorted, rolls)
	sort.Ints(sorted)

	switch mod {
	case ModDropLowest:
		dropped = sorted[:count]
		kept = sorted[count:]
	case ModDropHighest:
		kept = sorted[:len(sorted)-count]
		dropped = sorted[len(sorted)-count:]
	case ModKeepHighest:
		kept = sorted[len(sorted)-count:]
		dropped = sorted[:len(sorted)-count]
	case ModKeepLowest:
		kept = sorted[:count]
		dropped = sorted[count:]
	default:
		kept = rolls
	}

	return kept, dropped
}

// RollString parses and rolls a dice expression string
func RollString(expr string) (*Result, error) {
	d, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	result, err := d.Roll()
	if err != nil {
		return nil, err
	}
	result.Expression = expr
	return result, nil
}

// RollAdvantage rolls 2d20 and keeps the higher result
func RollAdvantage(modifier int) (*Result, error) {
	d := &Dice{
		NumDice:  2,
		Sides:    20,
		Modifier: modifier,
		RollMod:  ModKeepHighest,
		DropKeep: 1,
	}
	result, err := d.Roll()
	if err != nil {
		return nil, err
	}
	result.Expression = fmt.Sprintf("2d20kh1%+d (advantage)", modifier)
	return result, nil
}

// RollDisadvantage rolls 2d20 and keeps the lower result
func RollDisadvantage(modifier int) (*Result, error) {
	d := &Dice{
		NumDice:  2,
		Sides:    20,
		Modifier: modifier,
		RollMod:  ModKeepLowest,
		DropKeep: 1,
	}
	result, err := d.Roll()
	if err != nil {
		return nil, err
	}
	result.Expression = fmt.Sprintf("2d20kl1%+d (disadvantage)", modifier)
	return result, nil
}

// RollAbilityScore rolls 4d6 and drops the lowest for ability score generation
func RollAbilityScore() (*Result, error) {
	d := &Dice{
		NumDice:  4,
		Sides:    6,
		RollMod:  ModDropLowest,
		DropKeep: 1,
	}
	result, err := d.Roll()
	if err != nil {
		return nil, err
	}
	result.Expression = "4d6dl1 (ability score)"
	return result, nil
}

// String returns a human-readable representation of the result
func (r *Result) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rolled: %v", r.Rolls))
	if len(r.Dropped) > 0 {
		sb.WriteString(fmt.Sprintf(" (dropped: %v)", r.Dropped))
	}
	if r.Modifier != 0 {
		sb.WriteString(fmt.Sprintf(" %+d", r.Modifier))
	}
	sb.WriteString(fmt.Sprintf(" = %d", r.Total))
	return sb.String()
}
