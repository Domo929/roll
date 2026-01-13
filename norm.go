package main

import (
	"fmt"
)

func normGen(dieGens []*RollInfo) {
	if len(dieGens) == 0 {
		fmt.Println("no die combos provided")
		return
	}

	total := 0
	for _, dieGen := range dieGens {
		rollSum := 0
		resMsg := fmt.Sprintf("%s: ", dieGen)
		
		if dieGen.Flat != 0 {
			rollSum = dieGen.Flat
			resMsg = fmt.Sprintf("%s %d", resMsg, dieGen.Flat)
		} else {
			for i := 0; i < dieGen.Number; i++ {
				roll := result(dieGen.Sides)
				rollSum += roll
				resMsg = fmt.Sprintf("%s %d", resMsg, roll)
			}
		}
		
		if dieGen.Operation == SubtractOperation {
			total -= rollSum
		} else {
			total += rollSum
		}
		
		fmt.Println(resMsg)
	}

	fmt.Println("total: ", total)
}
