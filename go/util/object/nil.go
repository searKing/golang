package object

import (
	"errors"
	"reflect"
	"strings"
)

type ErrorNilPointer error
type ErrorMissMatch error

var (
	errorNilPointer = errors.New("nil pointer")
)

type Supplier interface {
	Get() interface{}
}

// IsNil returns {@code true} if the provided reference is {@code nil} otherwise
// returns {@code false}.
func IsNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if IsNilable(obj) {
		return reflect.ValueOf(obj).IsNil()
	}
	return false
}

// IsNil returns {@code true} if the provided reference is non-{@code nil} otherwise
// returns {@code false}.
func NoneNil(obj interface{}) bool {
	return !IsNil(obj)
}

// IsNil returns {@code true} if the provided reference can be assigned {@code nil} otherwise
// returns {@code false}.
func IsNilable(obj interface{}) (canBeNil bool) {
	defer func() {
		// As we can not access v.flag&reflect.flagMethod&v.ptr
		// So we use recover() instead
		if r := recover(); r != nil {
			canBeNil = false
		}
	}()
	reflect.ValueOf(obj).IsNil()

	canBeNil = true
	return
}

// RequireNonNil checks that the specified object reference is not {@code nil}. This
// method is designed primarily for doing parameter validation in methods
// and constructors
func RequireNonNil(obj interface{}, msg ...string) interface{} {
	if msg == nil {
		msg = []string{"nil pointer"}
	}
	if IsNil(obj) {
		panic(ErrorNilPointer(errors.New(strings.Join(msg, ""))))
	}
	return obj
}

// RequireNonNullElse returns the first argument if it is non-{@code nil} and
// otherwise returns the non-{@code nil} second argument.
func RequireNonNullElse(obj, defaultObj interface{}) interface{} {
	if NoneNil(obj) {
		return obj
	}
	return RequireNonNil(defaultObj, "defaultObj")
}

// RequireNonNullElseGet returns the first argument if it is non-{@code nil} and
// returns the non-{@code nil} value of {@code supplier.Get()}.
func RequireNonNullElseGet(obj interface{}, supplier Supplier) interface{} {
	if NoneNil(obj) {
		return obj
	}
	return RequireNonNil(RequireNonNil(supplier, "supplier").(Supplier).Get(), "supplier.Get()")
}

func IsEmptyValue(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// RequireNonNil checks that the specified object reference is not {@code nil}. This
// method is designed primarily for doing parameter validation in methods
// and constructors
func RequireEqual(actual, expected interface{}, msg ...string) interface{} {
	if msg == nil {
		msg = []string{"miss match"}
	}
	if !Equals(actual, expected) {
		panic(ErrorMissMatch(errors.New(strings.Join(msg, ""))))
	}
	return actual
}
