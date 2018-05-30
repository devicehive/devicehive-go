package devicehive_go

import "time"

const (
	timestampLayout     = "2006-01-02T15:04:05.000"
	Timeout             = 5 * time.Second
	UserStatusActive    = 0
	UserStatusLocked    = 1
	UserStatusDisabled  = 2
	UserRoleAdmin       = 0
	UserRoleClient      = 1
	InvalidResponseErr  = "invalid response"
	InvalidRequestErr   = "invalid request"
	ServiceErr          = "service error"
	ConnectionFailedErr = "connection failed"
	DefaultPollingWaitTimeoutSeconds = 30
	WrongURLErr         = "wrong url"
)
