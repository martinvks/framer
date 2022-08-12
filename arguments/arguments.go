package arguments

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"time"
)

const (
	H2 = iota
	H3
)

type Arguments struct {
	KeyLogFile        string
	PrintLines        int
	Proto             int
	Timeout           time.Duration
	Target            *url.URL
	RequestsDirectory string
}

func GetArguments() (Arguments, error) {
	var keyLogFile string
	flag.StringVar(&keyLogFile, "k", "", "Filename to log TLS master secrets")

	var printLines int
	flag.IntVar(&printLines, "l", 10, "Number of lines to print from the response body")

	var timeout time.Duration
	flag.DurationVar(&timeout, "t", 10*time.Second, "timeout")

	var requestsDirectory string
	flag.StringVar(&requestsDirectory, "d", "", "directory containing json request files")

	var proto string
	flag.StringVar(
		&proto,
		"p",
		"h2",
		"specifies which protocol to use. Must be one of \"h2\" or \"h3\"",
	)

	flag.Parse()
	args := flag.Args()

	if requestsDirectory == "" {
		return Arguments{}, errors.New("directory with json request files required")
	}

	var protoInt int
	switch proto {
	case "h2":
		protoInt = H2
	case "h3":
		protoInt = H3
	default:
		return Arguments{}, errors.New(
			fmt.Sprintf("unknown protocol %v", proto),
		)
	}

	if len(args) == 0 {
		return Arguments{}, errors.New("missing target URL")
	}

	target, err := url.Parse(args[0])

	if err != nil {
		return Arguments{}, err
	}

	return Arguments{
		KeyLogFile:        keyLogFile,
		PrintLines:        printLines,
		Proto:             protoInt,
		Timeout:           timeout,
		Target:            target,
		RequestsDirectory: requestsDirectory,
	}, nil
}
