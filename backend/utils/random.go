package utils

import (
	"math/rand/v2"
	"strconv"
)

func CreatCode() string {
	code := rand.IntN(900000) + 100000 //100000-999999
	return strconv.Itoa(code)
}
