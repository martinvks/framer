package cmd

import (
	"fmt"
	"os"

	"github.com/Martinvks/httptestrunner/client"
	"github.com/Martinvks/httptestrunner/utils"
	"github.com/spf13/cobra"
)

type multiArguments struct {
	logFile   string
	directory string
}

var multiArgs multiArguments

func init() {
	multiCmd.Flags().StringVar(
		&multiArgs.logFile,
		"logfile",
		"",
		"filename to log result in csv format. if not set, the result will be printed to console",
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
	Short:   "Send multiple requests to the target",
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

	csvWriter, err := utils.GetCsvWriter(multiArgs.logFile, commonArgs.addIdHeader)
	if err != nil {
		return fmt.Errorf("error creating csv writer: %w", err)
	}
	defer csvWriter.Close()

	err = csvWriter.WriteHeaders()
	if err != nil {
		return fmt.Errorf("error writing csv headers: %w", err)
	}

	for _, testCase := range testCases {
		request := utils.GetRequest(
			commonArgs.addIdHeader,
			commonArgs.addIdQuery,
			commonArgs.proto,
			commonArgs.target,
			testCase,
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
