package json

import (
	"encoding/json"
	"testing"
)

func TestPtrOmitEmpty(t *testing.T) {
	type Foo struct {
		Do *int64 `json:"do,omitempty"`
	}
	d := Foo{}
	encoded, err := Encode(d)
	if err != nil {
		t.Errorf("Encode: expect no error but got %v", err)
	}
	expected, _ := json.Marshal(d)
	if string(expected) != string(encoded) {
		t.Errorf("Encode: expect %s but got %s", string(expected), string(encoded))
	}
}

func TestPtrOmitEmptyWithValue(t *testing.T) {
	type Foo struct {
		Do *int64 `json:"do,omitempty"`
	}
	v := int64(5)
	d := Foo{
		Do: &v,
	}
	encoded, err := Encode(d)
	if err != nil {
		t.Errorf("Encode: expect no error but got %v", err)
	}
	expected, _ := json.Marshal(d)
	if string(expected) != string(encoded) {
		t.Errorf("Encode: expect %s but got %s", string(expected), string(encoded))
	}
}

func TestBoolPtrStringify(t *testing.T) {
	type Foo struct {
		Do      *bool `json:",string"`
		Missing int   `json:"-"`
	}
	v := false
	d := Foo{
		Do:      &v,
		Missing: 4,
	}
	encoded, err := Encode(d)
	if err != nil {
		t.Errorf("Encode: expect no error but got %v", err)
	}
	expected, _ := json.Marshal(d)
	if string(expected) != string(encoded) {
		t.Errorf("Encode: expect %s but got %s", string(expected), string(encoded))
	}
}

func TestStructInStruct(t *testing.T) {
	type Foo struct {
		Word string
	}
	type Bar struct {
		F Foo
	}
	d := Bar{
		F: Foo{
			Word: "foo",
		},
	}
	encoded, err := Encode(d)
	if err != nil {
		t.Errorf("Encode: expect no error but got %v", err)
	}
	expected, _ := json.Marshal(d)
	if string(expected) != string(encoded) {
		t.Errorf("Encode: expect %s but got %s", string(expected), string(encoded))
	}
}

func TestStructPtr(t *testing.T) {
	type Foo struct {
		Word string
		Do   *Foo
	}
	d := Foo{
		Word: "foo",
		Do: &Foo{
			Word: "foobar",
		},
	}
	encoded, err := Encode(d)
	if err != nil {
		t.Errorf("Encode: expect no error but got %v", err)
	}
	expected, _ := json.Marshal(d)
	if string(expected) != string(encoded) {
		t.Errorf("Encode: expect %s but got %s", string(expected), string(encoded))
	}
}

func TestComplexNotSupported(t *testing.T) {
	type Foo struct {
		C complex128
	}
	d := Foo{
		C: complex128(3),
	}
	_, err := Encode(d)
	if err == nil {
		t.Error("Encode: expect an error because complex is not supported")
	}
	_, err = json.Marshal(d)
	if err == nil {
		t.Error("Encode: expect an error because complex is not supported")
	}
}

func TestEncode(t *testing.T) {
	type TestCase struct {
		name  string
		input interface{}
	}
	tests := []*TestCase{
		{
			name: "uintptr",
			input: struct {
				V uintptr
			}{V: 5},
		},
		{
			name: "complex",
			input: struct {
				V complex128
			}{V: 5},
		},
		{
			name: "[]string",
			input: struct {
				V []string
			}{V: []string{"foo", "bar"}},
		},
		{
			name: "[2]string",
			input: struct {
				V [2]string
			}{V: [2]string{"foo", "bar"}},
		},
		{
			name: "[2]*string",
			input: struct {
				V [2]*string
			}{V: [2]*string{strPtr("foo"), strPtr("bar")}},
		},
		{
			name: "chan",
			input: struct {
				V chan string
			}{V: make(chan string)},
		},
		{
			name: "map[string]int",
			input: struct {
				V map[string]int
			}{V: map[string]int{"foo": 1}},
		},
		{
			name: "map[int]string",
			input: struct {
				V map[int]string
			}{V: map[int]string{1: "foo"}},
		},
	}
	for _, tc := range tests {
		result, err := Encode(tc.input)
		expected, expErr := json.Marshal(tc.input)
		if expErr != nil && err == nil {
			t.Errorf("%s: expected error but did not get error. Exp error: %v", tc.name, expErr)
		} else {
			if expErr == nil && string(result) != string(expected) {
				t.Errorf("%s: expected \"%s\", but got \"%s", tc.name, expected, result)
			}
		}
	}
}

func strPtr(s string) *string {
	return &s
}
