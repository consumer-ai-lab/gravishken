//go:build windows

package main

/*
// OOF: T_T T_T T_T this somehow fixes the compile issue on windows and i don't know why T_T T_T
// #cgo windows LDFLAGS: -static -lpthread

#include <pthread.h>
void* threadFunction(void* arg) {
    return NULL;
}
void createThread(int* arg) {
    pthread_t thread;
    pthread_create(&thread, NULL, threadFunction, arg);
    pthread_join(thread, NULL); // Wait for the thread to finish
}
*/
// import "C"
