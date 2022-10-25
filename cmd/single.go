package cmd

import (
	"fmt"
	"os"

	"github.com/Martinvks/httptestrunner/client"
	"github.com/Martinvks/httptestrunner/utils"
	"github.com/spf13/cobra"
)

type singleArguments struct {
	printLines int
	fileName   string
}

var singleArgs singleArguments

func init() {
	singleCmd.Flags().IntVarP(
		&singleArgs.printLines,
		"lines",
		"l",
		10,
		"number of lines to print from the response body",
	)

	singleCmd.Flags().StringVarP(
		&singleArgs.fileName,
		"filename",
		"f",
		"",
		"json request file (required)  https://github.com/Martinvks/httptestrunner#json-request-files",
	)

	_ = singleCmd.MarkFlagRequired("filename")

	rootCmd.AddCommand(singleCmd)
}

var singleCmd = &cobra.Command{
	Use:     "single [flags] target",
	Short:   "Send a single request to the target",
	Example: "httptestrunner single -f ./request.json https://martinvks.no",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runSingleCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func runSingleCmd() error {
	testCase, err := utils.GetSingleTestCase(singleArgs.fileName)
	if err != nil {
		return fmt.Errorf("error reading request file: %w", err)
	}

	request := utils.GetRequest(
		commonArgs.addIdHeader,
		commonArgs.addIdQuery,
		commonArgs.proto,
		commonArgs.target,
		testCase,
	)

	keyLogWriter, err := utils.GetKeyLogWriter(commonArgs.keyLogFile)
	if err != nil {
		return fmt.Errorf("error creating key log writer: %w", err)
	}

	ip, err := utils.LookUp(commonArgs.target.Hostname())
	if err != nil {
		return err
	}

	response, err := client.DoRequest(
		commonArgs.proto,
		commonArgs.target,
		commonArgs.timeout,
		keyLogWriter,
		ip,
		&request,
	)
	if err != nil {
		return err
	}

	utils.PrintHttpMessage(singleArgs.printLines, response)

	return nil
}
