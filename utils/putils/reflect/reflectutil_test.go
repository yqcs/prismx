package reflectutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"prismx_cli/utils/putils/reflect/tests"
)

type TestStruct struct {
	FirstOption    string
	SecondOption   int
	privateOption3 string
}

func TestToMap(t *testing.T) {
	testStruct := TestStruct{
		FirstOption:    "test",
		SecondOption:   10,
		privateOption3: "ignored",
	}
	// testing normal fields
	tomap, err := ToMap(testStruct, nil, false)
	require.Nilf(t, err, "error while parsing: %s", err)
	m := map[string]interface{}{"first_option": "test", "second_option": 10}
	require.EqualValues(t, m, tomap, "objects are not equal")

	// testing with non exported ones
	tomap, err = ToMap(testStruct, nil, true)
	require.Nilf(t, err, "error while parsing: %s", err)
	m = map[string]interface{}{"first_option": "test", "second_option": 10, "private_option3": "ignored"}
	require.EqualValues(t, m, tomap, "objects are not equal")

	// testing with custom stringify function
	fu := func(s string) string {
		return strings.ToLower(s)
	}
	tomap, err = ToMap(testStruct, fu, false)
	require.Nilf(t, err, "error while parsing: %s", err)
	m = map[string]interface{}{"firstoption": "test", "secondoption": 10}
	require.EqualValues(t, m, tomap, "objects are not equal")
}

func TestUnexportedField(t *testing.T) {
	// create a pointer instance to a struct with an "unexported" field
	testStruct := &tests.Test{}
	SetUnexportedField(testStruct, "unexported", "test")
	value := GetUnexportedField(testStruct, "unexported")
	require.Equal(t, value, "test")
}
