package dh

import (
	"encoding/json"
)

const (
	UserStatusActive   = 0
	UserStatusLocked   = 1
	UserStatusDisabled = 2
	UserRoleAdmin      = 0
	UserRoleClient     = 1
)

type User struct {
	client                  *Client
	Id                      int64                  `json:"id,omitempty"`
	Login                   string                 `json:"login,omitempty"`
	Role                    int                    `json:"role,omitempty"`
	Status                  int                    `json:"status,omitempty"`
	LastLogin               ISO8601Time            `json:"lastLogin,omitempty"`
	Data                    map[string]interface{} `json:"data,omitempty"`
	IntroReviewed           bool                   `json:"introReviewed,omitempty"`
	AllDeviceTypesAvailable bool                   `json:"allDeviceTypesAvailable,omitempty"`
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

func (u *User) AssignNetwork(networkId int64) *Error {
	_, err := u.client.request("assignNetwork", map[string]interface{}{
		"userId":    u.Id,
		"networkId": networkId,
	})

	return err
}

func (u *User) UnassignNetwork(networkId int64) *Error {
	_, err := u.client.request("unassignNetwork", map[string]interface{}{
		"userId":    u.Id,
		"networkId": networkId,
	})

	return err
}

func (u *User) AssignDeviceType(deviceTypeId int64) *Error {
	_, err := u.client.request("assignDeviceType", map[string]interface{}{
		"userId":       u.Id,
		"deviceTypeId": deviceTypeId,
	})

	return err
}

func (u *User) UnassignDeviceType(deviceTypeId int64) *Error {
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

func (c *Client) CreateUser(login, password string, role int, data map[string]interface{}, allDevTypesAvail bool) (user *User, err *Error) {
	user = &User{
		client: c,
		Login:  login,
		Role:   role,
		Data:   data,
		AllDeviceTypesAvailable: allDevTypesAvail,
	}

	res, err := c.request("createUser", map[string]interface{}{
		"user": map[string]interface{}{
			"login":    login,
			"role":     role,
			"status":   UserStatusActive,
			"password": password,
			"data":     data,
			"allDeviceTypesAvailable": allDevTypesAvail,
		},
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, user)
	if jsonErr != nil {
		return nil, newJSONErr()
	}

	return user, nil
}

func (c *Client) GetUser(userId int64) (user *User, err *Error) {
	rawRes, err := c.request("getUser", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		return nil, err
	}

	user = &User{
		client: c,
	}
	parseErr := json.Unmarshal(rawRes, user)
	if parseErr != nil {
		return nil, newJSONErr()
	}

	return user, nil
}

func (c *Client) GetCurrentUser() (user *User, err *Error) {
	rawRes, err := c.request("getCurrentUser", nil)
	if err != nil {
		return nil, err
	}

	user = &User{
		client: c,
	}
	parseErr := json.Unmarshal(rawRes, user)
	if parseErr != nil {
		return nil, newJSONErr()
	}

	return user, nil
}

func (c *Client) ListUsers(params *ListParams) (list []*User, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request("listUsers", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}
