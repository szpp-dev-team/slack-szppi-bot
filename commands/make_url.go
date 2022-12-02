package commands

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

func MakeUrl(x int, R string) string {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
	max := x
	var min int
	if R == "E" || R == "F" {
		min = 126
	} else if R == "G" || R == "H" || R == "EX" {
		min = 212
	} else {
		min = 0
	}
	var tmp int = rand.Intn(max-min) + min
	ran := fmt.Sprintf("%d", tmp)
	if tmp < 100 && tmp > 10 {
		ran = "0" + ran
	}
	if tmp < 10 {
		ran = "00" + ran
	}
	//ran := string(rand.Intn(max-min) + min)
	url := "https://atcoder.jp/contests/abc" + ran + "/tasks/abc" + ran + "_" + R
	//https://atcoder.jp/contests/abc258/tasks/abc258_b
	return url
}
