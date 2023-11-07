package internal

import (
	"errors"
	"os"
)

func IsFileExists(file string) (bool, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	if stat.IsDir() {
		return false, errors.New("the path isn't a valid file")
	}
	return true, err
}
