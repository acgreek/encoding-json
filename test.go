package main

import (
	"fmt"
	"reflect"
)

type Nested struct {
	N string
}

type TestStruct struct {
	Name string `require not empty`
	N    Nested
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
