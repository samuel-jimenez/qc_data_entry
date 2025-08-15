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

//TODO logerror
// if err != nil {
// 	log.Printf("Err: [%s]: %q\n", proc_name, err)

// log.Printf("Error: [%s]: %q\n", proc_name, err)
// }
// (proc_name string, statement *sql.Stmt, args ...any) error {
