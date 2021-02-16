package database

import (
	"fmt"
	"io/ioutil"
	"vitess.io/vitess/go/ioutil2"
)

func GetKeystore(id string) ([]byte, error) {
	filename := fmt.Sprintf("data/keystores/%s.mst", id)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func SetKeystore(id string, data []byte) error {
	// atomic write is natively supported on UNIX, and on Windows if using NTFS
	filename := fmt.Sprintf("data/keystores/%s.mst", id)
	err := ioutil2.WriteFileAtomic(filename, data, 0664)
	if err != nil {
		return err
	}
	return nil
}
