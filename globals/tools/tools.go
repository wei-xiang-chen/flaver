package tools

import (
	"os"
	"path"
	"strings"
)

func GetWorkDirectory() (string, error) {
	var directory string
	var err error
	if directory, err = os.Getwd(); err != nil {
		return "", err
	}
	for directory != "/" && !strings.HasSuffix(directory, "server") {
		directory = path.Dir(directory)
	}
	return directory, nil
}
