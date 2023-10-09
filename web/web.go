package web

import "embed"

//go:embed dist/assets/*
var Assets embed.FS

//go:embed src/statics/*
var Statics embed.FS

//go:embed dist/index.html
var IndexByte []byte

//go:embed favicon.ico
var Favicon embed.FS
