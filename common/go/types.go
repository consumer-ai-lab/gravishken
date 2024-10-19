package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

// this is a custom error type for use throughout the app
type Error struct {
	message string
}

func NewError(msg string) Error {
	return Error{message: msg}
}

func (self Error) Error() string {
	return fmt.Sprintf("Error: %s", self.message)
}

type Varient int

const (
	Err Varient = iota
	ExeNotFound
	Quit
	UserLogin
	WarnUser
	LoadRoute
	ReloadUi
	StartTest
	OpenApp
	QuitApp
	Unknown // NOTE: keep this as the last constant here.
)

func (self Varient) TSName() string {
	switch self {
	case Err:
		return "Err"
	case ExeNotFound:
		return "ExeNotFound"
	case Quit:
		return "Quit"
	case UserLogin:
		return "UserLogin"
	case WarnUser:
		return "WarnUser"
	case LoadRoute:
		return "LoadRoute"
	case ReloadUi:
		return "ReloadUi"
	case StartTest:
		return "StartTest"
	case OpenApp:
		return "OpenApp"
	case QuitApp:
		return "QuitApp"
	default:
		return "Unknown"
	}
}
func varientFromName(typ string) Varient {
	switch typ {
	case "Err":
		return Err
	case "ExeNotFound":
		return ExeNotFound
	case "Quit":
		return Quit
	case "UserLogin":
		return UserLogin
	case "WarnUser":
		return WarnUser
	case "ReloadUi":
		return ReloadUi
	case "LoadRoute":
		return LoadRoute
	case "StartTest":
		return StartTest
	case "OpenApp":
		return OpenApp
	case "QuitApp":
		return QuitApp
	default:
		return Unknown
	}
}

// only for unexpected errors / for errors that we can't do much about, other than telling the user about it
type TErr struct {
	Message string
}

type Message struct {
	Typ Varient
	Val string
}

type TExeNotFound struct {
	Name   string
	ErrMsg string
}

type TQuit struct{}

type TUserLogin struct {
	Username string
	Password string
	TestCode string
}
type TWarnUser struct {
	Message string
}

type TLoadRoute struct {
	Route string
}

type TReloadUi struct{}

type TStartTestRequest struct{}

type TStartTest struct {
	tests []Test
}

type AppType int

const (
	TXT AppType = iota
	DOCX
	XLSX
	PPTX
)

func (self AppType) TSName() string {
	switch self {
	case TXT:
		return "TXT"
	case DOCX:
		return "DOCX"
	case XLSX:
		return "XLSX"
	case PPTX:
		return "PPTX"
	default:
		return "Unknown"
	}
}

type TOpenApp struct {
	Typ AppType
}

type TQuitApp struct{}

func NewMessage(typ interface{}) Message {
	name := reflect.TypeOf(typ).Name()[1:]
	varient := varientFromName(name)
	json, err := json.Marshal(typ)
	if err != nil {
		panic(err)
	}
	return Message{
		Typ: varient,
		Val: string(json),
	}
}

func Get[T any](msg Message) (*T, error) {
	var val T

	name := reflect.TypeOf(val).Name()[1:]
	if name != msg.Typ.TSName() {
		err_msg := fmt.Sprintf("message of type '%s' but asked to be decoded as '%s'", msg.Typ.TSName(), name)
		return nil, NewError(err_msg)
	}

	err := json.Unmarshal([]byte(msg.Val), &val)
	return &val, err
}

// - [tkrajina/tkypescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs)
func DumpTypes(dir string) {
	allVarients := make([]Varient, Unknown+1)

	for i := 0; i <= int(Unknown); i++ {
		allVarients[i] = Varient(i)
	}

	converter := typescriptify.New().
		WithInterface(true).
		WithBackupDir("").
		Add(TErr{}).
		Add(Message{}).
		Add(TExeNotFound{}).
		Add(TQuit{}).
		Add(TUserLogin{}).
		Add(TWarnUser{}).
		Add(TLoadRoute{}).
		Add(TReloadUi{}).
		Add(TStartTestRequest{}).
		Add(TStartTest{}).
		Add(TOpenApp{}).
		Add(TQuitApp{}).
		AddEnum([]AppType{TXT, DOCX, XLSX, PPTX}).
		AddEnum(allVarients)

	converter = converter.
		Add(User{}).
		Add(UserSubmission{}).
		Add(UserBatchRequestData{}).
		Add(UserLoginRequest{}).
		Add(Test{}).
		Add(Admin{}).
		Add(AdminRequest{}).
		Add(Batch{}).
		AddEnum([]TestType{TypingTest, DocxTest, ExcelTest, WordTest})

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err.Error())
	}
	err = converter.ConvertToFile(filepath.Join(dir, "types.ts"))
	if err != nil {
		panic(err.Error())
	}
}
