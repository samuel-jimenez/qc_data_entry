package math

// because  go sucks

import (
	"strconv"

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

// func Round[T Number](val0 T,) T {

func SigFig(val float64, figures int) float64 {
	out, _ := strconv.ParseFloat(strconv.FormatFloat(val, 'G', figures, 64), 64)
	return out
}

// // TODO fix math.Round()
// func Round(val float64, precision int) float64 {
// 	out, _ := strconv.ParseFloat(formats.Format_float(val, precision), 64)
// 	return out
// }

func FMULT[T, U Number](a T, b U) float64 {
	return float64(a) * float64(b)
}

func FDIV_[T, U Number](a T, b U) float64 {
	return float64(a) / float64(b)
}

func FDIV(a, b int) float64 {
	return float64(a) / float64(b)
}
