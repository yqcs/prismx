package rand

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntN(t *testing.T) {
	type testCase struct {
		input      int
		expectedOk bool
	}

	testCases := []testCase{
		{input: 10, expectedOk: true},
		{input: 0, expectedOk: false},
		{input: -10, expectedOk: false},
	}

	for _, tc := range testCases {
		i, err := IntN(tc.input)
		ok := i >= 0 && i <= tc.input && err == nil
		require.Equal(t, tc.expectedOk, ok)
	}
}
