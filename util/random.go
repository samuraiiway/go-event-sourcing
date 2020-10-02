package util

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GetUUID() string {
	now := time.Now().UnixNano()
	ran := random.Uint64()
	return fmt.Sprintf("%x-%x", now, ran)
}

func RandomInt(n int) int {
	return random.Intn(n)
}
