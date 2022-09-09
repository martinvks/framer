package utils

import (
	"bytes"
	"fmt"
)

func PrintHttpMessage(printLines int, message HTTPMessage) {
	for _, h := range message.Headers {
		fmt.Printf("%s: %s\n", h.Name, h.Value)
	}
	fmt.Println()
	lines := bytes.Split(message.Body, []byte{'\n'})
	for i, l := range lines {
		if printLines < 0 || i < printLines {
			fmt.Println(string(l))
		} else {
			break
		}
	}
}
