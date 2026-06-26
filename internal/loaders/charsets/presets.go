package charsets

import (
	"embed"
	"maps"
	"strings"
)

//go:embed assets/*
var charsetFiles embed.FS

var builtin = mustLoad()

func All() map[string]string {
	out := make(map[string]string, len(builtin))
	maps.Copy(out, builtin)
	return out
}

func mustLoad() map[string]string {
	const dir = "assets"

	entries, err := charsetFiles.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	charsets := make(map[string]string, len(entries))
	for _, e := range entries {
		data, err := charsetFiles.ReadFile(dir + "/" + e.Name())
		if err != nil {
			panic(err)
		}
		charsets[strings.TrimSuffix(e.Name(), ".txt")] = string(data)
	}
	return charsets
}