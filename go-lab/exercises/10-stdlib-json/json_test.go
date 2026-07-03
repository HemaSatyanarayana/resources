package configjson

import (
	"encoding/json"
	"testing"
)

func TestParseUser(t *testing.T) {
	data := []byte(`{"name":"Ada","age":36,"email":"ada@x.io","is_admin":true}`)
	u, err := ParseUser(data)
	if err != nil {
		t.Fatalf("ParseUser error: %v", err)
	}
	if u.Name != "Ada" || u.Age != 36 || u.Email != "ada@x.io" || !u.Admin {
		t.Errorf("ParseUser = %+v", u)
	}
}

func TestToJSON(t *testing.T) {
	u := User{Name: "Grace", Age: 45, Admin: false} // no email
	b, err := ToJSON(u)
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	// Re-decode into a generic map to assert on keys/values robustly.
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("output was not valid JSON: %v (%s)", err, b)
	}
	if m["name"] != "Grace" {
		t.Errorf("name = %v, want Grace", m["name"])
	}
	if _, ok := m["is_admin"]; !ok {
		t.Errorf("missing is_admin key in %s", b)
	}
	if _, ok := m["email"]; ok {
		t.Errorf("empty email should be omitted, got %s", b)
	}
}

func TestRoundTrip(t *testing.T) {
	in := User{Name: "Lin", Age: 29, Email: "lin@x.io", Admin: true}
	b, err := ToJSON(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ParseUser(b)
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Errorf("round trip: got %+v, want %+v", out, in)
	}
}

func TestSortByAge(t *testing.T) {
	users := []User{
		{Name: "Cara", Age: 30},
		{Name: "Ann", Age: 30},
		{Name: "Bob", Age: 22},
	}
	SortByAge(users)
	want := []string{"Bob", "Ann", "Cara"} // 22, then 30s tie-broken by name
	for i, w := range want {
		if users[i].Name != w {
			t.Errorf("position %d = %s, want %s (%+v)", i, users[i].Name, w, users)
		}
	}
}

func TestHexColorMarshal(t *testing.T) {
	b, err := json.Marshal(HexColor(0xff8800))
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `"#ff8800"` {
		t.Errorf("HexColor marshal = %s, want \"#ff8800\"", b)
	}
	// Leading zeros must be preserved.
	b, _ = json.Marshal(HexColor(0x0000ff))
	if string(b) != `"#0000ff"` {
		t.Errorf("HexColor marshal = %s, want \"#0000ff\"", b)
	}
}
