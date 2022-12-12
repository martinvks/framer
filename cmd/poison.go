package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Martinvks/httptestrunner/client"
	"github.com/Martinvks/httptestrunner/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type poisonArguments struct {
	delay     time.Duration
	logFile   string
	directory string
}

var poisonArgs poisonArguments

func init() {
	poisonCmd.Flags().DurationVar(
		&poisonArgs.delay,
		"delay",
		0,
		"duration to wait between testing each request file",
	)

	poisonCmd.Flags().StringVarP(
		&poisonArgs.directory,
		"directory",
		"d",
		"",
		"directory containing json request files (required)  https://github.com/Martinvks/httptestrunner#json-request-files",
	)

	_ = poisonCmd.MarkFlagRequired("directory")

	rootCmd.AddCommand(poisonCmd)
}

var poisonCmd = &cobra.Command{
	Use:     "poison [flags] target",
	Short:   "Send multiple requests to the target and check for cache poisoning",
	Example: "httptestrunner poison -d ./requests https://martinvks.no/index.js",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runPoisonCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var tableHeaders = []string{
	"FILE",
	"STATUS",
	"LENGTH",
	"RETRY",
	"POISONED",
	"URL",
	"ERROR",
}

type ResponseData struct {
	errorString string
	status      string
	length      string
}

func runPoisonCmd() error {
	testCases, err := utils.GetAllTestCases(poisonArgs.directory)
	if err != nil {
		return fmt.Errorf("error reading request files: %w", err)
	}

	keyLogWriter, err := utils.GetKeyLogWriter(commonArgs.keyLogFile)
	if err != nil {
		return fmt.Errorf("error creating key log writer: %w", err)
	}

	ip, err := utils.LookUp(commonArgs.target.Hostname())
	if err != nil {
		return err
	}

	baseResponse, err := doRequest(
		uuid.NewString(),
		ip,
		keyLogWriter,
		utils.RequestData{AddDefaultHeaders: true},
	)

	if err != nil {
		return err
	}

	var validStatus = regexp.MustCompile(`^[23]\d{2}$`)

	if !validStatus.MatchString(baseResponse.status) {
		return fmt.Errorf("expected status code in the range 200-399 from target resource, but received: %s", baseResponse.status)
	}

	var tableData [][]string
	for _, testCase := range testCases {
		id := uuid.NewString()

		if poisonArgs.delay != 0 {
			time.Sleep(poisonArgs.delay)
		}

		response, err := doRequest(
			id,
			ip,
			keyLogWriter,
			testCase.RequestData,
		)

		tableData = append(tableData, []string{
			testCase.FileName,
			response.status,
			response.length,
			"",
			"",
			"",
			response.errorString,
		})

		if response != baseResponse && err == nil {
			response, err = doRequest(
				id,
				ip,
				keyLogWriter,
				utils.RequestData{AddDefaultHeaders: true},
			)

			poisoned := response != baseResponse && err == nil

			url := ""
			if poisoned {
				url = commonArgs.target.String()
			}

			tableData = append(tableData, []string{
				testCase.FileName,
				response.status,
				response.length,
				"true",
				strconv.FormatBool(poisoned),
				url,
				response.errorString,
			})
		}
	}

	err = utils.WriteTable(tableHeaders, tableData)
	if err != nil {
		return fmt.Errorf("error writing result table: %w", err)
	}

	return nil
}

func doRequest(
	id string,
	ip net.IP,
	keyLogWriter io.Writer,
	requestData utils.RequestData,
) (ResponseData, error) {
	request := utils.GetRequest(
		id,
		true,
		commonArgs.addIdHeader,
		commonArgs.proto,
		commonArgs.target,
		requestData,
	)

	response, err := client.DoRequest(
		commonArgs.proto,
		commonArgs.target,
		commonArgs.timeout,
		keyLogWriter,
		ip,
		&request,
	)

	errorString := ""
	if err != nil {
		errorString = err.Error()
	}

	length := ""
	if response != nil {
		length = strconv.Itoa(len(response.Body))
	}

	status := ""
	if response != nil {
		if val, ok := response.Headers.Get(":status"); ok {
			status = val
		}
	}

	return ResponseData{
		errorString: errorString,
		length:      length,
		status:      status,
	}, err
}