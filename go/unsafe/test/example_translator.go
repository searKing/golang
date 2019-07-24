package test

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

char** getStringArray(int num){
	char** list = malloc(num*sizeof(char*));
	for(int i = 0; i<num; i++){
		list[i] = "searKing";
	}
	return list;
}

void showStringArray(char** list, int num){
	printf("num = %d\n",num);
	int i = 0;
	for(i=0;i<num;i++){
		printf("%s\n",list[i]);
	}
}
*/
import "C"
import (
	"fmt"
	unsafe2 "github.com/searKing/golang/go/unsafe"
	"unsafe"
)

func ExampleGoStringArray() {
	//char** 转化成 []string
	cCharArray := C.getStringArray(3)
	defer C.free(unsafe.Pointer(cCharArray))
	fmt.Print(unsafe2.GoStringArray(unsafe.Pointer(cCharArray), 3))

}
func ExampleCStringArray() {
	//[]string 转化成 char**
	box := []string{"xing", "jack", "john", "searKing"}
	cCharArray, n := unsafe2.CStringArray(box...)
	defer C.free(unsafe.Pointer(cCharArray))
	C.showStringArray((**C.char)(unsafe.Pointer(cCharArray)), C.int(n))
}
