package rand

import (
	"crypto/rand"
	"errors"
	"math/big"
	crand "math/rand"
)

// IntN returns a uniform random value in [0, max). It errors if max <= 0.
func IntN(max int) (int, error) {
	if max <= 0 {
		return 0, errors.New("max can't be <= 0")
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return crand.Intn(max), nil
	}
	return int(nBig.Int64()), nil
}
