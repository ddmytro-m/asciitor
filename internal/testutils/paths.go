package testutils

import (
	"path/filepath"
	"runtime"
)

func GetRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

func DataPath(relPath string) string {
	return filepath.Join(GetRoot(), "test/data", relPath)
}
