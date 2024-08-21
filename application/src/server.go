package main

import (
	assets "app"
	types "common"
	"context"

	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (self *App) serve() {
	fmt.Println("Starting server...")

	// TODO: this ctx is currently useless
	ctx, close := context.WithCancel(context.Background())
	defer close()

	handleClient := func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-self.send:
				if !ok {
					return
				}
				log.Println(msg)
				ws.SetWriteDeadline(time.Now().Add(time.Second * 5))
				err := ws.WriteJSON(msg)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

	handleMessages := func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			var msg types.Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				return
			}
			self.recv <- msg
		}
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	serveWs := func(w http.ResponseWriter, r *http.Request) {
		if build_mode == "DEV" {
			r.Header.Del("origin")
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("new conn")
		go handleClient(ws)
		handleMessages(ws)
	}

	mux := http.NewServeMux()

	// TODO: more than 1 websocket client at the same time is not supported. maybe crash / don't accept the connection
	mux.HandleFunc("/ws", serveWs)

	if build_mode == "PROD" {
		build, _ := fs.Sub(assets.Dist, "dist")
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

	err := http.ListenAndServe(fmt.Sprintf("localhost:%s", port), mux)
	log.Fatal(err)
}
