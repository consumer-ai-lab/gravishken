package main

import (
	types "common"
	"log"

	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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
			app, err := newApp()
			if err != nil {
				log.Fatal(err)
			}
			defer app.destroy()
			app.openWv()
			app.prepareEnv()
			go app.serve()
			app.wait()
		},
	}
	command.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "testing command",
		Run: func(cmd *cobra.Command, args []string) {
			panic("TODO")
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := newApp()
			if err != nil {
				log.Fatal(err)
			}
			defer app.destroy()
			app.serve()
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "app",
		Short: "launch app",
		Run: func(cmd *cobra.Command, args []string) {
			command.Run(cmd, args)
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
