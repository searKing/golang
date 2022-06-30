package main

/*
void crash_now(int a, char *b);
*/
import "C"

import (
	"fmt"

	//_ "github.com/ianlancetaylor/cgosymbolizer"
	_ "github.com/searKing/golang/go/runtime/cgosymbolizer"
)

// https://groups.google.com/g/golang-nuts/c/NNOJ2iiuPnY
func crash() {
	C.crash_now(1, C.CString("some string"))
}

func main() {
	fmt.Println("Pre-crash")
	crash()
	fmt.Println("Post-crash")
}
