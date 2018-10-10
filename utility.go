package main

import (
	"os"
	"path/filepath"
)

func getSelfDir() (dirAbsPath string, err error) {
	self, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Dir(self))
}
