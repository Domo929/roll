# roll

A Go library and CLI for rolling dice using D&D notation.

## Features

- Standard `NdM` notation (e.g., `2d6`, `1d20`)
- Modifiers: `+X` / `-X` (e.g., `1d20+5`)
- Drop lowest/highest: `dl`, `dh` (e.g., `4d6dl1`)
- Keep highest/lowest: `kh`, `kl` (e.g., `2d20kh1`)
- Built-in advantage, disadvantage, and ability score rolls

## Installation

### Library

```bash
go get github.com/Domo929/roll
```

### CLI

```bash
go install github.com/Domo929/roll/cmd/roll@latest
```

## CLI Usage

```bash
roll d20              # Roll a single d20
roll 2d6+3            # Roll 2d6 and add 3
roll 4d6dl1           # Roll 4d6, drop lowest
roll 8d6 2d8          # Roll multiple expressions
roll advantage 5      # Roll with advantage, +5 modifier
roll disadvantage     # Roll with disadvantage
roll stats            # Generate a full set of ability scores
```

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/Domo929/roll"
)

func main() {
    // Parse and roll
    result, _ := roll.RollString("2d6+3")
    fmt.Println(result)

    // Advantage
    result, _ = roll.RollAdvantage(5)
    fmt.Println(result)

    // Ability scores
    result, _ = roll.RollAbilityScore()
    fmt.Println(result)
}
```

## Dice Notation Reference

| Notation   | Meaning                       | Example    |
|------------|-------------------------------|------------|
| `NdM`      | Roll N dice with M sides      | `2d6`      |
| `NdM+X`    | Add X to total                | `1d20+5`   |
| `NdM-X`    | Subtract X from total         | `1d20-2`   |
| `NdMdlN`   | Drop lowest N dice            | `4d6dl1`   |
| `NdMdhN`   | Drop highest N dice           | `4d6dh1`   |
| `NdMkhN`   | Keep highest N dice           | `2d20kh1`  |
| `NdMklN`   | Keep lowest N dice            | `2d20kl1`  |

## Testing

```bash
go test ./... -v -cover
```
