package transport

import "time"

type Params struct {
	ReconnectionTries    int
	ReconnectionInterval time.Duration
	DefaultTimeout       time.Duration
}
