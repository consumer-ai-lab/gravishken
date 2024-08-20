package main

import (
	types "common"
	"path/filepath"

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

type Error struct {
	message string
}

func NewError(msg string) Error {
	return Error{message: msg}
}

func (self *Error) Error() string {
	return fmt.Sprintf("Error: %s", self.message)
}

func app() {
	w := webview.New(build_mode == "DEV")
	defer w.Destroy()

	w.SetTitle("gravishken")

	url := fmt.Sprintf("http://localhost:%s/", port)
	w.Navigate(url)

	w.Run()
	w.Destroy()
}

func main() {
	if build_mode == "DEV" {
		root, ok := os.LookupEnv("PROJECT_ROOT")
		if !ok {
			panic("'PROJECT_ROOT' not set")
		}
		ts_dir := filepath.Join(root, "common", "ts")
		types.DumpTypes(ts_dir)
	}

	var command = &cobra.Command{
		// default action
		Run: func(cmd *cobra.Command, args []string) {
			go server()
			go app()
			test()
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
	command.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "testing command",
		Run: func(cmd *cobra.Command, args []string) {
			go app()
			test()
			// types.Test()
		},
	})

	// - [windows app start error](https://github.com/spf13/cobra/issues/844)
	cobra.MousetrapHelpText = ""
	var err = command.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
