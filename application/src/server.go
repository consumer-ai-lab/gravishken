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

func (self *App) serve() {
	handleClient := func(ws *websocket.Conn, ctx context.Context) {
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

	handleMessages := func(ws *websocket.Conn, close context.CancelFunc) {
		defer ws.Close()

		for {
			var msg types.Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				close()
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

		ctx, close := context.WithCancel(context.Background())

		log.Println("new conn")
		go handleClient(ws, ctx)
		handleMessages(ws, close)
	}

	mux := http.NewServeMux()

	// TODO: more than 1 websocket client at the same time is not supported. maybe crash / don't accept the connection
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

	log.Printf("Starting application on port %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%s", port), mux)
	log.Fatal(err)
}

func (self *App) handleMessages() {
	for {
		msg, ok := <-self.recv
		if !ok {
			return
		}

		switch msg.Typ {
		case types.LoadRoute:
			val, err := types.Get[types.TLoadRoute](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			self.send <- types.NewMessage(*val)
		case types.UserLogin:
			val, err := types.Get[types.TUserLogin](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			err = self.login(*val)
			if err != nil {
				self.notifyErr(err)
				continue
			}
		case types.GetTest:
			val, err := types.Get[types.TGetTest](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			err = self.startTest(*val)
			if err != nil {
				self.notifyErr(err)
				continue
			}
		case types.OpenApp:
			val, err := types.Get[types.TOpenApp](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			// TODO: use fixed paths instead of generating a random path
			// this will help in cases where someone restarts the test
			dest, err := self.runner.NewTemplate(val.Typ)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			go (func() {
				err = self.runner.FocusOrOpenApp(val.Typ, dest)
				self.notifyErr(err)
			})()
		case types.Err:
			val, err := types.Get[types.TErr](msg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(val)
		case types.ExeNotFound:
			log.Printf("message of type '%s' cannot be handled here: '%s'\n", msg.Typ.TSName(), msg.Val)
		case types.Unknown:
			log.Printf("unknown message type received: '%s'\n", msg.Val)
		// TODO:
		// case types.UserTestSomethingSomething:
		// 	user, err := types.Get[UserTest](msg)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	AddDataToStore(*user)
		// 	log.Println("Added data to the store:", user)
		default:
			log.Printf("message type '%s' not handled ('%s')\n", msg.Typ.TSName(), msg.Val)
		}
	}
}
