// Package configjson drills encoding/json (struct tags, marshal/unmarshal,
// custom marshaling) and the sort package.
package configjson

// User models a JSON object. Fill in the STRUCT TAGS so it (un)marshals as:
//
//	{"name":"Ada","age":36,"email":"ada@x.io","is_admin":true}
//
// Requirements:
//   - Name  -> "name"
//   - Age   -> "age"
//   - Email -> "email", but OMITTED when empty (omitempty)
//   - Admin -> "is_admin"
type User struct {
	Name  string // TODO: add `json:"name"`
	Age   int    // TODO: add `json:"age"`
	Email string // TODO: add `json:"email,omitempty"`
	Admin bool   // TODO: add `json:"is_admin"`
}

// ParseUser decodes JSON bytes into a User. Return the error from json.Unmarshal.
func ParseUser(data []byte) (User, error) {
	panic("TODO: implement ParseUser")
}

// ToJSON encodes u as JSON bytes. Return the error from json.Marshal.
func ToJSON(u User) ([]byte, error) {
	panic("TODO: implement ToJSON")
}

// SortByAge sorts users in place, ascending by Age, breaking ties by Name
// (ascending). Use sort.Slice.
func SortByAge(users []User) {
	panic("TODO: implement SortByAge")
}

// HexColor is a 24-bit RGB color. Implement json.Marshaler so it encodes as a
// JSON string like "#ff8800" (lower-case, always 6 hex digits).
type HexColor uint32

// MarshalJSON makes HexColor a json.Marshaler.
// Hint: the returned bytes must include the surrounding quotes, e.g. []byte(`"#ff8800"`).
func (c HexColor) MarshalJSON() ([]byte, error) {
	panic("TODO: implement HexColor.MarshalJSON")
}
