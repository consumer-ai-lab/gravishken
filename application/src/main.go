package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	webview "github.com/thrombe/webview_go"
)

/*
// OOF: T_T T_T T_T this somehow fixes the compile issue on windows and i don't know why T_T T_T
#cgo windows LDFLAGS: -static -lpthread

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
import "C"

var build_mode string
var port string

func app() {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("gravishken")

	url := fmt.Sprintf("http://localhost:%s/", port)
	w.Navigate(url)

	w.Run()
}

func main() {
	var command = &cobra.Command{
		// default action
		Run: func(cmd *cobra.Command, args []string) {
			go server()
			app()
		},
	}

	command.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			server()
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "app",
		Short: "launch app",
		Run: func(cmd *cobra.Command, args []string) {
			go server()
			app()
		},
	})

	var err = command.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
