package dh

import (
	"errors"
	"github.com/devicehive/devicehive-go/internal/utils"
	"time"
)

type ListParams struct {
	DeviceId     string    `json:"deviceId,omitempty"`
	Start        time.Time `json:"start,omitempty"`
	End          time.Time `json:"end,omitempty"`
	Notification string    `json:"notification,omitempty"`
	Command      string    `json:"command,omitempty"`
	Status       string    `json:"status,omitempty"`
	SortField    string    `json:"sortField,omitempty"`
	SortOrder    string    `json:"sortOrder,omitempty"`
	Take         int       `json:"take,omitempty"`
	Skip         int       `json:"skip,omitempty"`
	Name         string    `json:"name,omitempty"`
	NamePattern  string    `json:"namePattern,omitempty"`
	Login        string    `json:"login,omitempty"`
	LoginPattern string    `json:"loginPattern,omitempty"`
	NetworkId    string    `json:"networkId,omitempty"`
	NetworkName  string    `json:"networkName,omitempty"`
	UserRole     int       `json:"role,omitempty"`
	UserStatus   int       `json:"status,omitempty"`
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
