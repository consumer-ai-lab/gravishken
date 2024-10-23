package main

import (
	types "common"
	"log"
	"runtime"

	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// build time configuration. these get set using -ldflags in build script
var build_mode string
var port string

func main() {
	if build_mode == "DEV" {
		root, ok := os.LookupEnv("PROJECT_ROOT")
		if !ok {
			panic("'PROJECT_ROOT' not set")
		}
		ts_dir := filepath.Join(root, "common", "ts")
		types.DumpTypes(ts_dir)
	}

	// .env by default
	_ = godotenv.Overload()

	if runtime.GOOS == "windows" {
		// NOTE: APPDATA is set on windows automatically
		datadir := filepath.Join(os.Getenv("APPDATA"), "Gravishken")
		_ = os.Mkdir(datadir, os.ModePerm)

		_ = godotenv.Overload(filepath.Join(datadir, ".env"))
	}

	url, ok := os.LookupEnv("SERVER_URL")
	if ok {
		server_url = url
	}
	secure, ok := os.LookupEnv("SERVER_SECURE")
	if ok {
		server_is_secure = secure
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
			go app.handleMessages()
			go app.serve()
			app.wait()
		},
	}
	// command.AddCommand(&cobra.Command{
	// 	Use:   "test",
	// 	Short: "testing command",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		// panic("TODO")
	// 		uritaOpenWv("http://localhost:3000")
	// 	},
	// })
	command.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := newApp()
			if err != nil {
				log.Fatal(err)
			}
			defer app.destroy()
			go app.handleMessages()
			app.send <- types.NewMessage(types.TReloadUi{})
			go app.serve()
			app.wait()
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
