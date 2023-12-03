package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("No caller information")
		return
	}

	fmt.Printf("Function file path: %s\n", file)
	fmt.Printf("Function directory: %s\n", filepath.Dir(file))
	exePath, err := os.Executable() // 获取当前执行文件的路径
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath) // 从路径中获取目录
	fmt.Println("Executable directory:", exeDir)
}
