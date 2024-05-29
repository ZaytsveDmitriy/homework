package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		return
	}
	dirPath := os.Args[1]
	env, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
	}

	RunCmd(args[2:], env)
}
