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

	// 拼接当前工作目录的路径和app目录的相对路径
	path := filepath.Join(wd, folder)
	return path, nil
}
