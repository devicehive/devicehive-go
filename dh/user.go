package dh

import (
	"encoding/json"
)

const (
	UserStatusActive = 0
	UserStatusLocked = 1
	UserStatusDisabled = 2
	UserRoleAdmin = 0
	UserRoleClient = 1
)

type User struct {
	client *Client
	Id int64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Role int `json:"role,omitempty"`
	Status int `json:"status,omitempty"`
	LastLogin ISO8601Time `json:"lastLogin,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
	IntroReviewed bool `json:"introReviewed,omitempty"`
	AllDeviceTypesAvailable bool `json:"allDeviceTypesAvailable,omitempty"`
}

func (u *User) Remove() *Error {
	_, err := u.client.request("deleteUser", map[string]interface{}{
		"userId": u.Id,
	})

	return err
}

func (c *Client) CreateUser(login, password string, role int, data map[string]interface{}, allDevTypesAvail bool) (user *User, err *Error) {
	user = &User{
		client: c,
		Login: login,
		Role: role,
		Data: data,
		AllDeviceTypesAvailable: allDevTypesAvail,
	}

	res, err := c.request("createUser", map[string]interface{}{
		"user": map[string]interface{}{
			"login": login,
			"role": role,
			"status": UserStatusActive,
			"password": password,
			"data": data,
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
