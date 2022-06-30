#include <stdio.h>

void crash_now(int a, char *b) {
    int *thing = NULL;
    fprintf(stderr, "Here: %d %s\n", *thing + a, b);
    // https://github.com/benesch/cgosymbolizer/blob/master/example/example.c
    volatile char p = *(volatile char*) NULL;
}