package common

import (
	"math/rand"
)

func RandomChoice[T any](slice []T) T {
	rand.Seed(Now().Unix())
	rndIdx := rand.Intn(len(slice))
	return slice[rndIdx]
}
