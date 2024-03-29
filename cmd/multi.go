package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/martinvks/framer/client"
	"github.com/martinvks/framer/utils"
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
		"directory containing json request files (required)  https://github.com/martinvks/framer#json-request-files",
	)

	_ = multiCmd.MarkFlagRequired("directory")

	rootCmd.AddCommand(multiCmd)
}

var multiCmd = &cobra.Command{
	Use:     "multi [flags] target",
	Short:   "Send multiple requests to the target URL and print the response status code or error to console",
	Example: "framer multi -d ./requests https://martinvks.no",
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
	requestFiles, err := utils.GetRequestFiles(multiArgs.directory)
	if err != nil {
		return err
	}

	keyLogWriter, err := utils.GetKeyLogWriter(commonArgs.KeyLogFile)
	if err != nil {
		return fmt.Errorf("error creating key log writer: %w", err)
	}

	ip, err := utils.LookUp(commonArgs.Target.Hostname())
	if err != nil {
		return err
	}

	tableHeaders := getMultiTableHeaders(commonArgs.AddIdHeader)
	var tableData [][]string
	for _, requestFile := range requestFiles {
		id := uuid.NewString()

		if multiArgs.delay != 0 {
			time.Sleep(multiArgs.delay)
		}

		request := utils.GetRequest(
			id,
			multiArgs.addIdQuery,
			requestFile.RequestData,
			commonArgs,
		)

		response, err := client.DoRequest(
			commonArgs.Proto,
			commonArgs.Target,
			commonArgs.Timeout,
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
			commonArgs.AddIdHeader,
			id,
			requestFile.FileName,
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
	requestFilename string,
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
		requestFilename,
		responseCode,
		responseBodyLength,
		error,
	)
}
