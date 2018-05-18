package dh

import (
	"errors"
	"github.com/devicehive/devicehive-go/internal/utils"
	"time"
)

type ListParams struct {
	DeviceId          string    `json:"deviceId,omitempty"`
	Start             time.Time `json:"start,omitempty"`
	End               time.Time `json:"end,omitempty"`
	Notification      string    `json:"notification,omitempty"`
	Command           string    `json:"command,omitempty"`
	Status            string    `json:"status,omitempty"`
	SortField         string    `json:"sortField,omitempty"`
	SortOrder         string    `json:"sortOrder,omitempty"`
	Take              int       `json:"take,omitempty"`
	Skip              int       `json:"skip,omitempty"`
	DeviceName        string    `json:"name,omitempty"`
	DeviceNamePattern string    `json:"namePattern,omitempty"`
	NetworkId         string    `json:"networkId,omitempty"`
	NetworkName       string    `json:"networkName,omitempty"`
}

func (p *ListParams) Map() (m map[string]interface{}, err error) {
	m = utils.StructToJSONMap(p)

	if m == nil {
		return nil, errors.New("invalid JSON representation of struct")
	}

	if p.Start.Unix() < 0 {
		delete(m, "start")
	} else {
		m["start"] = p.Start.Format(timestampLayout)
	}
	if p.End.Unix() < 0 {
		delete(m, "end")
	} else {
		m["end"] = p.End.Format(timestampLayout)
	}

	return m, nil
}
