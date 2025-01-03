package stringsutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type betweentest struct {
	After     string
	Before    string
	Result    string
	WantError bool
}

func TestBetween(t *testing.T) {
	tests := map[string]betweentest{
		"a b c":                            {After: "a", Before: "c", Result: " b ", WantError: false},
		"this is a test":                   {After: "this", Before: "test", Result: " is a ", WantError: false},
		"this is a test bbb test":          {After: "test", Before: "test", Result: " bbb ", WantError: false},
		"this is a test with before error": {After: "test", Before: "testt", Result: "this is a test with before error", WantError: true},
		"this is a test with after error":  {After: "testt", Before: "test", Result: "this is a test with after error", WantError: true},
	}
	for str, test := range tests {
		res, err := Between(str, test.After, test.Before)
		if test.WantError {
			require.Error(t, err)
		}
		require.Equalf(t, test.Result, res, "test: %s after: %s before: %s result: %s", str, test.After, test.Before, res)
	}
}

func TestBefore(t *testing.T) {
	tests := map[string]betweentest{
		"a b c":                        {Before: "c", Result: "a b ", WantError: false},
		"this is a test":               {Before: "test", Result: "this is a ", WantError: false},
		"this is a test with second t": {Before: "testt", Result: "this is a test with second t", WantError: true},
	}
	for str, test := range tests {
		res, err := Before(str, test.Before)
		if test.WantError {
			require.Error(t, err)
		}
		require.Equalf(t, test.Result, res, "test: %s before: %s result: %s", str, test.Before, res)
	}
}

func TestAfter(t *testing.T) {
	tests := map[string]betweentest{
		"a b c":                        {After: "a", Result: " b c", WantError: false},
		"this is a test":               {After: "this", Result: " is a test", WantError: false},
		"this is a test with second t": {After: "testt", Result: "this is a test with second t", WantError: true},
	}
	for str, test := range tests {
		res, err := After(str, test.After)
		if test.WantError {
			require.Error(t, err)
		}
		require.Equalf(t, test.Result, res, "test: %s after: %s result: %s", str, test.After, res)
	}
}

type prefixsuffixtest struct {
	Prefixes []string
	Suffixes []string
	Result   interface{}
}

func TestHasPrefixAny(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"a b c":     {Prefixes: []string{"a"}, Result: true},
		"a b c d":   {Prefixes: []string{"a b", "a"}, Result: true},
		"a b c d e": {Prefixes: []string{"b", "o", "a"}, Result: true},
		"test test": {Prefixes: []string{"a", "b"}, Result: false},
	}
	for str, test := range tests {
		res := HasPrefixAny(str, test.Prefixes...)
		require.Equalf(t, test.Result, res, "test: %s prefixes: %+v result: %s", str, test.Prefixes, res)
	}
}

func TestHasPrefixAnyI(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"A b c":     {Prefixes: []string{"a"}, Result: true},
		"A B c d":   {Prefixes: []string{"a b", "a"}, Result: true},
		"a b c d e": {Prefixes: []string{"b", "o", "A"}, Result: true},
		"test test": {Prefixes: []string{"a", "b"}, Result: false},
	}
	for str, test := range tests {
		res := HasPrefixAnyI(str, test.Prefixes...)
		require.Equalf(t, test.Result, res, "test: %s prefixes: %+v result: %s", str, test.Prefixes, res)
	}
}

func TestHasSuffixAny(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"a b c":     {Suffixes: []string{"c"}, Result: true},
		"a b c d":   {Suffixes: []string{"c d", "a"}, Result: true},
		"a b c d e": {Suffixes: []string{"c", "d", "e"}, Result: true},
		"test test": {Suffixes: []string{"a", "b"}, Result: false},
	}
	for str, test := range tests {
		res := HasSuffixAny(str, test.Suffixes...)
		require.Equalf(t, test.Result, res, "test: %s suffixes: %+v result: %s", str, test.Suffixes, res)
	}
}

func TestTrimPrefixAny(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"a b c":     {Prefixes: []string{"a"}, Result: " b c"},
		"a b c d":   {Prefixes: []string{"a b", "a"}, Result: " c d"},
		"a b c d e": {Prefixes: []string{"b", "o", "a"}, Result: " b c d e"},
		"test test": {Prefixes: []string{"a", "b"}, Result: "test test"},
	}
	for str, test := range tests {
		res := TrimPrefixAny(str, test.Prefixes...)
		require.Equalf(t, test.Result, res, "test: %s prefixes: %+v result: %s", str, test.Prefixes, res)
	}
}

func TestTrimSuffixAny(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"a b c":     {Suffixes: []string{"c"}, Result: "a b "},
		"a b c d":   {Suffixes: []string{"c d", "a"}, Result: "a b "},
		"a b c d e": {Suffixes: []string{"e"}, Result: "a b c d "},
		"test test": {Suffixes: []string{"a", "b"}, Result: "test test"},
	}
	for str, test := range tests {
		res := TrimSuffixAny(str, test.Suffixes...)
		require.Equalf(t, test.Result, res, "test: %s suffixes: %+v result: %s", str, test.Suffixes, res)
	}
}

type jointest struct {
	Items     []interface{}
	Separator string
	Result    string
}

func TestJoin(t *testing.T) {
	tests := []jointest{
		{Items: []interface{}{"a"}, Separator: "", Result: "a"},
		{Items: []interface{}{"a", "b"}, Separator: ",", Result: "a,b"},
		{Items: []interface{}{"a", "b", 1}, Separator: ",", Result: "a,b,1"},
		{Items: []interface{}{2, "b", 1}, Separator: "", Result: "2b1"},
	}
	for _, test := range tests {
		res := Join(test.Items, test.Separator)
		require.Equalf(t, test.Result, res, "test: %+v", test)
	}
}

func TestHasPrefixI(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"a b c":   {Prefixes: []string{"a"}, Result: true},
		"A b c d": {Prefixes: []string{"a"}, Result: true},
		"Ab c d":  {Prefixes: []string{"b"}, Result: false},
	}
	for str, test := range tests {
		res := HasPrefixI(str, test.Prefixes[0])
		require.Equalf(t, test.Result, res, "test: %s prefixes: %+v result: %s", str, test.Prefixes, res)
	}
}

func TestHasSuffixI(t *testing.T) {
	tests := map[string]prefixsuffixtest{
		"a b c":  {Prefixes: []string{"c"}, Result: true},
		"A b C":  {Prefixes: []string{"c"}, Result: true},
		"Ab c d": {Prefixes: []string{"c"}, Result: false},
	}
	for str, test := range tests {
		res := HasSuffixI(str, test.Prefixes[0])
		require.Equalf(t, test.Result, res, "test: %s suffixes: %+v result: %s", str, test.Suffixes, res)
	}
}

func TestReverse(t *testing.T) {
	tests := map[string]string{
		"abc":          "cba",
		"A b C":        "C b A",
		"Ab c d":       "d c bA",
		"hello world!": "!dlrow olleh",
		"!@#$%^&*()":   ")(*&^%$#@!",
		"明日は晴天り":       "り天晴は日明",
	}
	for str, expRes := range tests {
		res := Reverse(str)
		require.Equalf(t, expRes, res, "test: %s expected: %+v result: %s", str, expRes, res)
	}
}

type containstest struct {
	Items  []string
	Result bool
}

type replacealltest struct {
	Old    string
	New    string
	Result string
}

func TestContainsAny(t *testing.T) {
	tests := map[string]containstest{
		"abc":   {Items: []string{"a", "b"}, Result: true},
		"abcd":  {Items: []string{"x", "b"}, Result: true},
		"A b C": {Items: []string{"x"}, Result: false},
	}
	for str, test := range tests {
		res := ContainsAny(str, test.Items...)
		require.Equalf(t, test.Result, res, "test: %+v", res)
	}
}

func TestContainsAnyI(t *testing.T) {
	tests := map[string]containstest{
		"abc":    {Items: []string{"A", "b"}, Result: true},
		"abcd":   {Items: []string{"X", "b"}, Result: true},
		"A b C":  {Items: []string{"X"}, Result: false},
		"aaa":    {Items: []string{"A"}, Result: true},
		"Hello!": {Items: []string{"hELLO", "world"}, Result: true},
	}
	for str, test := range tests {
		res := ContainsAnyI(str, test.Items...)
		require.Equalf(t, test.Result, res, "test: %+v", res)
	}
}

func TestEqualFoldAny(t *testing.T) {
	tests := map[string]containstest{
		"abc":   {Items: []string{"a", "Abc"}, Result: true},
		"abcd":  {Items: []string{"x", "ABcD"}, Result: true},
		"A b C": {Items: []string{"x"}, Result: false},
		"hello": {Items: []string{"llo", "heLLo"}, Result: true},
		"world": {Items: []string{"Hello", "WoRld"}, Result: true},
	}
	for str, test := range tests {
		res := EqualFoldAny(str, test.Items...)
		require.Equalf(t, test.Result, res, "test: %+v", res)
	}
}

type attest struct {
	After  int
	Search string
	Result interface{}
}

func TestIndexAt(t *testing.T) {
	tests := map[string]attest{
		"a a b":          {After: 1, Search: "a", Result: 2},
		"test":           {After: 1, Search: "t", Result: 3},
		"test test":      {After: 4, Search: "test", Result: 5},
		"test test test": {After: 0, Search: "test", Result: 0},
	}
	for str, test := range tests {
		res := IndexAt(str, test.Search, test.After)
		require.Equalf(t, test.Result, res, "test: %s after: %d search: %s result: %d", str, test.After, test.Search, res)
	}
}

type splitanytest struct {
	Splitset []string
	Result   interface{}
}

func TestSplitAny(t *testing.T) {
	tests := map[string]splitanytest{
		"a a b":        {Splitset: []string{" "}, Result: []string{"a", "a", "b"}},
		"test1test2 3": {Splitset: []string{"1", "2", " "}, Result: []string{"test", "test", "3"}},
	}
	for str, test := range tests {
		res := SplitAny(str, test.Splitset...)
		require.Equalf(t, test.Result, res, "test: %s splitset: %v result: %v got: %v", str, test.Splitset, test.Result, res)
	}
}

func TestSlideWithLength(t *testing.T) {
	var res []string
	c := SlideWithLength("test123", 4)
	require.NotNil(t, c)
	for cc := range c {
		res = append(res, cc)
	}
	require.Equal(t, []string{"test", "est1", "st12", "t123", "123"}, res)
}

func TestReplaceAll(t *testing.T) {
	tests := map[string]replacealltest{
		"hello":      {Old: "l", New: "k", Result: "hekko"},
		"abcd":       {Old: "a", New: "b", Result: "bbcd"},
		"A b C":      {Old: " ", New: "", Result: "AbC"},
		"!@#$%^&*()": {Old: "@", New: "_", Result: "!_#$%^&*()"},
	}
	for str, test := range tests {
		res := ReplaceAll(str, test.New, test.Old)
		require.Equalf(t, test.Result, res, "test: %s, new: %s, old: %s, expected: %s, got: %s", str, test.New, test.Old, test.Result, res)
	}
}

func TestLongestRepeatingSequence(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"abcdefg", ""},
		{"abcabcabc", "abc"},
		{"abcdefabcdef", "abcdef"},
		{"abcdefgabcdefg", "abcdefg"},
		{"abcabcdefdef", "abc"},
	}

	for _, test := range tests {
		result := LongestRepeatingSequence(test.s)
		require.Equalf(t, test.expected, result.Sequence, "test: %s, expected %q, got: %s", test.s, test.expected, result.Sequence)
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"abcdefg", true},
		{"abcabcabc", true},
		{"abcdefabcdef", true},
		{"abcdefgabcdefg", true},
		{"abcabcdefdef", true},
		{"\x03", false},
	}

	for _, test := range tests {
		result := IsPrintable(test.s)
		require.Equalf(t, test.expected, result, "test: %s, expected %q, got: %s", test.s, test.expected, result)
	}
}

func TestIsCTRLC(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"aaa", false},
		{"\x03", true},
	}

	for _, test := range tests {
		result := IsCTRLC(test.s)
		require.Equalf(t, test.expected, result, "test: %s, expected %q, got: %s", test.s, test.expected, result)
	}
}

type truncateTest struct {
	test    string
	maxSize int
	result  string
}

func TestTruncate(t *testing.T) {
	tests := []truncateTest{
		{test: "abcd", maxSize: -1, result: "abcd"},
		{test: "abcd", maxSize: 0, result: ""},
		{test: "abcde", maxSize: 3, result: "abc"},
		{test: "abcdef", maxSize: 8, result: "abcdef"},
		{test: "abcdefg", maxSize: 6, result: "abcdef"},
		{test: "aaaa", maxSize: 20, result: "aaaa"},
		{test: "aaaa", maxSize: 4, result: "aaaa"},
	}

	for _, test := range tests {
		res := Truncate(test.test, test.maxSize)
		require.Equalf(t, test.result, res, "test:%s maxsize: %d result: %s", test.test, test.maxSize, res)
	}
}

func TestIndexAny(t *testing.T) {
	tests := []struct {
		s           string
		seps        []string
		expectedIdx int
		expectedSep string
	}{
		{"abcdefg", []string{"a", "b"}, 0, "a"},
		{"abcdefg", []string{"z", "b"}, 1, "b"},
		{"abcdefg", []string{"z", "zz"}, -1, ""},
	}
	for _, test := range tests {
		idx, sep := IndexAny(test.s, test.seps...)
		require.Equal(t, test.expectedIdx, idx)
		require.Equal(t, test.expectedSep, sep)
	}
}
