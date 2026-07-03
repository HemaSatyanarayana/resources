package strutil

import "testing"

func TestReverseString(t *testing.T) {
	cases := map[string]string{
		"Hello, 世界": "界世 ,olleH",
		"abc":       "cba",
		"":          "",
		"a":         "a",
		"résumé":    "émusér",
	}
	for in, want := range cases {
		if got := ReverseString(in); got != want {
			t.Errorf("ReverseString(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	cases := map[string]bool{
		"A man, a plan, a canal: Panama": true,
		"racecar":                        true,
		"No lemon, no melon":             true,
		"hello":                          false,
		"":                               true,
		"Was it a car or a cat I saw?":   true,
	}
	for in, want := range cases {
		if got := IsPalindrome(in); got != want {
			t.Errorf("IsPalindrome(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestCountVowels(t *testing.T) {
	cases := map[string]int{
		"hello":      2,
		"AEIOU":      5,
		"rhythm":     0,
		"Go is fun!": 3,
	}
	for in, want := range cases {
		if got := CountVowels(in); got != want {
			t.Errorf("CountVowels(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestTitleCase(t *testing.T) {
	cases := map[string]string{
		"hello WORLD go": "Hello World Go",
		"the QUICK brown": "The Quick Brown",
		"":                "",
		"a":               "A",
	}
	for in, want := range cases {
		if got := TitleCase(in); got != want {
			t.Errorf("TitleCase(%q) = %q, want %q", in, got, want)
		}
	}
}
