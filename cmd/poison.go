package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/martinvks/framer/client"
	"github.com/martinvks/framer/utils"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type poisonArguments struct {
	delay             time.Duration
	retryNonCacheable bool
	directory         string
}

var poisonArgs poisonArguments

func init() {
	poisonCmd.Flags().DurationVar(
		&poisonArgs.delay,
		"delay",
		0,
		"duration to wait between testing each request file",
	)

	poisonCmd.Flags().BoolVar(
		&poisonArgs.retryNonCacheable,
		"retry-non-cacheable",
		false,
		"send retry request to check cache poisoning for non cacheable response codes (e.g. 400 Bad Request)",
	)

	poisonCmd.Flags().StringVarP(
		&poisonArgs.directory,
		"directory",
		"d",
		"",
		"directory containing json request files (required)  https://github.com/martinvks/framer#json-request-files",
	)

	_ = poisonCmd.MarkFlagRequired("directory")

	rootCmd.AddCommand(poisonCmd)
}

var poisonCmd = &cobra.Command{
	Use:     "poison [flags] target",
	Short:   "Send multiple requests to the target and check for cache poisoning",
	Example: "framer poison -d ./requests https://martinvks.no/index.js",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runPoisonCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var cacheableStatusCodes = []string{
	"200",
	"203",
	"204",
	"206",
	"300",
	"301",
	"404",
	"405",
	"410",
	"414",
	"501",
}

var tableHeaders = []string{
	"FILE",
	"STATUS",
	"LENGTH",
	"RETRY",
	"POISONED",
	"ERROR",
}

type ResponseData struct {
	errorString string
	status      string
	length      string
	location    string
}

func runPoisonCmd() error {
	requestFiles, err := utils.GetRequestFiles(poisonArgs.directory)
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
	for _, requestFile := range requestFiles {
		id := uuid.NewString()

		if poisonArgs.delay != 0 {
			time.Sleep(poisonArgs.delay)
		}

		response, err := doRequest(
			id,
			ip,
			keyLogWriter,
			requestFile.RequestData,
		)

		tableData = append(tableData, []string{
			requestFile.FileName,
			response.status,
			response.length,
			"",
			"",
			response.errorString,
		})

		isDifferent := response != baseResponse && err == nil
		isCacheable := slices.Contains(cacheableStatusCodes, response.status)

		if isDifferent && (isCacheable || poisonArgs.retryNonCacheable) {
			response, err = doRequest(
				id,
				ip,
				keyLogWriter,
				utils.RequestData{AddDefaultHeaders: true},
			)

			poisoned := response != baseResponse && err == nil

			tableData = append(tableData, []string{
				requestFile.FileName,
				response.status,
				response.length,
				"true",
				strconv.FormatBool(poisoned),
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
		requestData,
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

	errorString := ""
	if err != nil {
		errorString = err.Error()
	}

	length := ""
	if response != nil {
		length = strconv.Itoa(len(response.Body))
	}

	location := ""
	if response != nil {
		if val, ok := response.Headers.Get("location"); ok {
			location = val
		}
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
		location:    location,
		status:      status,
	}, err
}
