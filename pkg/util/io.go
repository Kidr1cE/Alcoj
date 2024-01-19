package util

import (
	"alcoj/pkg/docker"
	"log"
	"os"
	"path/filepath"
)

func Write(content []byte, path string) error {
	// if file not exist, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			log.Println("create dockerfile error: ", err)
			return err
		}
		file.Close()
	}

	// write content to dockerfile
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		log.Println("open dockerfile error: ", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		log.Println("write dockerfile error: ", err)
		return err
	}
	return nil
}

func WriteToAppFolder(filename string, content []byte) error {
	AppFolderPath := docker.AppFolderPath
	if len(content) == 0 {
		log.Println("content is empty")
		return nil
	}

	file, err := os.Create(filepath.Join(AppFolderPath, filename))
	if err != nil {
		log.Println("create file error: ", err)
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		log.Println("write file error: ", err)
		return err
	}

	return nil
}
