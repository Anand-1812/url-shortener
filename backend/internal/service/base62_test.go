package service

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0"},
		{61, "Z"},
		{62, "10"},
		{10000001, "FXsl"},
	}

	for _, test := range tests {
		if got := Encode(test.input); got != test.want {
			t.Errorf("Encode(%d) = %q, want %q", test.input, got, test.want)
		}
	}
}
