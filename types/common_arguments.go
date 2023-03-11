package types

import (
	"net/url"
	"time"
)

type CommonArguments struct {
	AddIdHeader   bool
	CommonHeaders Headers
	IdHeaderName  string
	KeyLogFile    string
	Proto         int
	Timeout       time.Duration
	Target        *url.URL
}
