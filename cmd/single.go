package cmd

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/martinvks/framer/client"
	"github.com/martinvks/framer/utils"
	"github.com/spf13/cobra"
)

type singleArguments struct {
	addIdQuery bool
	printLines int
	fileName   string
}

var singleArgs singleArguments

func init() {
	singleCmd.Flags().BoolVar(
		&singleArgs.addIdQuery,
		"id-query",
		false,
		"add a query parameter with name \"id\" and a uuid v4 value to avoid cached responses",
	)

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
		"json request file (required)  https://github.com/martinvks/framer#json-request-files",
	)

	_ = singleCmd.MarkFlagRequired("filename")

	rootCmd.AddCommand(singleCmd)
}

var singleCmd = &cobra.Command{
	Use:     "single [flags] target",
	Short:   "Send a single request to the target URL and print the response to console",
	Example: "framer single -f ./request.json https://martinvks.no",
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
	id := uuid.NewString()

	requestFile, err := utils.GetRequestFile(singleArgs.fileName)
	if err != nil {
		return err
	}

	request := utils.GetRequest(
		id,
		singleArgs.addIdQuery,
		requestFile.RequestData,
		commonArgs,
	)

	keyLogWriter, err := utils.GetKeyLogWriter(commonArgs.KeyLogFile)
	if err != nil {
		return fmt.Errorf("error creating key log writer: %w", err)
	}

	ip, err := utils.LookUp(commonArgs.Target.Hostname())
	if err != nil {
		return err
	}

	response, err := client.DoRequest(
		commonArgs.Proto,
		commonArgs.Target,
		commonArgs.Timeout,
		keyLogWriter,
		ip,
		&request,
	)
	if err != nil {
		return err
	}

	utils.WriteResponse(singleArgs.printLines, response)

	return nil
}
