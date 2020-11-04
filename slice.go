// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"fmt"
	"reflect"
)

// индекс элемента T в []T
// выводит индекс найденого элемента, либо -1, если элемент не найден
func SliceIndex(slice, item any) int {
	ival := reflect.ValueOf(slice)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}
	if ival.Type().Elem().String() != reflect.TypeOf(item).String() {
		panic("different types of slice and item")
	}

	for i := 0; i < ival.Len(); i++ {
		if reflect.DeepEqual(ival.Index(i).Interface(), item) {
			return i
		}
	}

	return -1
}

func SliceContains(slice, item any) bool {
	return SliceIndex(slice, item) != -1
}

func DeleteIndex(slice any, i int) any {
	return cutSliceWithMode(slice, i, i+1, true)
}

func SliceCut(slice any, i, j int) any {
	return cutSliceWithMode(slice, i, j, false)
}

func cutSliceWithMode(slice any, i, j int, deleteInsteadCut bool) any {
	panicIndexStr := fmt.Sprintf("[%v:%v]", i, j)
	if deleteInsteadCut {
		panicIndexStr = fmt.Sprintf("[%v]", i)
	}

	ival := reflect.ValueOf(slice)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	if i > j {
		panic("end less than start " + panicIndexStr)
	}
	if i < 0 || j < 0 {
		panic("slice index " + panicIndexStr + " out of bounds")
	}
	if ival.Len()-1 < i || ival.Len() < j {
		panic(fmt.Sprintf("index out of range %v with length %v", panicIndexStr, ival.Len()))
	}

	return reflect.AppendSlice(ival.Slice(0, i), ival.Slice(j, ival.Len())).Interface()
}

func SliceExpand(slice any, i, j int) any {
	panicIndexStr := fmt.Sprintf("[%v]", i)

	ival := reflect.ValueOf(slice)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	if i < 0 {
		panic("slice index " + panicIndexStr + " out of bounds")
	}
	if j < 0 {
		panic(fmt.Sprintf("can't expand slice on %v points", j))
	}
	if ival.Len()-1 < i {
		panic(fmt.Sprintf("index out of range %v with length %v", panicIndexStr, ival.Len()))
	}

	zeroitems := reflect.MakeSlice(ival.Type(), j, j)
	part := reflect.AppendSlice(zeroitems, ival.Slice(i, ival.Len()))
	return reflect.AppendSlice(ival.Slice(0, i), part).Interface()
}

func SliceToInterfaceSlice(in any) []any {
	if in == nil {
		return nil
	}

	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	res := make([]any, ival.Len())

	for i := 0; i < ival.Len(); i++ {
		res[i] = ival.Index(i).Interface()
	}
	return res
}

// map[<K>]<V> -> []<K> // (K, V) could be any type
func MapKeys(in any) any {
	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Map {
		panic("not a map: " + ival.Type().String())
	}

	keys := ival.MapKeys()

	items := reflect.MakeSlice(reflect.SliceOf(ival.Type().Key()), len(keys), len(keys))
	for i, key := range keys {
		items.Index(i).Set(key)
	}

	return items.Interface()
}

// []<T> -> map[<T>]struct{}
func SliceUnique(in any) any {
	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	res := reflect.MakeMap(reflect.MapOf(ival.Type().Elem(), reflect.TypeOf(null{})))

	for i := 0; i < ival.Len(); i++ {
		res.SetMapIndex(ival.Index(i), reflect.ValueOf(null{}))
	}
	return res.Interface()
}

func SetUnify(a, b any) any {
	aval := reflect.ValueOf(a)
	if aval.Type().Kind() != reflect.Slice {
		panic("first element is not a slice: " + aval.Type().String())
	}
	bval := reflect.ValueOf(b)
	if bval.Type().Kind() != reflect.Slice {
		panic("second element is not a slice: " + bval.Type().String())
	}

	if aval.Type().Elem() != bval.Type().Elem() {
		panic("slices has different types")
	}

	pretotal := reflect.MakeMap(reflect.MapOf(aval.Type().Elem(), reflect.TypeOf(null{})))

	for i := 0; i < aval.Len(); i++ {
		pretotal.SetMapIndex(aval.Index(i), reflect.ValueOf(null{}))
	}
	for i := 0; i < bval.Len(); i++ {
		pretotal.SetMapIndex(bval.Index(i), reflect.ValueOf(null{}))
	}

	res := reflect.New(aval.Type()).Elem()

	for _, key := range pretotal.MapKeys() {
		res = reflect.Append(res, key)
	}

	return res.Interface()
}

// пока только строчки, потом сделаю нормальный интерфейс
func SetsEqual(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
