package main

import (
	"fmt"
	"reflect"

	"github.com/acgreek/gostructtags/encoding/json"
)

type Nested struct {
	N             string `json:"n"`
	Num           int
	UnsignedNum   uint
	Num32         int32
	UnsignedNum32 uint32
	Num64         int64
	UnsignedNum64 uint64
	Float32       float32
	Float64       float64
	Bool          bool
	link          *Nested
	ExportLink    *Nested

	PtrN             *string `json:"n"`
	PtrNum           *int
	PtrUnsignedNum   *uint
	PtrNum32         *int32
	PtrUnsignedNum32 *uint32
	PtrNum64         *int64 `json:"ptr_num64",omitempty`
	PtrUnsignedNum64 *uint64
	PtrFloat32       *float32
	PtrFloat64       *float64
	PtrBool          *bool
}

type TestStruct struct {
	Name string `json:"name",omitempty`
	N    Nested `json:"n"`
}

func main() {
	foo := TestStruct{
		Name: "foo",
		N: Nested{
			N: "bar",
		},
	}
	t := reflect.TypeOf(foo)
	v := reflect.ValueOf(foo)
	displayStruct(nil, 0, t, v)
	f, _ := t.FieldByName("Name")
	fmt.Printf("foo struct: %v is type %v %s\n", foo, t, f.Tag)
	output, err := json.Encode(foo)
	fmt.Printf("json=%s err=%s\n", output, err)
}

func displayStruct(stack []int, depth int, t reflect.Type, v reflect.Value) {
	for idx, field := range reflect.VisibleFields(t) {
		vs := v.FieldByIndex(append(stack, idx))
		fmt.Printf("%d idx=%d, kind=%s\n", depth, idx, vs.Kind())
		switch vs.Kind() {
		case reflect.Struct:
			t := vs.Type()
			displayStruct(append(stack, idx), depth+1, t, v)
		case reflect.String:
			fmt.Printf("%d field %s, %v %s value=%s\n", depth, field.Name, field.Type, field.Tag, vs.String())
		}
	}
}
