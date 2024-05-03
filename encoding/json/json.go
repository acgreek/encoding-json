package json

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

func Encode(in interface{}) ([]byte, error) {
	buf := &strings.Builder{}
	t := reflect.TypeOf(in)
	err := encodeValue(t.Kind(), buf, nil, 0, 0, in, false)
	return []byte(buf.String()), err
}

func encodeStruct(buf io.Writer, stack []int, depth int, in interface{}) error {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)
	for idx, field := range reflect.VisibleFields(t) {
		vs := v.FieldByIndex([]int{idx})
		if err := encodeStructField(buf, vs, field, stack, idx, depth); err != nil {
			return fmt.Errorf("struct field: %w", err)
		}
	}
	return nil
}

func encodeStructField(buf io.Writer, vs reflect.Value, field reflect.StructField, stack []int, idx int, depth int) error {
	name := field.Name
	omitEmpty := false
	stringify := false
	tagValue, ok := field.Tag.Lookup("json")
	if ok {
		sections := strings.Split(tagValue, ",")
		if sections[0] == "-" {
			return nil
		}
		if sections[0] != "" {
			name = sections[0]
		}
		for _, value := range sections[1:] {
			if value == "omitempty" {
				omitEmpty = true
			}
			if value == "string" {
				stringify = true
			}
		}
	}
	if !vs.CanInterface() || (omitEmpty && IsZeroOfUnderlyingType(vs.Interface())) {
		return nil
	}
	if idx > 0 {
		fmt.Fprintf(buf, ",")
	}
	return encodeInterface(buf, name, vs.Kind(), vs.Interface(), stringify, stack, idx, depth)
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func encodeInterface(buf io.Writer, name string, kind reflect.Kind, vs interface{}, stringify bool,
	stack []int, idx int, depth int) error {
	fmt.Fprintf(buf, "\"%s\":", name)
	if stringify {
		tempBuf := &strings.Builder{}
		err := encodeValue(kind, tempBuf, stack, idx, depth, vs, stringify)
		if err != nil {
			return fmt.Errorf("field %s: %v", name, err)
		}
		buf.Write([]byte("\"" + strings.ReplaceAll(tempBuf.String(), "\"", "\\\"") + "\""))
		return nil
	}
	err := encodeValue(kind, buf, stack, idx, depth, vs, stringify)
	if err != nil {
		return fmt.Errorf("field %s: %v", name, err)
	}
	return nil
}

func encodeValue(kind reflect.Kind, buf io.Writer, stack []int, idx int,
	depth int, vs interface{}, stringify bool) error {
	switch kind {
	case reflect.Struct:
		fmt.Fprintf(buf, "{")
		if err := encodeStruct(buf, append(stack, idx), depth+1, vs); err != nil {
			return err
		}
		fmt.Fprintf(buf, "}")
	case reflect.String:
		fmt.Fprintf(buf, "\"%s\"", vs.(string))
	case reflect.Int:
		fmt.Fprintf(buf, "%d", vs.(int))
	case reflect.Int16:
		fmt.Fprintf(buf, "%d", vs.(int16))
	case reflect.Int32:
		fmt.Fprintf(buf, "%d", vs.(int32))
	case reflect.Int64:
		fmt.Fprintf(buf, "%d", vs.(int64))
	case reflect.Int8:
		fmt.Fprintf(buf, "%d", vs.(int8))
	case reflect.Uint:
		fmt.Fprintf(buf, "%d", vs.(uint))
	case reflect.Uint16:
		fmt.Fprintf(buf, "%d", vs.(uint16))
	case reflect.Uint32:
		fmt.Fprintf(buf, "%d", vs.(uint32))
	case reflect.Uint64:
		fmt.Fprintf(buf, "%d", vs.(uint64))
	case reflect.Uint8:
		fmt.Fprintf(buf, "%d", vs.(uint8))
	case reflect.Bool:
		fmt.Fprintf(buf, "%v", vs.(bool))
	case reflect.Float32:
		fmt.Fprintf(buf, "%v", vs.(float32))
	case reflect.Float64:
		fmt.Fprintf(buf, "%v", vs.(float64))
	case reflect.Uintptr:
		fmt.Fprintf(buf, "%v", vs.(uintptr))
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		fmt.Fprintf(buf, "[")
		vo := reflect.ValueOf(vs)
		len := vo.Len()
		for i := 0; i < len; i++ {
			if i > 0 {
				fmt.Fprintf(buf, ",")
			}
			element := vo.Index(i)
			kind := element.Kind()
			err := encodeValue(kind, buf, stack, idx, depth, element.Interface(), stringify)
			if err != nil {
				return fmt.Errorf("failed to encode slice: %w", err)
			}
		}
		fmt.Fprintf(buf, "]")
	case reflect.Map:
		fmt.Fprintf(buf, "{")
		vo := reflect.ValueOf(vs)
		for idx, key := range vo.MapKeys() {
			if idx > 0 {
				fmt.Fprintf(buf, ",")
			}
			kind := key.Kind()
			err := encodeValue(kind, buf, stack, idx, depth, key.Interface(), stringify)
			if err != nil {
				return fmt.Errorf("failed to key of map: %w", err)
			}
			fmt.Fprintf(buf, ":")
			element := vo.MapIndex(key)
			kind = element.Kind()
			err = encodeValue(kind, buf, stack, idx, depth, element.Interface(), stringify)
			if err != nil {
				return fmt.Errorf("failed to encode value of slice: %w", err)
			}
		}
		fmt.Fprintf(buf, "}")
	case reflect.Pointer:
		vo := reflect.ValueOf(vs)
		rv := vo.Elem()
		if vo.IsZero() || vo.IsNil() || !rv.CanInterface() {
			fmt.Fprintf(buf, "null")
			return nil
		}
		kind := rv.Kind()
		vs = rv.Interface()
		if kind != reflect.Struct {
			in := reflect.Indirect(rv)
			in.Kind()
		}
		return encodeValue(kind, buf, stack, idx, depth, vs, stringify)
	default:
		return fmt.Errorf("unsupported field kind %s", kind)
	}
	return nil
}
