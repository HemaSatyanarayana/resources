// Package strutil drills UTF-8, runes vs bytes, and the strings/unicode packages.
package strutil

// ReverseString reverses s by rune (Unicode code point), NOT by byte, so that
// multi-byte characters survive intact.
// Example: ReverseString("Hello, 世界") == "界世 ,olleH"
func ReverseString(s string) string {
	panic("TODO: implement ReverseString")
}

// IsPalindrome reports whether s reads the same forwards and backwards,
// ignoring case and any character that is not a letter or digit.
// Example: "A man, a plan, a canal: Panama" is a palindrome.
func IsPalindrome(s string) bool {
	panic("TODO: implement IsPalindrome")
}

// CountVowels returns the number of ASCII vowels (a, e, i, o, u,
// case-insensitive) in s.
func CountVowels(s string) int {
	panic("TODO: implement CountVowels")
}

// TitleCase upper-cases the first rune of each whitespace-separated word and
// lower-cases the rest. Build the result with strings.Builder.
// Example: TitleCase("hello WORLD go") == "Hello World Go"
func TitleCase(s string) string {
	panic("TODO: implement TitleCase")
}
