package display

import "testing"

func TestTruncateString(t *testing.T) {
	testCases := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"Hello, World!", 5, "He..."},
		{"Short", 10, "Short"},
		{"This is a long string", 10, "This is..."},
		{"", 5, ""},
	}

	for _, tc := range testCases {
		result := truncateString(tc.input, tc.maxLen)
		if result != tc.expected {
			t.Errorf("truncateString(%q, %d) = %q; want %q", tc.input, tc.maxLen, result, tc.expected)
		}
	}
}
