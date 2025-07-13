package client

import "testing"

func TestContainsEarth(t *testing.T) {
	cases := []struct {
		origin string
		want   bool
	}{
		{"Earth", true},
		{"Earth (C-137)", true},
		{"Mars", false},
		{"", false},
	}

	for _, tt := range cases {
		got := containsEarth(tt.origin)
		if got != tt.want {
			t.Errorf("containsEarth(%q)=%v, want %v", tt.origin, got, tt.want)
		}
	}
}
