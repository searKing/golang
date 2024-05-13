package sql

import "time"

func StringSliceArgs(src ...string) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func BoolSliceArgs(src ...bool) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func TimeSliceArgs(src ...time.Time) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func IntSliceArgs(src ...int) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func UintSliceArgs(src ...uint) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Int8SliceArgs(src ...int8) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Uint8SliceArgs(src ...uint8) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Int16SliceArgs(src ...int16) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Uint16SliceArgs(src ...uint16) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Int32SliceArgs(src ...int32) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}

func Uint32SliceArgs(src ...uint32) []any {
	var args []any
	for _, s := range src {
		args = append(args, s)
	}
	return args
}
