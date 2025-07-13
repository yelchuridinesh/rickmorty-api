// internal/api/helpers_test.go
package api

import "testing"

func TestSanitizeSortBy(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"id", "id"},
		{"name", "name"},
		{"Id", "id"}, // case-sensitive, so falls back
		{"", "id"},   // empty → fallback
		{"created_at", "id"},
	}

	for _, tc := range cases {
		got := sanitizeSortBy(tc.in)
		if got != tc.want {
			t.Errorf("sanitizeSortBy(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}
