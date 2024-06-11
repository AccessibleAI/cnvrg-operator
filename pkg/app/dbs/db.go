package dbs

import (
	"embed"
)

const fsRoot = "tmpl"

//go:embed  tmpl/*
var fs embed.FS
