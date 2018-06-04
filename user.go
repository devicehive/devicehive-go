package devicehive_go

import (
	"encoding/json"
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
	client 					*Client
}

func (u *User) Save() *Error {
	_, err := u.client.request("updateUser", map[string]interface{}{
		"userId": u.Id,
		"user":   u,
	})

	return err
}

func (u *User) Remove() *Error {
	_, err := u.client.request("deleteUser", map[string]interface{}{
		"userId": u.Id,
	})

	return err
}

func (u *User) UpdatePassword(password string) *Error {
	_, err := u.client.request("updateUser", map[string]interface{}{
		"userId": u.Id,
		"user": map[string]interface{}{
			"password": password,
		},
	})

	return err
}

func (u *User) AssignNetwork(networkId int) *Error {
	_, err := u.client.request("assignNetwork", map[string]interface{}{
		"userId":    u.Id,
		"networkId": networkId,
	})

	return err
}

func (u *User) UnassignNetwork(networkId int) *Error {
	_, err := u.client.request("unassignNetwork", map[string]interface{}{
		"userId":    u.Id,
		"networkId": networkId,
	})

	return err
}

func (u *User) AssignDeviceType(deviceTypeId int) *Error {
	_, err := u.client.request("assignDeviceType", map[string]interface{}{
		"userId":       u.Id,
		"deviceTypeId": deviceTypeId,
	})

	return err
}

func (u *User) UnassignDeviceType(deviceTypeId int) *Error {
	_, err := u.client.request("unassignDeviceType", map[string]interface{}{
		"userId":       u.Id,
		"deviceTypeId": deviceTypeId,
	})

	return err
}

func (u *User) ListNetworks() (list []*Network, err *Error) {
	rawRes, err := u.client.request("getUser", map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return nil, err
	}

	pErr := json.Unmarshal(rawRes, &struct {
		List *[]*Network `json:"networks"`
	}{&list})
	if pErr != nil {
		return nil, newJSONErr()
	}
	for _, v := range list {
		v.client = client
	}
	return list, nil
}

func (u *User) ListDeviceTypes() (list []*DeviceType, err *Error) {
	rawRes, err := u.client.request("getUserDeviceTypes", map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return nil, err
	}

	pErr := json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}
	for _, v := range list {
		v.client = u.client
	}

	return list, nil
}

func (u *User) AllowAllDeviceTypes() *Error {
	_, err := u.client.request("allowAllDeviceTypes", map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return err
	}

	u.AllDeviceTypesAvailable = true
	return nil
}

func (u *User) DisallowAllDeviceTypes() *Error {
	_, err := u.client.request("disallowAllDeviceTypes", map[string]interface{}{
		"userId": u.Id,
	})
	if err != nil {
		return err
	}

	u.AllDeviceTypesAvailable = false
	return nil
}
