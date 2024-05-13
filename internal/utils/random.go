package utils

import "math/rand/v2"

// GetRandomBoundedInt returns a num from 0 - n (exclusive)
// if n = 0, returns 0
func GetRandomBoundedInt(r *rand.Rand, n int) int {
	if n == 0 {
		return 0
	}
	return r.IntN(n)
}
