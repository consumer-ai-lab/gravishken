package types

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type Varient int

const (
	Var1 Varient = iota
	Var2
	Err
	Unknown
)

var allVarients = []Varient{Var1, Var2, Err, Unknown}

func (self Varient) TSName() string {
	switch self {
	case Var1:
		return "Var1"
	case Var2:
		return "Var2"
	case Err:
		return "Err"
	default:
		return "Unknown"
	}
}
func varientFromName(typ string) Varient {
	switch typ {
	case "Var1":
		return Var1
	case "Var2":
		return Var2
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

type TVar1 struct {
	Field1 int
	Field2 bool
}

type TVar2 struct {
	Field1 bool
	Field3 string
}

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
		Add(TVar1{}).
		Add(TVar2{}).
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

func Test() {
	v1 := TVar1{Field1: 42, Field2: false}
	msg := NewMessage(v1)
	log.Println(msg)
	jmsg, _ := json.Marshal(msg)
	log.Println(string(jmsg))

	back, _ := Get[TVar1](msg)
	log.Println(back)
}
