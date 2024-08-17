package types

import (
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type TypeLmao struct {
	Something string
}

func DumpTypes() {
	converter := typescriptify.New().Add(TypeLmao{}).WithInterface(true).WithBackupDir("")
	err := converter.ConvertToFile("../types.ts")
	if err != nil {
		panic(err.Error())
	}
}
