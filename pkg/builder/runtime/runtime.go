package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

type panicCalled struct{}

func Panic(v any) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	file = filepath.Clean(filepath.ToSlash(file))
	wd = filepath.Clean(filepath.ToSlash(wd))

	file = strings.Replace(file, wd, "", 1)
	file = filepath.Clean(file)

	fmt.Println(fmt.Sprintf("%s:%d: %s", file, line, v))

	panic(panicCalled{})
}

var deferFunc func()

func Defer(f func()) {
	deferFunc = f
}

func Recover() {
	if err := recover(); err != nil {
		if _, ok := err.(panicCalled); !ok {
			fmt.Printf("\npanic recovered:\n\n%s\n", string(debug.Stack()))
		}
	}

	if deferFunc != nil {
		defer Recover()
		f := deferFunc
		deferFunc = nil
		f()
	}
}
