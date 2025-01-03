package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPruneEmptyStrings(t *testing.T) {
	test := []string{"a", "", "", "b"}
	// converts back
	res := PruneEmptyStrings(test)
	require.Equal(t, []string{"a", "b"}, res, "strings not pruned correctly")
}

func TestPruneEqual(t *testing.T) {
	testStr := []string{"a", "", "", "b"}
	// converts back
	resStr := PruneEqual(testStr, "b")
	require.Equal(t, []string{"a", "", ""}, resStr, "strings not pruned correctly")

	testInt := []int{1, 2, 3, 4}
	// converts back
	resInt := PruneEqual(testInt, 2)
	require.Equal(t, []int{1, 3, 4}, resInt, "ints not pruned correctly")
}

func TestDedupe(t *testing.T) {
	testStr := []string{"a", "a", "b", "b"}
	// converts back
	resStr := Dedupe(testStr)
	require.Equal(t, []string{"a", "b"}, resStr, "strings not deduped correctly")

	testInt := []int{1, 1, 2, 2}
	// converts back
	res := Dedupe(testInt)
	require.Equal(t, []int{1, 2}, res, "ints not deduped correctly")
}

func TestPickRandom(t *testing.T) {
	testStr := []string{"a", "b"}
	// converts back
	resStr := PickRandom(testStr)
	require.Contains(t, testStr, resStr, "element was not picked correctly")

	testInt := []int{1, 2}
	// converts back
	resInt := PickRandom(testInt)
	require.Contains(t, testInt, resInt, "element was not picked correctly")
}

func TestContains(t *testing.T) {
	testSliceStr := []string{"a", "b"}
	testElemStr := "a"
	// converts back
	resStr := Contains(testSliceStr, testElemStr)
	require.True(t, resStr, "unexptected result")

	testSliceInt := []int{1, 2}
	testElemInt := 1
	// converts back
	resInt := Contains(testSliceInt, testElemInt)
	require.True(t, resInt, "unexptected result")
}

func TestContainsItems(t *testing.T) {
	test1Str := []string{"a", "b", "c"}
	test2Str := []string{"a", "c"}
	// converts back
	resStr := ContainsItems(test1Str, test2Str)
	require.True(t, resStr, "unexptected result")

	test1Int := []int{1, 2, 3}
	test2Int := []int{1, 3}
	// converts back
	resInt := ContainsItems(test1Int, test2Int)
	require.True(t, resInt, "unexptected result")
}

func TestToInt(t *testing.T) {
	test1 := []string{"1", "2"}
	test2 := []int{1, 2}
	// converts back
	res, err := ToInt(test1)
	require.Nil(t, err)
	require.Equal(t, test2, res, "unexptected result")
}

func TestEqual(t *testing.T) {
	test1 := []string{"1", "2"}
	require.True(t, Equal(test1, test1), "unexptected result")
	require.False(t, Equal(test1, []string{"2", "1"}), "unexptected result")
}

func TestIsEmpty(t *testing.T) {
	require.True(t, IsEmpty([]string{}))
	require.False(t, IsEmpty([]string{"a"}))
}

func TestElementsMatch(t *testing.T) {
	require.True(t, ElementsMatch([]string{}, []string{}))
	require.True(t, ElementsMatch([]int{1}, []int{1}))
	require.True(t, ElementsMatch([]int{1, 2}, []int{2, 1}))
	require.False(t, ElementsMatch([]int{1}, []int{2}))
}

func TestDiff(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{3, 4, 5}
	extraS1, extraS2 := Diff(s1, s2)
	require.ElementsMatch(t, extraS1, []int{1, 2})
	require.ElementsMatch(t, extraS2, []int{4, 5})
}

func TestMerge(t *testing.T) {
	tests := []struct {
		input    [][]int
		expected []int
	}{
		{[][]int{{1, 2, 3}, {3, 4, 5}, {5, 6, 7}}, []int{1, 2, 3, 4, 5, 6, 7}},
		{[][]int{{1, 1, 2}, {2, 3, 3}, {3, 4, 5}}, []int{1, 2, 3, 4, 5}},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, []int{1, 2, 3, 4, 5, 6}},
	}

	for _, test := range tests {
		output := Merge(test.input...)
		require.ElementsMatch(t, test.expected, output)
	}
}

func TestMergeItems(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3, 3, 4, 5, 5, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
		{[]int{1, 1, 2, 2, 3, 3}, []int{1, 2, 3}},
		{[]int{1, 2, 3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
	}

	for _, test := range tests {
		// merge single basic types (int, string, etc)
		output := MergeItems(test.input...)
		require.ElementsMatch(t, test.expected, output)
	}

}

func TestFirstNonZeroInt(t *testing.T) {
	testCases := []struct {
		Input          []int
		ExpectedOutput interface{}
		ExpectedFound  bool
	}{
		{
			Input:          []int{0, 0, 3, 5, 10},
			ExpectedOutput: 3,
			ExpectedFound:  true,
		},
		{
			Input:          []int{},
			ExpectedOutput: 0,
			ExpectedFound:  false,
		},
	}

	for _, tc := range testCases {
		output, found := FirstNonZero(tc.Input)
		require.Equal(t, tc.ExpectedOutput, output)
		require.Equal(t, tc.ExpectedFound, found)
	}
}

func TestFirstNonZeroString(t *testing.T) {
	testCases := []struct {
		Input          []string
		ExpectedOutput interface{}
		ExpectedFound  bool
	}{
		{
			Input:          []string{"", "foo", "test"},
			ExpectedOutput: "foo",
			ExpectedFound:  true,
		},
		{
			Input:          []string{},
			ExpectedOutput: "",
			ExpectedFound:  false,
		},
	}

	for _, tc := range testCases {
		output, found := FirstNonZero(tc.Input)
		require.Equal(t, tc.ExpectedOutput, output)
		require.Equal(t, tc.ExpectedFound, found)
	}
}

func TestFirstNonZeroFloat(t *testing.T) {
	testCases := []struct {
		Input          []float64
		ExpectedOutput interface{}
		ExpectedFound  bool
	}{
		{
			Input:          []float64{0.0, 0.0, 0.0, 1.2, 3.4},
			ExpectedOutput: 1.2,
			ExpectedFound:  true,
		},
		{
			Input:          []float64{},
			ExpectedOutput: 0.0,
			ExpectedFound:  false,
		},
	}

	for _, tc := range testCases {
		output, found := FirstNonZero(tc.Input)
		require.Equal(t, tc.ExpectedOutput, output)
		require.Equal(t, tc.ExpectedFound, found)
	}
}

func TestFirstNonZeroBool(t *testing.T) {
	testCases := []struct {
		Input          []bool
		ExpectedOutput interface{}
		ExpectedFound  bool
	}{
		{
			Input:          []bool{false, false, false},
			ExpectedOutput: false,
			ExpectedFound:  false,
		},
		{
			Input:          []bool{},
			ExpectedOutput: false,
			ExpectedFound:  false,
		},
	}

	for _, tc := range testCases {
		output, found := FirstNonZero(tc.Input)
		require.Equal(t, tc.ExpectedOutput, output)
		require.Equal(t, tc.ExpectedFound, found)
	}
}

func TestClone(t *testing.T) {
	intSlice := []int{1, 2, 3}
	require.Equal(t, intSlice, Clone(intSlice))

	stringSlice := []string{"a", "b", "c"}
	require.Equal(t, stringSlice, Clone(stringSlice))

	bytesSlice := []byte{1, 2, 3}
	require.Equal(t, bytesSlice, Clone(bytesSlice))
}

func TestVisitSequential(t *testing.T) {
	intSlice := []int{1, 2, 3}
	var res []int
	visit := func(index int, item int) error {
		res = append(res, item)
		return nil
	}
	VisitSequential(intSlice, visit)
	require.Equal(t, intSlice, res)
}

func TestVisitRandom(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var timesDifferent int
	for i := 0; i < 100; i++ {
		var res []int
		visit := func(index int, item int) error {
			res = append(res, item)
			return nil
		}
		VisitRandom(intSlice, visit)
		if !Equal(intSlice, res) {
			timesDifferent++
		}
		require.ElementsMatch(t, intSlice, res)
	}
	require.True(t, timesDifferent > 0)
}

func TestVisitRandomZero(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var timesDifferent int
	for i := 0; i < 100; i++ {
		var res []int
		visit := func(index int, item int) error {
			res = append(res, item)
			return nil
		}
		VisitRandomZero(intSlice, visit)
		if !Equal(intSlice, res) {
			timesDifferent++
		}
		require.ElementsMatch(t, intSlice, res)
	}
	require.True(t, timesDifferent > 0)
}
