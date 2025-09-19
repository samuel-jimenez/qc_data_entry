package util

import (
	"fmt"
	"log"
	"strings"
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

func Concat(str ...string) string {
	num_Str := len(str)
	Builder := strings.Builder{}
	Builder_len := 0
	for i := range num_Str {
		Builder_len += len(str[i])
	}

	Builder.Grow(Builder_len)
	for i := range num_Str {
		Builder.WriteString(str[i])
	}

	return Builder.String()
}
