package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/martinvks/framer/types"
	"github.com/spf13/cobra"
)

var (
	headers    []string
	proto      string
	commonArgs types.CommonArguments
)

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&commonArgs.AddIdHeader,
		"id-header",
		false,
		"add a header field with name \"x-id\" and a random uuid v4 value. the value will be added to the output when using the \"multi\" command",
	)

	rootCmd.PersistentFlags().StringVar(
		&commonArgs.IdHeaderName,
		"id-header-name",
		"x-id",
		"change the field name for the id-header. can be used as cache buster for vary headers, i.e., \"--id-header --id-header-name origin\"",
	)

	rootCmd.PersistentFlags().StringArrayVarP(
		&headers,
		"header",
		"H",
		[]string{},
		"common header fields added to each request. syntax similar to curl: -H \"x-extra-header: val\"",
	)

	rootCmd.PersistentFlags().StringVarP(
		&commonArgs.KeyLogFile,
		"keylogfile",
		"k",
		"",
		"filename to log TLS master secrets",
	)

	rootCmd.PersistentFlags().DurationVarP(
		&commonArgs.Timeout,
		"timeout",
		"t",
		10*time.Second,
		"timeout",
	)

	rootCmd.PersistentFlags().StringVarP(
		&proto,
		"protocol",
		"p",
		"h2",
		"specifies which protocol to use. Must be one of \"h2\" or \"h3\"",
	)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var rootCmd = &cobra.Command{
	Use:   "framer",
	Short: "An HTTP client for sending (possibly malformed) HTTP/2 and HTTP/3 requests",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		for _, header := range headers {
			name, value, found := strings.Cut(header, ":")
			if !found {
				return fmt.Errorf("invalid header '%s', expected syntax: 'x-extra-header: val'", header)
			}
			commonArgs.CommonHeaders = append(
				commonArgs.CommonHeaders,
				types.Header{
					Name:  strings.TrimSpace(strings.ToLower(name)),
					Value: strings.TrimSpace(value),
				})
		}

		switch proto {
		case "h2":
			commonArgs.Proto = types.H2
		case "h3":
			commonArgs.Proto = types.H3
		default:
			return fmt.Errorf("unknown protocol '%s'", proto)
		}

		target, err := url.Parse(args[0])
		if err != nil {
			return err
		}
		commonArgs.Target = target

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
