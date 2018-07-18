package devicehive_go

import "time"

type ConnectionParams struct {
	ReconnectionTries    int
	ReconnectionInterval time.Duration
	RequestTimeout       time.Duration
}

func (p *ConnectionParams) Timeout() time.Duration {
	if p == nil || p.RequestTimeout == 0 {
		return Timeout
	}

	return p.RequestTimeout
}
