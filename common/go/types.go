package types

import (
	"os"
	"path/filepath"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type TypeLmao struct {
	Something string
}

// - [tkrajina/typescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs)
func DumpTypes(dir string) {
	converter := typescriptify.New().Add(TypeLmao{}).WithInterface(true).WithBackupDir("")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err.Error())
	}
	err = converter.ConvertToFile(filepath.Join(dir, "types.ts"))
	if err != nil {
		panic(err.Error())
	}
}
