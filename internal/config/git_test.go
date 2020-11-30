package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gitlabTruncateSlug(t *testing.T) {
	type testcase struct {
		input    string
		expected string
	}
	testcases := []testcase{
		{
			input:    "abc-xyz",
			expected: "abc-xyz",
		},
		{
			input:    "core-features-app-operation-support-center-support-center-upload-api",
			expected: "core-features-app-operation-support-center-support-center-uploa",
		},
	}
	for _, tc := range testcases {
		res := gitlabTruncateSlug(tc.input)
		assert.Equal(t, tc.expected, res)
	}
}
