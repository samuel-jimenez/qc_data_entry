package util

import "fmt"

func Assert(statement bool, message string) {
	if !statement {
		panic(message)
	}
}

func AssertEqual(val0, val1 any) {
	if val0 != val1 {
		panic(fmt.Sprintf("Assert failed: %v !=  %v", val0, val1))
	}
}
