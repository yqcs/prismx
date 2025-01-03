package ptr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSafe(t *testing.T) {
	type args[T any] struct {
		v *T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "struct=>int - NilPointer",
			args: args[int]{v: nil},
			want: 0,
		},
		{
			name: "struct=>int - NonNilPointer",
			args: args[int]{v: new(int)},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Safe(tt.args.v)
			require.Equal(t, tt.want, got, "Safe() = %v, want %v", got, tt.want)
		})
	}
}
