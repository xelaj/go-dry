// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"fmt"
	"reflect"
	"sort"
	"unicode"
)

// ReflectTypeOfError returns the built-in error type
func ReflectTypeOfError() reflect.Type {
	return reflect.TypeOf((*error)(nil)).Elem()
}

// ReflectSetStructFieldString sets the field with name to value.
func ReflectSetStructFieldString(structPtr any, name, value string) error {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("structPtr must be pointer to a struct, but is %T", structPtr)
	}
	v = v.Elem()

	if f := v.FieldByName(name); f.IsValid() {
		if f.Kind() == reflect.String {
			f.SetString(value)
		} else {
			_, err := fmt.Sscan(value, f.Addr().Interface())
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("%T has no struct field '%s'", v.Interface(), name)
	}

	return nil
}

// ReflectSetStructFieldsFromStringMap sets the fields of a struct
// with the field names and values taken from a map[string]string.
// If errOnMissingField is true, then all fields must exist.
func ReflectSetStructFieldsFromStringMap(structPtr any, m map[string]string, errOnMissingField bool) error {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("structPtr must be pointer to a struct, but is %T", structPtr)
	}
	v = v.Elem()

	for name, value := range m {
		if f := v.FieldByName(name); f.IsValid() {
			if f.Kind() == reflect.String {
				f.SetString(value)
			} else {
				_, err := fmt.Sscan(value, f.Addr().Interface())
				if err != nil {
					return err
				}
			}
		} else if errOnMissingField {
			return fmt.Errorf("%T has no struct field '%s'", v.Interface(), name)
		}
	}

	return nil
}

/*
ReflectExportedStructFields returns a map from exported struct field names to values,
inlining anonymous sub-structs so that their field names are available
at the base level.
Example:
	type A struct {
		X int
	}
	type B Struct {
		A
		Y int
	}
	// Yields X and Y instead of A and Y:
	ReflectExportedStructFields(reflect.ValueOf(B{}))
*/
func ReflectExportedStructFields(v reflect.Value) map[string]reflect.Value {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("Expected a struct, got %s", t))
	}
	result := make(map[string]reflect.Value)
	reflectExportedStructFields(v, t, result)
	return result
}

func reflectExportedStructFields(v reflect.Value, t reflect.Type, result map[string]reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		if ReflectStructFieldIsExported(structField) {
			if structField.Anonymous && structField.Type.Kind() == reflect.Struct {
				reflectExportedStructFields(v.Field(i), structField.Type, result)
			} else {
				result[structField.Name] = v.Field(i)
			}
		}
	}
}

func ReflectNameIsExported(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}

func ReflectStructFieldIsExported(structField reflect.StructField) bool {
	return structField.PkgPath == ""
}

// ReflectSort will sort slice according to compareFunc using reflection.
// slice can be a slice of any element type including interface{}.
// compareFunc must have two arguments that are assignable from
// the slice element type or pointers to such a type.
// The result of compareFunc must be a bool indicating
// if the first argument is less than the second.
// If the element type of slice is interface{}, then the type
// of the compareFunc arguments can be any type and dynamic
// casting from the interface value or its address will be attempted.
func ReflectSort(slice, compareFunc any) {
	sortable, err := newReflectSortable(slice, compareFunc)
	if err != nil {
		panic(err)
	}
	sort.Sort(sortable)
}

func newReflectSortable(slice, compareFunc any) (*reflectSortable, error) {
	t := reflect.TypeOf(compareFunc)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("compareFunc must be a function, got %T", compareFunc)
	}
	if t.NumIn() != 2 {
		return nil, fmt.Errorf("compareFunc must take two arguments, got %d", t.NumIn())
	}
	if t.In(0) != t.In(1) {
		return nil, fmt.Errorf("compareFunc's arguments must be identical, got %s and %s", t.In(0), t.In(1))
	}
	if t.NumOut() != 1 {
		return nil, fmt.Errorf("compareFunc must have one result, got %d", t.NumOut())
	}
	if t.Out(0).Kind() != reflect.Bool {
		return nil, fmt.Errorf("compareFunc result must be bool, got %s", t.Out(0))
	}

	argType := t.In(0)
	ptrArgs := argType.Kind() == reflect.Ptr
	if ptrArgs {
		argType = argType.Elem()
	}

	sliceV := reflect.ValueOf(slice)
	if sliceV.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Need slice got %T", slice)
	}
	elemT := sliceV.Type().Elem()
	if elemT != argType && elemT.Kind() != reflect.Interface {
		return nil, fmt.Errorf("Slice element type must be interface{} or %s, got %s", argType, elemT)
	}

	return &reflectSortable{
		Slice:       sliceV,
		CompareFunc: reflect.ValueOf(compareFunc),
		ArgType:     argType,
		PtrArgs:     ptrArgs,
	}, nil
}

type reflectSortable struct {
	Slice       reflect.Value
	CompareFunc reflect.Value
	ArgType     reflect.Type
	PtrArgs     bool
}

func (r *reflectSortable) Len() int {
	return r.Slice.Len()
}

func (r *reflectSortable) Less(i, j int) bool {
	arg0 := r.Slice.Index(i)
	arg1 := r.Slice.Index(j)
	if r.Slice.Type().Elem().Kind() == reflect.Interface {
		arg0 = arg0.Elem()
		arg1 = arg1.Elem()
	}
	if (arg0.Kind() == reflect.Ptr) != r.PtrArgs {
		if r.PtrArgs {
			// Expects PtrArgs for reflectSortable, but Slice is value type
			arg0 = arg0.Addr()
		} else {
			// Expects value type for reflectSortable, but Slice is PtrArgs
			arg0 = arg0.Elem()
		}
	}
	if (arg1.Kind() == reflect.Ptr) != r.PtrArgs {
		if r.PtrArgs {
			// Expects PtrArgs for reflectSortable, but Slice is value type
			arg1 = arg1.Addr()
		} else {
			// Expects value type for reflectSortable, but Slice is PtrArgs
			arg1 = arg1.Elem()
		}
	}
	return r.CompareFunc.Call([]reflect.Value{arg0, arg1})[0].Bool()
}

func (r *reflectSortable) Swap(i, j int) {
	temp := r.Slice.Index(i)
	r.Slice.Index(i).Set(r.Slice.Index(j))
	r.Slice.Index(j).Set(temp)
}

func ReflectIsInteger(i any) bool {
	switch i.(type) {
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint:
		return true
	default:
		return false
	}
}

func ReflectIsIntgerType(t reflect.Type) bool {
	return t.ConvertibleTo(reflect.TypeOf(int(0)))
}

// true if a < b
func ReflectLess(a, b interface{}) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	switch k := aVal.Type().Kind(); k {
	case reflect.Bool:
		ar := aVal.Bool()
		switch k := bVal.Type().Kind(); k {
		case reflect.Bool:
			br := bVal.Bool()
			return (!ar && br)
		default:
			panic("can't compare " + aVal.Type().String() + " and " + bVal.Type().String())
		}

	case reflect.String:
		ar := aVal.String()
		switch k := bVal.Type().Kind(); k {
		case reflect.String:
			br := bVal.String()
			return ar < br
		default:
			panic("can't compare " + aVal.Type().String() + " and " + bVal.Type().String())
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ar := aVal.Int()
		switch k := bVal.Type().Kind(); k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			br := bVal.Int()
			return ar < br
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			br := bVal.Uint()
			return uint64(ar) < br
		case reflect.Float32, reflect.Float64:
			br := bVal.Float()
			return float64(ar) < br
		default:
			panic("can't compare " + aVal.Type().String() + " and " + bVal.Type().String())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ar := aVal.Uint()
		switch k := bVal.Type().Kind(); k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			br := bVal.Int()
			return ar < uint64(br)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			br := bVal.Uint()
			return ar < br
		case reflect.Float32, reflect.Float64:
			br := bVal.Float()
			return float64(ar) < br
		default:
			panic("can't compare " + aVal.Type().String() + " and " + bVal.Type().String())
		}

	case reflect.Float32, reflect.Float64:
		ar := aVal.Float()
		switch k := bVal.Type().Kind(); k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			br := bVal.Int()
			return ar < float64(br)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			br := bVal.Uint()
			return ar < float64(br)
		case reflect.Float32, reflect.Float64:
			br := bVal.Float()
			return ar < br
		default:
			panic("can't compare " + aVal.Type().String() + " and " + bVal.Type().String())
		}
	default:
		panic("can't compare " + aVal.Type().String() + " and " + bVal.Type().String())
		//TODO: if type implements Less method, than use it
	}
}
