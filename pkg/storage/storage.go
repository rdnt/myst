package storage

import (
	"bytes"
	"fmt"
	"os"

	"github.com/natefinch/atomic"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

func Init() error {
	err := os.MkdirAll("data/users", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll("data/keystorerepo", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Load(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

func Save(path string, b []byte) error {
	return atomic.WriteFile(path, bytes.NewReader(b))
}
