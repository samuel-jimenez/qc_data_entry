package util

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/samuel-jimenez/windigo"
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

// /?TODO cf DB.Select_Panic_ErrorBox
func LogCry(proc_name, fmt_string string, args ...any) {
	message := fmt.Sprintf(fmt_string, args...)
	err := errors.New(message)
	windigo.Error(nil, message)
	log.Printf("Critical hit: [%s]: %q\n", proc_name, err)
	panic(err)
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

// func Atof(str string) float64 {
// strconv.ParseFloat("3.1415", 64)

// multiple-value-list
func µ(a ...any) []any {
	return a
}

// gO caNt HanDLe
// MUltIPLE-ValUe iN SIngLE-vaLUE coNTExT
func VALUES[T any](val T, a ...any) T {
	return val
}
