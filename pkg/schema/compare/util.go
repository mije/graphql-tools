package compare

import (
	"reflect"

	"github.com/vektah/gqlparser/ast"
)

type result struct {
	added   []interface{}
	removed []interface{}
	common  []struct{ x, y interface{} }
}

func compareSlices(x, y interface{}, equalFunc func(x, y interface{}) bool) result {
	res := new(result)

	xv := reflect.ValueOf(x)
	if xv.Kind() != reflect.Slice {
		panic("x is not a slice")
	}

	yv := reflect.ValueOf(y)
	if yv.Kind() != reflect.Slice {
		panic("y is not a slice")
	}

	if xv.Type().Elem() != yv.Type().Elem() {
		panic("incompatible types of slice elements")
	}

xloop:
	for i := 0; i < xv.Len(); i++ {
		xve := xv.Index(i)
		xe := xve.Interface()
		for j := 0; j < yv.Len(); j++ {
			yve := yv.Index(j)
			ye := yve.Interface()
			if equalFunc(xe, ye) {
				res.common = append(res.common, struct{ x, y interface{} }{xe, ye})
				continue xloop
			}
		}
		res.removed = append(res.removed, xe)
	}

yloop:
	for j := 0; j < yv.Len(); j++ {
		yve := yv.Index(j)
		ye := yve.Interface()
		for i := 0; i < xv.Len(); i++ {
			xve := xv.Index(i)
			xe := xve.Interface()
			if equalFunc(xe, ye) {
				continue yloop
			}
		}
		res.added = append(res.added, ye)
	}

	return *res
}

func compareMaps(x, y interface{}) result {
	res := new(result)

	xv := reflect.ValueOf(x)
	if xv.Kind() != reflect.Map {
		panic("x is not a map")
	}

	yv := reflect.ValueOf(y)
	if yv.Kind() != reflect.Map {
		panic("y is not a map")
	}

	xt, yt := xv.Type(), yv.Type()
	if xt.Key() != yt.Key() {
		panic("incompatible types of map keys")
	}
	if xt.Elem() != yt.Elem() {
		panic("incompatible types of map values")
	}

	for _, k := range xv.MapKeys() {
		xve, yve := xv.MapIndex(k), yv.MapIndex(k)
		if yve.IsValid() {
			res.common = append(res.common, struct {
				x, y interface{}
			}{
				x: xve.Interface(),
				y: yve.Interface(),
			})
		} else {
			res.removed = append(res.removed, xve.Interface())
		}
	}

	for _, k := range yv.MapKeys() {
		xve, yve := xv.MapIndex(k), yv.MapIndex(k)
		if !xve.IsValid() {
			res.added = append(res.added, yve.Interface())
		}
	}

	return *res
}

func typeEquals(x, y *ast.Type) bool {
	if x == nil && y == nil {
		return true
	}
	if x != nil && y != nil {
		return x.String() == y.String()
	}
	return false
}

func isBreakingTypeChange(x, y *ast.Type, isInputField bool) bool {

	if x == nil && y == nil {
		return false
	}

	if x != nil && y != nil {
		if x.NamedType != "" { // x is not a list
			if x.NamedType != y.NamedType { // y is must be not be a list and of same type
				return true
			}
			if x.NonNull && !y.NonNull { // y is of same type, but is nullable now
				return !isInputField // this is OK only for for input fields
			}
			if !x.NonNull && y.NonNull { // y is of same type, but is mandatory now
				return isInputField // this is OK only for output fields
			}
			return false // OK otherwise
		}

		if y.NamedType != "" { // x is a list, y must be a list as well
			return true
		}

		return isBreakingTypeChange(x.Elem, y.Elem, isInputField)
	}

	return true
}

func valueEquals(x, y *ast.Value) bool {
	if x == nil && y == nil {
		return true
	}
	if x != nil && y != nil {
		return x.String() == y.String()
	}
	return false
}
