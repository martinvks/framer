package utils

import (
	"bytes"
	"fmt"
	"github.com/Martinvks/httptestrunner/types"
)

func PrintHttpMessage(printLines int, response *types.HttpResponse) {
	for _, h := range response.Headers {
		fmt.Printf("%s: %s\n", h.Name, h.Value)
	}
	fmt.Println()
	lines := bytes.Split(response.Body, []byte{'\n'})
	for i, l := range lines {
		if printLines < 0 || i < printLines {
			fmt.Println(string(l))
		} else {
			break
		}
	}
}
