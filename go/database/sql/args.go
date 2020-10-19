package sql

import "time"

func StringSliceArgs(src ...string) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func BoolSliceArgs(src ...bool) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func TimeSliceArgs(src ...time.Time) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func IntSliceArgs(src ...int) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func UintSliceArgs(src ...uint) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Int8SliceArgs(src ...int8) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Uint8SliceArgs(src ...uint8) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Int16SliceArgs(src ...int16) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Uint16SliceArgs(src ...uint16) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Int32SliceArgs(src ...int32) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Uint32SliceArgs(src ...uint32) []interface{} {
	var args []interface{}
	for _, s := range src {
		args = append(args, s)
	}
	return args
}
