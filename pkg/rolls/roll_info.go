package rolls

import "fmt"

type Operation string

const (
	AddOperation      Operation = "+"
	SubtractOperation Operation = "-"
)

func (o Operation) String() string {
	return string(o)
}

type RollInfo struct {
	Operation Operation
	Number    int
	Sides     int
	Flat      int
}

func (ri *RollInfo) String() string {
	if ri.Flat != 0 {
		return fmt.Sprintf("%s%d", ri.Operation, ri.Flat)
	}

	return fmt.Sprintf("%s%dd%d", ri.Operation, ri.Number, ri.Sides)
}
