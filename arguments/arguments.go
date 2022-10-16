package arguments

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"
)

const (
	H2 = iota
	H3
)

var (
	idQuery    bool
	keyLogFile string
	proto      string
	timeout    time.Duration

	printLines int
	fileName   string

	csvLogFile string
	directory  string
)

type CommonArguments struct {
	IdQuery    bool
	KeyLogFile string
	Proto      int
	Timeout    time.Duration
	Target     *url.URL
}

type SingleModeArguments struct {
	CommonArguments
	PrintLines int
	FileName   string
}

type MultiModeArguments struct {
	CommonArguments
	CsvLogFile string
	Directory  string
}

var (
	singleFlagSet = flag.NewFlagSet("single", flag.ExitOnError)
	multiFlagSet  = flag.NewFlagSet("multi", flag.ExitOnError)
)

var subcommands = map[string]*flag.FlagSet{
	singleFlagSet.Name(): singleFlagSet,
	multiFlagSet.Name():  multiFlagSet,
}

func Usage() {
	w := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(w, "Usage: httptestrunner <mode> <flags> <target>\n\n")
	_, _ = fmt.Fprintf(w, "Modes:\n")
	_, _ = fmt.Fprintf(w, "single\tsend a single request to the target\n")
	_, _ = fmt.Fprintf(w, "multi\tsend multiple requests to the target\n\n")
	_, _ = fmt.Fprintf(w, "Run 'httptestrunner <mode> -h' for more information\n")
}

func setupSubcommandUsage(fs *flag.FlagSet) {
	fs.Usage = func() {
		w := fs.Output()
		_, _ = fmt.Fprintf(w, "Usage: httptestrunner %s <flags> <target>\n\n", fs.Name())
		_, _ = fmt.Fprintf(w, "Flags:\n")
		fs.PrintDefaults()
		_, _ = fmt.Fprintf(w, "\nTarget:\n")
		_, _ = fmt.Fprintf(w, "The target URL. e.g. https://example.com\n\n")
	}
}

func setupCommonFlags() {
	for _, fs := range subcommands {
		fs.BoolVar(
			&idQuery,
			"id_query",
			false,
			"add query parameter with name \"id\" and a uuid v4 value to avoid cached responses",
		)

		fs.StringVar(
			&keyLogFile,
			"k",
			"",
			"filename to log TLS master secrets",
		)

		fs.StringVar(
			&proto,
			"p",
			"h2",
			"specifies which protocol to use. Must be one of \"h2\" or \"h3\"",
		)

		fs.DurationVar(
			&timeout,
			"t",
			10*time.Second,
			"timeout",
		)
	}
}

func setupSingleModeFlags() {
	singleFlagSet.IntVar(
		&printLines,
		"l",
		10,
		"number of lines to print from the response body",
	)

	singleFlagSet.StringVar(
		&fileName,
		"f",
		"",
		"filename with request data in json format",
	)
}

func setupMultiModeFlags() {
	multiFlagSet.StringVar(
		&csvLogFile,
		"csv",
		"",
		"filename to log result in csv format. if not set, the result will be printed to console",
	)

	multiFlagSet.StringVar(
		&directory,
		"d",
		"",
		"directory containing json request files",
	)
}

func GetArguments(osArgs []string) (interface{}, error) {
	if len(osArgs) < 1 {
		return nil, errors.New("you must select a mode")
	}

	setupCommonFlags()
	setupSingleModeFlags()
	setupMultiModeFlags()

	for _, fs := range subcommands {
		setupSubcommandUsage(fs)
	}

	flagSet := subcommands[osArgs[0]]
	if flagSet == nil {
		return nil, fmt.Errorf("unknown mode '%s'", os.Args[1])
	}

	err := flagSet.Parse(osArgs[1:])
	if err != nil {
		return nil, err
	}

	args := flagSet.Args()

	var intProto int
	switch proto {
	case "h2":
		intProto = H2
	case "h3":
		intProto = H3
	default:
		return nil, fmt.Errorf("unknown protocol '%s'", proto)
	}

	if len(args) == 0 {
		return nil, errors.New("missing target URL")
	}

	target, err := url.Parse(args[0])
	if err != nil {
		return nil, err
	}

	commonArguments := CommonArguments{
		IdQuery:    idQuery,
		KeyLogFile: keyLogFile,
		Proto:      intProto,
		Timeout:    timeout,
		Target:     target,
	}

	if flagSet.Name() == "single" {
		if fileName == "" {
			return nil, errors.New("filename required")
		}

		return SingleModeArguments{
			CommonArguments: commonArguments,
			PrintLines:      printLines,
			FileName:        fileName,
		}, nil
	}

	if flagSet.Name() == "multi" {
		if directory == "" {
			return nil, errors.New("directory required")
		}

		return MultiModeArguments{
			CommonArguments: commonArguments,
			CsvLogFile:      csvLogFile,
			Directory:       directory,
		}, nil
	}

	panic(fmt.Sprintf("unknown flag set name %s", flagSet.Name()))
}
