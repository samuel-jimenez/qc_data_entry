package math

// because  go sucks

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float |
		int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
	// go sucks
}

// func Abs(val0 any) any {
// 	return max(val0, -val0)
//
// }
// func Abs(val0 Number) Number {
// 	return max(val0, -val0)
//
// }

func Abs[T Number](val0 T) T {
	return max(val0, -val0)
}
