package main

import "fmt"

func assert(statement bool, message string) {
	if !statement {
		panic(message)
	}
}

func assertEqual(val0, val1 any) {
	if val0 != val1 {
		panic(fmt.Sprintf("Assert failed: %v !=  %v", val0, val1))
	}
}
