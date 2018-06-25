package devicehive_go

import "time"

type ConnectionParams struct {
	ReconnectionTries    int
	ReconnectionInterval time.Duration
}
