package common

import (
	"math/rand"
	"time"
)

func RandomChoice[T any](slice []T) T {
	rand.Seed(time.Now().Unix())
	rndIdx := rand.Intn(len(slice))
	return slice[rndIdx]
}
