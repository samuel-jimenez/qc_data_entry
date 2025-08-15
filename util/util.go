package util

import (
	"fmt"
	"log"
)

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

func LogError(proc_name string, err error) {
	if err != nil {
		log.Printf("Error: [%s]: %q\n", proc_name, err)
	}
}
