package pin_utils

import "math/rand"

func Generate(low, hi int) int {
	return low + rand.Intn(hi-low)
}
