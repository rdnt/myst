package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/otiai10/copy"
	"mvdan.cc/sh/shell"

	"myst/pkg/builder/runtime"
)

func Step(name string, f func()) {
	fmt.Printf("%s...\n", name)

	f()

	fmt.Printf("%s... done\n", name)
}

func Run(command ...string) {
	joined := strings.Join(command, "\\\n")

	expanded, err := shell.Expand(joined, func(s string) string {
		return os.Getenv(s)
	})
	if err != nil {
		runtime.Panic(fmt.Errorf("failed to expand command: %s\n", err))
	}

	args, err := shellquote.Split(expanded)
	if err != nil {
		runtime.Panic(fmt.Errorf("failed to parse command: %s\n", err))
	}

	cmd, err := NewCmd("test", args...)
	if err != nil {
		runtime.Panic(fmt.Errorf("failed to create command: %s\n", err))
	}

	err = cmd.Start()
	if err != nil {
		runtime.Panic(fmt.Errorf("failed to start command: %s\n", err))
	}

	go func() {
		s := bufio.NewScanner(cmd.Stdout())
		s.Split(bufio.ScanBytes)

		for s.Scan() {
			str := s.Text()
			fmt.Print(str)
		}
	}()

	go func() {
		s := bufio.NewScanner(cmd.Stderr())
		s.Split(bufio.ScanBytes)

		for s.Scan() {
			str := s.Text()
			fmt.Print(str)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		runtime.Panic(fmt.Errorf("failed to run command: %s\n", err))
	}
}

func Env(s string) string {
	return os.Getenv(s)
}

func SetEnv(s, v string) {
	err := os.Setenv(s, v)
	if err != nil {
		runtime.Panic(err)
	}
}

var currentDir string

func SetCurrentDir(dir string) {
	currentDir = dir
}

func CopyDir(src, dst string) {
	err := copy.Copy(src, dst)
	if err != nil {
		runtime.Panic(err)
	}
}

func Touch(name string) {
	f, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		runtime.Panic(err)
	}

	err = f.Close()
	if err != nil {
		runtime.Panic(err)
	}
}

func CommandExists(binary string) {
	_, err := exec.LookPath(binary)
	if errors.Is(err, exec.ErrNotFound) {
		runtime.Panic(fmt.Errorf("%s: command not found", binary))
	} else if err != nil {
		runtime.Panic(err)
	}
}

func CleanDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		runtime.Panic(fmt.Errorf("failed to remove all %s: %s", dir, err))
	}
}
