package seedconfigurator

import (
	"math/rand"
	"time"
)

func ConfigureRandomSeed() {
	now := time.Now()
	rand.Seed(now.UnixNano())
}
