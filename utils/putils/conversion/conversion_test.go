package conversion

import (
	"bytes"
	"testing"
)

func TestBytes(t *testing.T) {
	testCases := []struct {
		input    string
		expected []byte
	}{
		{"test", []byte("test")},
		{"", []byte("")},
	}

	for _, tc := range testCases {
		result := Bytes(tc.input)
		if !bytes.Equal(result, tc.expected) {
			t.Errorf("Expected %v, but got %v", tc.expected, result)
		}
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("test"), "test"},
		{[]byte(""), ""},
	}

	for _, tc := range testCases {
		result := String(tc.input)
		if result != tc.expected {
			t.Errorf("Expected %s, but got %s", tc.expected, result)
		}
	}
}
