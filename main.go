package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/Martinvks/httptestrunner/arguments"
	"github.com/Martinvks/httptestrunner/http2"
	"github.com/Martinvks/httptestrunner/http3"
	"github.com/Martinvks/httptestrunner/types"
	"github.com/Martinvks/httptestrunner/utils"
)

func main() {
	args, err := arguments.GetArguments(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	switch args := args.(type) {
	case arguments.SingleModeArguments:
		err = runSingleMode(args)
	case arguments.MultiModeArguments:
		err = runMultiMode(args)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func runSingleMode(args arguments.SingleModeArguments) error {
	testCase, err := utils.GetSingleTestCase(args.FileName)
	if err != nil {
		return fmt.Errorf("error reading request file: %w", err)
	}

	request := utils.GetRequest(args.Target, testCase)

	keyLogWriter, err := getKeyLogWriter(args.KeyLogFile)
	if err != nil {
		return fmt.Errorf("error creating key log writer: %w", err)
	}

	ip, err := utils.LookUp(args.Target.Hostname())
	if err != nil {
		return err
	}

	response, err := doRequest(
		args.Proto,
		args.Target,
		args.Timeout,
		keyLogWriter,
		ip,
		&request,
	)
	if err != nil {
		return err
	}

	utils.PrintHttpMessage(args.PrintLines, response)

	return nil
}

func runMultiMode(args arguments.MultiModeArguments) error {
	testCases, err := utils.GetAllTestCases(args.Directory)
	if err != nil {
		return fmt.Errorf("error reading request files: %w", err)
	}

	keyLogWriter, err := getKeyLogWriter(args.KeyLogFile)
	if err != nil {
		return fmt.Errorf("error creating key log writer: %w", err)
	}

	ip, err := utils.LookUp(args.Target.Hostname())
	if err != nil {
		return err
	}

	csvWriter, err := utils.GetCsvWriter(args.CsvLogFile)
	if err != nil {
		return fmt.Errorf("error creating csv writer: %w", err)
	}
	defer csvWriter.Close()

	err = csvWriter.WriteHeaders()
	if err != nil {
		return fmt.Errorf("error writing csv headers: %w", err)
	}

	for _, testCase := range testCases {
		request := utils.GetRequest(args.Target, testCase)

		response, err := doRequest(
			args.Proto,
			args.Target,
			args.Timeout,
			keyLogWriter,
			ip,
			&request,
		)

		requestError := ""
		if err != nil {
			requestError = err.Error()
		}

		responseCode := ""
		if response != nil {
			if val, ok := response.Headers.Get(":status"); ok {
				responseCode = val
			}
		}

		err = csvWriter.WriteData(
			testCase.FileName,
			testCase.Id,
			responseCode,
			requestError,
		)
		if err != nil {
			return fmt.Errorf("error writing csv record: %w", err)
		}
	}

	return nil
}

func doRequest(
	proto int,
	target *url.URL,
	timeout time.Duration,
	keyLogWriter io.Writer,
	ip net.IP,
	request *types.HttpRequest,
) (*types.HttpResponse, error) {
	switch proto {
	case arguments.H2:
		return http2.SendHTTP2Request(ip, target, timeout, keyLogWriter, request)
	case arguments.H3:
		return http3.SendHTTP3Request(ip, target, timeout, keyLogWriter, request)
	default:
		return nil, fmt.Errorf("unknown proto %d", proto)
	}
}

func getKeyLogWriter(filename string) (io.Writer, error) {
	if filename != "" {
		return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	}
	return nil, nil
}
