// Command roll is a CLI tool for rolling dice using D&D notation.
//
// Usage:
//
//	roll <expression> [expressions...]
//	roll advantage [modifier]
//	roll disadvantage [modifier]
//	roll stats
//
// Examples:
//
//	roll 2d6          # Roll 2 six-sided dice
//	roll 1d20+5       # Roll d20 with +5 modifier
//	roll 4d6dl1       # Roll 4d6, drop lowest (ability scores)
//	roll 2d20kh1      # Roll 2d20, keep highest (advantage)
//	roll advantage 5  # Roll with advantage and +5 modifier
//	roll stats        # Generate a full set of ability scores
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Domo929/roll"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	args := os.Args[1:]
	
	switch strings.ToLower(args[0]) {
	case "help", "-h", "--help":
		printUsage()
		return
	case "advantage", "adv":
		handleAdvantage(args[1:])
	case "disadvantage", "dis", "disadv":
		handleDisadvantage(args[1:])
	case "stats", "abilities":
		handleStats()
	default:
		handleRolls(args)
	}
}

func printUsage() {
	fmt.Println(`roll - A D&D dice roller

Usage:
  roll <expression> [expressions...]
  roll advantage [modifier]
  roll disadvantage [modifier]
  roll stats

Dice Notation:
  NdM       Roll N dice with M sides (e.g., 2d6, 1d20)
  NdM+X     Add X to the total (e.g., 1d20+5)
  NdM-X     Subtract X from the total (e.g., 1d20-2)
  NdMdl[N]  Drop lowest N dice (e.g., 4d6dl1)
  NdMdh[N]  Drop highest N dice (e.g., 4d6dh1)
  NdMkh[N]  Keep highest N dice (e.g., 2d20kh1 for advantage)
  NdMkl[N]  Keep lowest N dice (e.g., 2d20kl1 for disadvantage)

Shortcuts:
  advantage [mod]     Roll 2d20 keep highest + modifier
  disadvantage [mod]  Roll 2d20 keep lowest + modifier
  stats               Roll 4d6dl1 six times for ability scores

Examples:
  roll d20            Roll a single d20
  roll 2d6+3          Roll 2d6 and add 3
  roll 4d6dl1         Roll 4d6, drop lowest (ability score)
  roll 8d6 2d8        Roll multiple expressions
  roll advantage 5    Roll with advantage, +5 modifier
  roll stats          Generate ability scores`)
}

func handleAdvantage(args []string) {
	if err := handleAdvantageTo(os.Stdout, args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func handleAdvantageTo(w io.Writer, args []string) error {
	mod := 0
	if len(args) > 0 {
		var err error
		mod, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid modifier: %s", args[0])
		}
	}

	result, err := roll.RollAdvantage(mod)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	fmt.Fprintf(w, "🎲 Advantage: %s\n", result.String())
	return nil
}

func handleDisadvantage(args []string) {
	if err := handleDisadvantageTo(os.Stdout, args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func handleDisadvantageTo(w io.Writer, args []string) error {
	mod := 0
	if len(args) > 0 {
		var err error
		mod, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid modifier: %s", args[0])
		}
	}

	result, err := roll.RollDisadvantage(mod)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	fmt.Fprintf(w, "🎲 Disadvantage: %s\n", result.String())
	return nil
}

func handleStats() {
	handleStatsTo(os.Stdout)
}

func handleStatsTo(w io.Writer) {
	fmt.Fprintln(w, "🎲 Rolling ability scores (4d6 drop lowest):")
	fmt.Fprintln(w)

	abilities := []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"}
	total := 0

	for _, ability := range abilities {
		result, err := roll.RollAbilityScore()
		if err != nil {
			fmt.Fprintf(w, "Error: %v\n", err)
			return
		}
		total += result.Total
		fmt.Fprintf(w, "  %s: %2d  %s\n", ability, result.Total, result.String())
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Total: %d (average: %.1f)\n", total, float64(total)/6.0)
}

func handleRolls(exprs []string) {
	handleRollsTo(os.Stdout, exprs)
}

func handleRollsTo(w io.Writer, exprs []string) {
	for _, expr := range exprs {
		result, err := roll.RollString(expr)
		if err != nil {
			fmt.Fprintf(w, "Error rolling '%s': %v\n", expr, err)
			continue
		}
		fmt.Fprintf(w, "🎲 %s: %s\n", expr, result.String())
	}
}
