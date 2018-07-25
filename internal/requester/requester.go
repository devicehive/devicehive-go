package requester

import (
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"time"
)

type Requester interface {
	Request(resourceName string, data map[string]interface{}, timeout time.Duration, accessToken string) ([]byte, error)
	PrepareRequestData(resourceName string, data map[string]interface{}, accessToken string) (resource string, reqParams *apirequests.RequestParams)
}
