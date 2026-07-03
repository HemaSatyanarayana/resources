# 03 — Strings & Runes

Go strings are **immutable UTF-8 byte sequences**. Indexing `s[i]` gives you a *byte*, not a character. The moment you deal with non-ASCII text you must think in **runes** (`rune` = `int32` = one Unicode code point).

## Concepts

- **`len(s)`** is the byte count, not the character count.
- **`for i, r := range s`** iterates *runes*, with `i` being the byte offset and `r` the `rune`.
- **`[]rune(s)`** converts to a slice of code points — the clean way to reverse or index "by character".
- **`strings.Builder`** builds strings without the O(n²) cost of `+=` in a loop.
- **`unicode`** package: `unicode.ToLower`, `unicode.IsLetter`, `unicode.IsDigit`.

## Your task

Implement the four functions in [`strutil.go`](strutil.go):

| Function | Skill |
|----------|-------|
| `ReverseString` | `[]rune` conversion, two-pointer swap |
| `IsPalindrome` | Filter runes, case-fold, compare |
| `CountVowels` | Range over runes |
| `TitleCase` | `strings.Builder`, `strings.Fields`, `unicode.ToUpper/ToLower` |

## Run

```bash
go test -v ./exercises/03-strings-runes/
```

## Hints

- Reversing `s[i]` byte-by-byte will corrupt any multi-byte character — convert to `[]rune` first.
- For the palindrome, extract the "clean" runes (letters/digits, lower-cased) into a `[]rune`, then compare from both ends.
- `strings.Builder`: declare `var b strings.Builder`, call `b.WriteRune(r)` / `b.WriteString(...)`, finish with `b.String()`.
- `strings.Title` is **deprecated** — that's exactly why you're implementing `TitleCase` yourself.

<details>
<summary>Reference solution</summary>

```go
package strutil

import (
	"strings"
	"unicode"
)

func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func IsPalindrome(s string) bool {
	var clean []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
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

func CountVowels(s string) int {
	n := 0
	for _, r := range s {
		switch unicode.ToLower(r) {
		case 'a', 'e', 'i', 'o', 'u':
			n++
		}
	}
	return n
}

func TitleCase(s string) string {
	var b strings.Builder
	for i, word := range strings.Fields(s) {
		if i > 0 {
			b.WriteByte(' ')
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
```

</details>
