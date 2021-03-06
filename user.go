// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/internal/resourcenames"
)

type User struct {
	Id                      int                    `json:"id,omitempty"`
	Login                   string                 `json:"login,omitempty"`
	Role                    int                    `json:"role,omitempty"`
	Status                  int                    `json:"status,omitempty"`
	LastLogin               ISO8601Time            `json:"lastLogin,omitempty"`
	Data                    map[string]interface{} `json:"data,omitempty"`
	IntroReviewed           bool                   `json:"introReviewed,omitempty"`
	AllDeviceTypesAvailable bool                   `json:"allDeviceTypesAvailable,omitempty"`
	client                  *Client
}

func (u *User) Save() *Error {
	_, err := u.client.request(resourcenames.UpdateUser, map[string]interface{}{
		"userId": u.Id,
		"user":   u,
	})

	return err
}

func (u *User) Remove() *Error {
	_, err := u.client.request(resourcenames.DeleteUser, map[string]interface{}{
		"userId": u.Id,
	})

	return err
}

func (u *User) UpdatePassword(password string) *Error {
	_, err := u.client.request(resourcenames.UpdateUser, map[string]interface{}{
		"userId": u.Id,
		"user": map[string]interface{}{
			"password": password,
		},
	})

	return err
}

func (u *User) AssignNetwork(networkId int) *Error {
	_, err := u.client.request(resourcenames.AssignNetwork, map[string]interface{}{
		"userId":    u.Id,
		"networkId": networkId,
	})

	return err
}

func (u *User) UnassignNetwork(networkId int) *Error {
	_, err := u.client.request(resourcenames.UnassignNetwork, map[string]interface{}{
		"userId":    u.Id,
		"networkId": networkId,
	})

	return err
}

func (u *User) AssignDeviceType(deviceTypeId int) *Error {
	_, err := u.client.request(resourcenames.AssignDeviceType, map[string]interface{}{
		"userId":       u.Id,
		"deviceTypeId": deviceTypeId,
	})

	return err
}

func (u *User) UnassignDeviceType(deviceTypeId int) *Error {
	_, err := u.client.request(resourcenames.UnassignDeviceType, map[string]interface{}{
		"userId":       u.Id,
		"deviceTypeId": deviceTypeId,
	})

	return err
}

func (u *User) ListNetworks() (list []*Network, err *Error) {
	rawRes, err := u.client.request(resourcenames.GetUser, map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return nil, err
	}

	pErr := json.Unmarshal(rawRes, &struct {
		List *[]*Network `json:"networks"`
	}{&list})
	if pErr != nil {
		return nil, newJSONErr(pErr)
	}
	for _, v := range list {
		v.client = u.client
	}
	return list, nil
}

func (u *User) ListDeviceTypes() (list []*DeviceType, err *Error) {
	rawRes, err := u.client.request(resourcenames.GetUserDeviceTypes, map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return nil, err
	}

	pErr := json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr(pErr)
	}
	for _, v := range list {
		v.client = u.client
	}

	return list, nil
}

func (u *User) AllowAllDeviceTypes() *Error {
	_, err := u.client.request(resourcenames.AllowAllDeviceTypes, map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return err
	}

	u.AllDeviceTypesAvailable = true
	return nil
}

func (u *User) DisallowAllDeviceTypes() *Error {
	_, err := u.client.request(resourcenames.DisallowAllDeviceTypes, map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return err
	}

	u.AllDeviceTypesAvailable = false
	return nil
}
