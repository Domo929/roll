package rolls

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

var rng = rand.New(rand.NewSource(0))

func SetRandomSource(source rand.Source) {
	rng = rand.New(source)
}

func Roll(args []string) {
	if args[0] == "age" {
		if err := ageGen(args); err != nil {
			log.Println(err)
		}
		return
	}

	rollInfos, err := parseOperators(args)
	if err != nil {
		log.Fatalln(err)
	}

	normGen(rollInfos)
}

func result(sides int) int {
	return rng.Intn(sides) + 1
}

func parseOperators(args []string) ([]*RollInfo, error) {
	combined := strings.Join(args, "")
	parts := split(combined)

	var rollInfos []*RollInfo
	for _, s := range parts {
		rollInfo, err := toRollInfo(s)
		if err != nil {
			return nil, err
		}
		rollInfos = append(rollInfos, rollInfo)
	}

	return rollInfos, nil
}

func split(str string) []string {
	var base string = str
	var splits []string

	for {
		riStr, remainder := getNextDieStr(base)
		splits = append(splits, riStr)

		if remainder == "" {
			break
		}

		base = remainder
	}

	return splits
}

func getNextDieStr(str string) (string, string) {
	start := 0
	if len(str) > 0 && (str[0] == '+' || str[0] == '-') {
		start = 1
	}
	
	for i := start; i < len(str); i++ {
		if str[i] == '+' || str[i] == '-' {
			return str[0:i], str[i:]
		}
	}

	return str, ""
}

func toRollInfo(str string) (*RollInfo, error) {
	if str == "" {
		return nil, fmt.Errorf("empty dice specification")
	}

	var (
		rollInfo RollInfo
		offset   int
		err      error
	)

	switch str[0] {
	case '+':
		rollInfo.Operation = AddOperation
		offset = 1
	case '-':
		rollInfo.Operation = SubtractOperation
		offset = 1
	default:
		rollInfo.Operation = AddOperation
	}

	if offset >= len(str) {
		return nil, fmt.Errorf("invalid dice specification: %s", str)
	}

	parts := strings.Split(str[offset:], "d")
	if len(parts) == 1 {
		rollInfo.Flat, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid flat modifier: %s", str)
		}
	} else if len(parts) == 2 {
		rollInfo.Number, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid number of dice: %s", str)
		}
		rollInfo.Sides, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid number of sides: %s", str)
		}
		
		if rollInfo.Number < 0 {
			return nil, fmt.Errorf("number of dice cannot be negative: %d", rollInfo.Number)
		}
		if rollInfo.Number == 0 {
			return nil, fmt.Errorf("number of dice must be at least 1: got %d", rollInfo.Number)
		}
		if rollInfo.Sides <= 0 {
			return nil, fmt.Errorf("dice must have at least 1 side: %d", rollInfo.Sides)
		}
	} else {
		return nil, fmt.Errorf("invalid dice specification: %s", str)
	}

	return &rollInfo, nil
}
