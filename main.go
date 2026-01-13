package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

func main() {
	flag.Parse()
	SetRandomSource(rand.NewSource(time.Now().UnixNano()))

	if len(flag.Args()) == 0 {
		log.Fatal("need to provide 'age [+/-]modifier or a list of die rolls (3d6, 2d8, etc)")
	}

	Roll(flag.Args())
}
