// Package strutil drills UTF-8, runes vs bytes, and the strings/unicode packages.
package strutil

import (
	"strings"
	"unicode"
)

// ReverseString reverses s by rune (Unicode code point), NOT by byte, so that
// multi-byte characters survive intact.
// Example: ReverseString("Hello, 世界") == "界世 ,olleH"
func ReverseString(s string) string {
	// panic("TODO: implement ReverseString")
	r := []rune(s)
	var ans []rune
	for i := len(r) - 1; i >= 0; i-- {
		ans = append(ans, r[i])
	}

	return string(ans)
}

// IsPalindrome reports whether s reads the same forwards and backwards,
// ignoring case and any character that is not a letter or digit.
// Example: "A man, a plan, a canal: Panama" is a palindrome.
func IsPalindrome(s string) bool {
	// panic("TODO: implement IsPalindrome")
	var clean []rune

	for _, r := range s {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			clean = append(clean, unicode.ToLower(r))
		}
	}

	for i, j := 0, len(clean)-1; i < j; i, j = i+1, j-1 {
		if clean[i] != clean[j] {
			return false
		}
	}
	return true
}

// CountVowels returns the number of ASCII vowels (a, e, i, o, u,
// case-insensitive) in s.
func CountVowels(s string) int {
	// panic("TODO: implement CountVowels")
	var r = []rune(s)
	var count int
	for _, v := range r {
		switch unicode.ToLower(v) {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// TitleCase upper-cases the first rune of each whitespace-separated word and
// lower-cases the rest. Build the result with strings.Builder.
// Example: TitleCase("hello WORLD go") == "Hello World Go"
func TitleCase(s string) string {
	// panic("TODO: implement TitleCase")
	var b strings.Builder

	for i, word := range strings.Fields(s) {
		if i > 0 {
			b.WriteRune(' ')
		}
		for j, r := range word {
			if j == 0 {
				b.WriteRune(unicode.ToUpper(r))
			} else {
				b.WriteRune(unicode.ToLower(r))
			}
		}

	}

	return b.String()

}
