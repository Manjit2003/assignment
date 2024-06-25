package utils

import (
	"fmt"
	"math/rand"
)

var (
	min = 10
	max = 10_000
)

func GenerateRandomCreds() (string, string) {
	num := rand.Intn(max-min) + min
	return fmt.Sprintf("username_%d", num), fmt.Sprintf("password_%d", num)
}
