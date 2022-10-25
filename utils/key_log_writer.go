package utils

import (
	"io"
	"os"
)

func GetKeyLogWriter(filename string) (io.Writer, error) {
	if filename != "" {
		return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	}
	return nil, nil
}
