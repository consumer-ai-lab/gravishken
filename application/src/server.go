package main

import (
	assets "app"
	"common"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"net/http/httptest"

	"github.com/gorilla/websocket"
)

func (self *App) serve() {
	handleClient := func(ws *websocket.Conn, ctx context.Context) {
		defer ws.Close()

		for {
			select {
			case <-self.exitCtx.Done():
				return
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
					log.Println("frontend ws closed :/", err)
					return
				}
			}
		}
	}

	handleMessages := func(ws *websocket.Conn, ctx context.Context, close context.CancelFunc) {
		defer ws.Close()
		defer close()

		go func() {
			defer close()
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Millisecond * 1):
					//
				}
				var msg common.Message
				err := ws.ReadJSON(&msg)
				if err != nil {
					log.Println("frontend ws closed 2 :/", err)
					return
				}
				self.recv <- msg
			}
		}()

		select {
		case <-ctx.Done():
			return
		case <-self.exitCtx.Done():
			return
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
		handleMessages(ws, ctx, close)
	}

	mux := http.NewServeMux()

	// TODO: more than 1 websocket client at the same time is not supported. maybe crash / don't accept the connection
	mux.HandleFunc("/ws", serveWs)

	mux.HandleFunc("/get-user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Header().Add("access-control-allow-origin", "*")
		if err := json.NewEncoder(w).Encode(self.client.user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/get-tests", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Header().Add("access-control-allow-origin", "*")
		if err := json.NewEncoder(w).Encode(self.client.tests); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/get-submitted-ids", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Header().Add("access-control-allow-origin", "*")
		if err := json.NewEncoder(w).Encode(self.test_state.submitted); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/submit-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("access-control-allow-origin", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Allow", "POST, OPTIONS")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var submission common.TestSubmission
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&submission); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if self.runner.IsAppOpen() {
			_ = self.runner.FocusOpenApp()
			msg := "An Application is open. Please save your work, close the app and retry"
			self.notifyErr(fmt.Errorf(msg))
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		test, err := self.findTestById(submission.TestId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch test.Type {
		case common.TypingTest, common.MCQTest:
		case common.DocxTest, common.ExcelTest, common.PptTest:
			path, ok := self.test_state.tests[test.Id]
			if !ok {
				self.notifyErr(fmt.Errorf("No data found. Did you complete the test?"))
				http.Error(w, "No test data found", http.StatusBadRequest)
				return
			}

			file, err := os.Open(path)
			if err != nil {
				err = fmt.Errorf("failed opening test file")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()

			filedata, err := io.ReadAll(file)
			if err != nil {
				err = fmt.Errorf("failed reading test file")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			encodedData := base64.StdEncoding.EncodeToString(filedata)
			info := common.AppTestInfo{
				FileData: encodedData,
			}
			switch test.Type {
			case common.DocxTest:
				submission.TestInfo.DocxTestInfo = &info
			case common.ExcelTest:
				submission.TestInfo.ExcelTestInfo = &info
			case common.PptTest:
				submission.TestInfo.PptTestInfo = &info
			default:
				panic("unreachable")
			}
		default:
			if err != nil {
				http.Error(w, "Unknown test", http.StatusBadRequest)
				return
			}
		}

		err = self.client.submitTest(submission)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		self.test_state.submitted[submission.TestId] = true

		self.maybeFinishTest()
		w.WriteHeader(http.StatusNoContent)
		self.send <- common.NewMessage(common.TNotification{
			Message: fmt.Sprintf("Test submitted Sucessfully"),
			Typ:     "success",
		})
	})

	var contentReplacements = map[string]string{
		"%SERVER_URL%": os.Getenv("SERVER_URL"),
		"%APP_PORT%":   port,
	}

	var httpFS http.FileSystem
	if build_mode == "PROD" {
		build, _ := fs.Sub(assets.Dist, "dist")
		httpFS = http.FS(build)
	} else if build_mode == "DEV" {
		httpFS = http.Dir("dist")
	} else {
		panic("invalid BUILD_MODE")
	}

	fileServer := http.FileServer(httpFS)

	modifiedFileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the file and capture its content
		recorder := httptest.NewRecorder()
		fileServer.ServeHTTP(recorder, r)

		content := recorder.Body.String()

		for oldString, newString := range contentReplacements {
			content = strings.ReplaceAll(content, oldString, newString)
		}

		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		if build_mode == "DEV" {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		}
		w.WriteHeader(recorder.Code)
		w.Write([]byte(content))
	})

	mux.Handle("/", modifiedFileServer)

	go func() {
		log.Printf("Starting application on port %s...\n", port)
		err := http.ListenAndServe(fmt.Sprintf("localhost:%s", port), mux)
		log.Fatal(err)
	}()
	<-self.exitCtx.Done()
}

func (self *App) findTestById(id common.ID) (*common.Test, error) {
	for _, test := range self.client.tests {
		if test.Id == id {
			return &test, nil
		}
	}
	return nil, fmt.Errorf("Unknown test")
}

func (self *App) maybeFinishTest() {
	for _, test := range self.client.tests {
		_, ok := self.test_state.submitted[test.Id]
		if !ok {
			return
		}
	}

	self.send <- common.NewMessage(common.TTestFinished{})
}

func (self *App) handleMessages() {
	for {
		var msg common.Message
		var ok bool
		select {
		case <-self.exitCtx.Done():
			return
		case msg, ok = <-self.recv:
			if !ok {
				return
			}
		}

		log.Println(msg.Typ.TSName(), msg)

		switch msg.Typ {
		case common.LoadRoute:
			val, err := common.Get[common.TLoadRoute](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			self.send <- common.NewMessage(*val)
		case common.UserLoginRequest:
			val, err := common.Get[common.TUserLoginRequest](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			err = self.login(val)
			if err != nil {
				self.notifyErr(err)
				continue
			}

			self.maintainConnection()

			routeMessage := common.TLoadRoute{
				Route: "/instructions",
			}
			message := common.NewMessage(routeMessage)
			self.send <- message
		case common.StartTest:
			err := self.startTest()
			if err != nil {
				self.notifyErr(err)
				continue
			}
		case common.CheckSystem:
			self.runner.CheckApps()
		case common.OpenApp:
			val, err := common.Get[common.TOpenApp](msg)
			if err != nil {
				self.notifyErr(err)
				continue
			}
			var dest string
			// TODO: self.state.tests will be wiped if app restarts. :) but i don't care rn
			dest, ok := self.test_state.tests[val.TestId]
			if !ok {
				dest, err = self.runner.NewTemplate(val.Typ)
				if err != nil {
					self.notifyErr(err)
					continue
				}
				self.test_state.tests[val.TestId] = dest
			}
			go (func() {
				err = self.runner.FocusOrOpenApp(val.Typ, dest)
				self.notifyErr(err)
			})()
		case common.Quit:
			self.exitFn()
		case common.QuitApp:
			self.runner.KillApp()
		case common.Err:
			val, err := common.Get[common.TErr](msg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(val)
		case common.ExeNotFound, common.TestFinished:
			log.Printf("message of type '%s' cannot be handled here: '%s'\n", msg.Typ.TSName(), msg.Val)
		case common.Unknown:
			log.Printf("unknown message type received: '%s'\n", msg.Val)
		default:
			log.Printf("message type '%s' not handled ('%s')\n", msg.Typ.TSName(), msg.Val)
		}
	}
}
