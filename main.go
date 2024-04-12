package main

import (
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	rl, err := readline.New(PromptAssembly() + " $ ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		// 处理输入
		if err = execInput(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// ErrNoPath is returned when 'cd' was called without a second argument.
var ErrNoPath = errors.New("path required")

func execInput(input string) error {
	// 移除换行符
	input = strings.TrimSpace(input)

	// 分割输入，获取命令和参数
	args := strings.Fields(input)

	// 检查并去除每个参数两端的引号
	for i, arg := range args {
		args[i] = strings.Trim(arg, "\"'")
	}

	// Check for built-in commands.
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "cd":
		// 'cd' to home with empty path not yet supported.
		if len(args) < 2 {
			return ErrNoPath
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "history":
		return historyCommand()
	}

	// Prepare the command to execute.
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	appendToHistory(input)

	// Execute the command and return the error.
	return cmd.Run()
}

func appendToHistory(command string) {
	// Open the .aish_history file in append mode
	file, err := os.OpenFile(".aish_history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the command and its result to the file
	if _, err := file.WriteString(command + "\n"); err != nil {
		log.Fatal(err)
	}
}

func historyCommand() error {
	// Open the .aish_history file
	file, err := os.Open(".aish_history")
	if err != nil {
		return err
	}
	defer file.Close()

	// Read and print the file contents
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Print(string(content))
	return nil
}

func PromptAssembly() string {
	// 获取当前用户
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}

	// 获取主机名
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	// 获取当前路径
	path, err := os.Getwd()
	if err != nil {
		path = "."
	} else {
		// 简化路径：将用户目录替换为 '~'
		homeDir, err := os.UserHomeDir()
		if err == nil && strings.HasPrefix(path, homeDir) {
			path = "~" + strings.TrimPrefix(path, homeDir)
		}
	}

	return fmt.Sprintf("\033[1;34m%s\033[0m@\033[1;32m%s\033[0m:\033[1;36m%s\033[0m", user, host, path)
}
