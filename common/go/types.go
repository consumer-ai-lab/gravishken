package types

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type Varient int

const (
	ExeNotFound Varient = iota
	Err
	Unknown
)

var allVarients = []Varient{ExeNotFound, Err, Unknown}

func (self Varient) TSName() string {
	switch self {
	case ExeNotFound:
		return "ExeNotFound"
	case Err:
		return "Err"
	default:
		return "Unknown"
	}
}
func varientFromName(typ string) Varient {
	switch typ {
	case "ExeNotFound":
		return ExeNotFound
	case "Err":
		return Err
	default:
		return Unknown
	}
}

type Message struct {
	Type Varient
	Val  string
}

type TExeNotFound struct {
	Name   string
	ErrMsg string
}

// only for unexpected errors / for errors that we can't do much about, other than telling the user about it
type TErr struct {
	Message string
}

func NewMessage(typ interface{}) Message {
	name := reflect.TypeOf(typ).Name()[1:]
	varient := varientFromName(name)
	json, err := json.Marshal(typ)
	if err != nil {
		panic(err)
	}
	return Message{
		Type: varient,
		Val:  string(json),
	}
}

func Get[T any](msg Message) (T, error) {
	var val T
	err := json.Unmarshal([]byte(msg.Val), &val)
	return val, err
}

// - [tkrajina/tkypescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs)
func DumpTypes(dir string) {
	converter := typescriptify.New().
		WithInterface(true).
		WithBackupDir("").
		Add(Message{}).
		Add(TExeNotFound{}).
		Add(TErr{}).
		AddEnum(allVarients)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err.Error())
	}
	err = converter.ConvertToFile(filepath.Join(dir, "types.ts"))
	if err != nil {
		panic(err.Error())
	}
}
