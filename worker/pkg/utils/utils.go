package utils

import (
	"os"
	"path/filepath"
)

func WriteDockerfile() {

}

func WriteCodeFile() {

}

func GetPath(folder string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := filepath.Join(wd, folder)
	return path, nil
}
