package main

import (
	"errors"
	"fmt"
	"github.com/Martinvks/httptestrunner/http2"
	"github.com/Martinvks/httptestrunner/utils"
	"net/url"
	"os"
	"time"

	"github.com/Martinvks/httptestrunner/arguments"
)

func main() {
	args, err := arguments.GetArguments(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	switch args := args.(type) {
	case arguments.SingleModeArguments:
		err = runSingleMode(args)
	case arguments.MultipleModeArguments:
		err = runMultipleMode(args)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func runSingleMode(args arguments.SingleModeArguments) error {
	testCase, err := utils.GetSingleTestCase(args.FileName)
	if err != nil {
		return fmt.Errorf("error reading request file: %v\n", err)
	}

	request := utils.GetRequest(args.Target, testCase)

	response, err := doRequest(
		args.Proto,
		args.Target,
		args.Timeout,
		args.KeyLogFile,
		request,
	)
	if err != nil {
		return err
	}

	utils.PrintHttpMessage(args.PrintLines, response)

	return nil
}

func runMultipleMode(args arguments.MultipleModeArguments) error {
	testCases, err := utils.GetAllTestCases(args.Directory)
	if err != nil {
		return fmt.Errorf("error reading request files: %v\n", err)
	}

	for _, testCase := range testCases {
		request := utils.GetRequest(args.Target, testCase)
		response, err := doRequest(
			args.Proto,
			args.Target,
			args.Timeout,
			args.KeyLogFile,
			request,
		)
		utils.WriteHttpMessage(response, err)
	}

	return nil
}

func doRequest(
	proto int,
	target *url.URL,
	timeout time.Duration,
	keyLogFile string,
	request utils.HTTPMessage,
) (utils.HTTPMessage, error) {
	switch proto {
	case arguments.H2:
		return http2.SendHTTP2Request(target, timeout, keyLogFile, request)
	case arguments.H3:
		return utils.HTTPMessage{}, errors.New("h3 not implemented")
	default:
		return utils.HTTPMessage{}, fmt.Errorf("unknown proto %d", proto)
	}
}
