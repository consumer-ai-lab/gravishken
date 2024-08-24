package assets

import "embed"

/*go:embed dist/**/
var Dist embed.FS

//go:embed templates/*
var Templates embed.FS
