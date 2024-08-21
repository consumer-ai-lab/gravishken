package main

import (
	assets "app"

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

func handleClient(ws *websocket.Conn) {
	defer ws.Close()

	pingTick := time.NewTicker(time.Millisecond * 1000)
	defer pingTick.Stop()

	for {
		select {
		case <-pingTick.C:
			// ws.SetWriteDeadline(time.Now().Add(time.Second * 2))
			err := ws.WriteJSON("string json sheesh")
			// 
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
		var msg UserTest
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(msg)
		// Parse the message and add the data to the DataStore [added by kurve, just for testing]
		user := UserTest{
			UserID: msg.UserID,
			TestID: msg.TestID,
			StartTime: msg.StartTime,
			EndTime: msg.EndTime,
			ElapsedTime: msg.ElapsedTime,
			SubmissionReceived: msg.SubmissionReceived,
			ReadingElapsedTime: msg.ReadingElapsedTime,
			ReadingSubmissionReceived: msg.ReadingSubmissionReceived,
			SubmissionFolderID: msg.SubmissionFolderID,
			MergedFileID: msg.MergedFileID,
			WPM: msg.WPM,
			WPMNormal: msg.WPMNormal,
			ResultDownloaded: msg.ResultDownloaded,
		}
		AddDataToStore(user)
		log.Println("Added data to the store:", user)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
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

func server() {
	fmt.Println("Starting server...")

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", serveWs)

	// Add the new route to your server [added by kurve, just for testing]
	mux.HandleFunc("/data", getDataFromStore)

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

	log.Fatal(http.ListenAndServe("localhost:6200", mux))
}
