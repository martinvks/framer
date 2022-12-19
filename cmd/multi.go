package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Martinvks/httptestrunner/client"
	"github.com/Martinvks/httptestrunner/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type multiArguments struct {
	addIdQuery bool
	delay      time.Duration
	directory  string
}

var multiArgs multiArguments

func init() {
	multiCmd.Flags().BoolVar(
		&multiArgs.addIdQuery,
		"id-query",
		false,
		"add a query parameter with name \"id\" and a uuid v4 value to avoid cached responses",
	)

	multiCmd.Flags().DurationVar(
		&multiArgs.delay,
		"delay",
		0,
		"duration to wait between each request",
	)

	multiCmd.Flags().StringVarP(
		&multiArgs.directory,
		"directory",
		"d",
		"",
		"directory containing json request files (required)  https://github.com/Martinvks/httptestrunner#json-request-files",
	)

	_ = multiCmd.MarkFlagRequired("directory")

	rootCmd.AddCommand(multiCmd)
}

var multiCmd = &cobra.Command{
	Use:     "multi [flags] target",
	Short:   "Send multiple requests to the target URL and print the response status code or error to console",
	Example: "httptestrunner multi -d ./requests https://martinvks.no",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runMultiCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func runMultiCmd() error {
	testCases, err := utils.GetAllTestCases(multiArgs.directory)
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

	tableHeaders := getMultiTableHeaders(commonArgs.addIdHeader)
	var tableData [][]string
	for _, testCase := range testCases {
		id := uuid.NewString()

		if multiArgs.delay != 0 {
			time.Sleep(multiArgs.delay)
		}

		request := utils.GetRequest(
			id,
			multiArgs.addIdQuery,
			commonArgs.addIdHeader,
			commonArgs.proto,
			commonArgs.target,
			commonArgs.commonHeaders,
			testCase.RequestData,
		)

		response, err := client.DoRequest(
			commonArgs.proto,
			commonArgs.target,
			commonArgs.timeout,
			keyLogWriter,
			ip,
			&request,
		)

		requestError := ""
		if err != nil {
			requestError = err.Error()
		}

		responseBodyLength := ""
		if response != nil {
			responseBodyLength = strconv.Itoa(len(response.Body))
		}

		responseCode := ""
		if response != nil {
			if val, ok := response.Headers.Get(":status"); ok {
				responseCode = val
			}
		}

		tableData = append(tableData, getMultiTableData(
			commonArgs.addIdHeader,
			id,
			testCase.FileName,
			responseCode,
			responseBodyLength,
			requestError,
		))
	}

	err = utils.WriteTable(tableHeaders, tableData)
	if err != nil {
		return fmt.Errorf("error writing result table: %w", err)
	}

	return nil
}

func getMultiTableHeaders(addIdField bool) []string {
	var headers []string
	if addIdField {
		headers = append(headers, "ID")
	}

	return append(
		headers,
		"FILE",
		"STATUS",
		"LENGTH",
		"ERROR",
	)
}

func getMultiTableData(
	addIdField bool,
	requestId string,
	testFilename string,
	responseCode string,
	responseBodyLength string,
	error string,
) []string {
	var row []string
	if addIdField {
		row = append(row, requestId)
	}

	return append(
		row,
		testFilename,
		responseCode,
		responseBodyLength,
		error,
	)
}
