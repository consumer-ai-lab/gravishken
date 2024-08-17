package main

import (
	assets "gravtest"
	"time"

	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/thrombe/webview_go"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleClient(ws *websocket.Conn) {
	defer ws.Close()

	pingTick := time.NewTicker(time.Millisecond * 1000)
	defer pingTick.Stop()

	for {
		select {
		case <-pingTick.C:
			// ws.SetWriteDeadline(time.Now().Add(time.Second * 2))
			err := ws.WriteJSON("string json sheesh")
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func handleMessages(ws *websocket.Conn) {
	defer ws.Close()

	for {
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(msg)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("new conn")
	go handleClient(ws)
	handleMessages(ws)
}

func server() {
	fmt.Println("Starting server...")

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", serveWs)

	if build_mode == "PROD" {
		build, _ := fs.Sub(assets.Assets, "dist")
		fileServer := http.FileServer(http.FS(build))
		mux.Handle("/", fileServer)
	} else if build_mode == "DEV" {
		build := http.Dir("dist")
		fileServer := http.FileServer(build)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("cache-control", "no-store, no-cache, must-revalidate")
			fileServer.ServeHTTP(w, r)
		})
	} else {
		panic("invalid BUILD_MODE")
	}

	log.Fatal(http.ListenAndServe("localhost:6200", mux))
}

func app() {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("GravTest")

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
