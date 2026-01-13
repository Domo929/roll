package main

import (
	"fmt"
	"strconv"
)

const (
	AgeDiceCount = 3
	AgeDiceSides = 6
)

func ageGen(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("age command requires arguments")
	}

	modifier := "0"
	if len(args) >= 2 {
		modifier = args[1]
	}
	modVal, negative, err := parseModifier(modifier)
	if err != nil {
		return fmt.Errorf("invalid modifier: %w", err)
	}

	dice := make([]int, 0, AgeDiceCount)
	sum := 0
	for i := 0; i < AgeDiceCount; i++ {
		roll := result(AgeDiceSides)
		dice = append(dice, roll)
		sum += roll
	}
	stuntPoints := dice[0] == dice[1] || dice[0] == dice[2] || dice[1] == dice[2]

	fmt.Printf("Dice: %d %d *%d*\n", dice[0], dice[1], dice[2])
	if negative {
		sum -= modVal
		fmt.Printf("Modifier: -%d\n", modVal)
	} else {
		sum += modVal
		fmt.Printf("Modifier: +%d\n", modVal)
	}

	fmt.Println("Total: ", sum)
	if stuntPoints {
		fmt.Printf("Generated %d stunt points\n", dice[2])
	}
	if dice[2] == AgeDiceSides {
		fmt.Println("Rolled a 6 on your drama die")
	}

	return nil
}

func parseModifier(modifier string) (int, bool, error) {
	if modifier == "" {
		return 0, false, nil
	}

	var (
		numIndexStart = 0
		negative      = false
	)

	if len(modifier) >= 2 && (modifier[0] == '+' || modifier[0] == '-') {
		numIndexStart = 1
		if modifier[0] == '-' {
			negative = true
		}
	}

	if numIndexStart >= len(modifier) {
		return 0, false, fmt.Errorf("modifier has no numeric value")
	}

	val, err := strconv.ParseInt(modifier[numIndexStart:], 10, 64)
	if err != nil {
		return 0, false, err
	}

	return int(val), negative, nil
}
