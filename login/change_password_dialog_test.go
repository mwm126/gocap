package login

import (
	"testing"
)

func TestValidateNewPassword(t *testing.T) {
	testCases := []struct {
		label  string
		old    string
		new    string
		new2   string
		result string
	}{
		{"mismatch", "abc", "abc", "abcd", "Passwords do not match"},
		{"same", "abc", "abc", "abc", "Password is the same as the previous password"},
		{"short", "abcd", "abc", "abc", "Password must have length >=12 characters"},
		{"needlower", "", "ABCDEFGHIJKL", "ABCDEFGHIJKL", "Password must contain a lowercase letter"},
		{"needupper", "", "abcdefghijkl", "abcdefghijkl", "Password must contain an uppercase letter"},
		{"needdigit", "", "Abcdefghijkl", "Abcdefghijkl", "Password must contain a digit"},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := password_passes(tc.old, tc.new, tc.new2).Error()
			want := tc.result
			if got != want {
				t.Errorf("Got %s; want %s", got, want)
			}
		})
	}
}
