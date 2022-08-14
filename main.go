package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Martinvks/httptestrunner/arguments"
	"github.com/Martinvks/httptestrunner/http2"
	"github.com/Martinvks/httptestrunner/utils"
)

func main() {
	args, err := arguments.GetArguments()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	testCase, err := utils.GetTestCase(args.RequestsDirectory)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading request file: %v\n", err)
		os.Exit(1)
	}

	request := utils.GetRequest(args.Target, testCase)

	response, err := doRequest(args, request)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printResponse(args.PrintLines, response)
}

func doRequest(args arguments.Arguments, request utils.HTTPMessage) (utils.HTTPMessage, error) {
	switch args.Proto {
	case arguments.H2:
		return http2.SendHTTP2Request(args.Target, args.Timeout, args.KeyLogFile, request)
	case arguments.H3:
		fmt.Println("TODO: implement http3")
	}
	return utils.HTTPMessage{}, nil
}

func printResponse(printLines int, response utils.HTTPMessage) {
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
