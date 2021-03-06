// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"time"

	"github.com/devicehive/devicehive-go/internal/resourcenames"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter"
)

// Main struct which serves as entry point to DeviceHive API
type Client struct {
	transportAdapter          transportadapter.TransportAdapter
	PollingWaitTimeoutSeconds int
	subscriptionTimestamp     time.Time
	defaultRequestTimeout     time.Duration
}

// Constructor, doesn't create device at DH
func (c *Client) NewDevice() *Device {
	return &Device{client: c}
}

// Constructor, doesn't create device type at DH
func (c *Client) NewDeviceType() *DeviceType {
	return &DeviceType{client: c}
}

// Constructor, doesn't create user at DH
func (c *Client) NewUser() *User {
	return &User{client: c}
}

// Constructor, doesn't create network at DH
func (c *Client) NewNetwork() *Network {
	return &Network{client: c}
}

// Constructor, doesn't create command at DH
func (c *Client) NewCommand() *Command {
	return &Command{client: c}
}

// Constructor, doesn't create notification at DH
func (c *Client) NewNotification() *Notification {
	return &Notification{}
}

// Subscribes to notifications by custom filter
// In case params is nil returns subscription for all notifications
func (c *Client) SubscribeNotifications(params *SubscribeParams) (*NotificationSubscription, *Error) {
	tspSubs, subsId, err := c.subscribe(resourcenames.SubscribeNotifications, params)
	if err != nil || tspSubs == nil {
		return nil, err
	}

	subs := newNotificationSubscription(subsId, tspSubs, c)

	return subs, nil
}

// Subscribes to commands by custom filter
// In case params is nil returns subscription for all commands
func (c *Client) SubscribeCommands(params *SubscribeParams) (*CommandSubscription, *Error) {
	tspSubs, subsId, err := c.subscribe(resourcenames.SubscribeCommands, params)
	if err != nil || tspSubs == nil {
		return nil, err
	}

	subs := newCommandSubscription(subsId, tspSubs, c)

	return subs, nil
}

func (c *Client) setCreds(login, password string) {
	c.transportAdapter.SetCreds(login, password)
}

func (c *Client) setRefreshToken(refTok string) {
	c.transportAdapter.SetRefreshToken(refTok)
}

func (c *Client) authenticate(token string) (bool, *Error) {
	result, rawErr := c.transportAdapter.Authenticate(token, c.defaultRequestTimeout)
	if rawErr != nil {
		return result, newError(rawErr)
	}

	return result, nil
}

func (c *Client) subscribe(resourceName string, params *SubscribeParams) (subscription *transport.Subscription, subscriptionId string, err *Error) {
	if params == nil {
		params = &SubscribeParams{}
	}

	params.WaitTimeout = c.PollingWaitTimeoutSeconds

	if params.Timestamp.Unix() <= 0 {
		params.Timestamp = c.subscriptionTimestamp
	}

	data, jsonErr := params.Map()
	if jsonErr != nil {
		return nil, "", &Error{name: InvalidRequestErr, reason: jsonErr.Error()}
	}

	subs, subscriptionId, rawErr := c.transportAdapter.Subscribe(resourceName, c.PollingWaitTimeoutSeconds, data)
	if rawErr != nil {
		return nil, "", newTransportErr(rawErr)
	}

	return subs, subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	err := c.transportAdapter.Unsubscribe(resourceName, subscriptionId, c.defaultRequestTimeout)
	if err != nil {
		return newError(err)
	}

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) ([]byte, *Error) {
	resBytes, rawErr := c.transportAdapter.Request(resourceName, data, c.defaultRequestTimeout)
	if rawErr != nil {
		return nil, newError(rawErr)
	}

	return resBytes, nil
}

func (c *Client) getModel(resourceName string, model interface{}, data map[string]interface{}) *Error {
	rawRes, err := c.request(resourceName, data)

	if err != nil {
		return err
	}

	parseErr := json.Unmarshal(rawRes, model)
	if parseErr != nil {
		return newJSONErr(parseErr)
	}

	return nil
}

func (c *Client) GetDevice(deviceId string) (device *Device, err *Error) {
	d := c.NewDevice()

	err = c.getModel(resourcenames.GetDevice, d, map[string]interface{}{
		"deviceId": deviceId,
	})
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Id property of device must be non empty
func (c *Client) PutDevice(id, name string, data map[string]interface{}, networkId, deviceTypeId int, blocked bool) (*Device, *Error) {
	if name == "" {
		name = id
	}

	device := c.NewDevice()
	device.Id = id
	device.Name = name
	device.Data = data
	device.NetworkId = networkId
	device.DeviceTypeId = deviceTypeId
	device.IsBlocked = blocked

	_, err := c.request(resourcenames.PutDevice, map[string]interface{}{
		"deviceId": device.Id,
		"device":   device,
	})

	if err != nil {
		return nil, err
	}

	return device, nil
}

// In case params is nil default values defined at DeviceHive take place
func (c *Client) ListDevices(params *ListParams) (list []*Device, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request(resourcenames.ListDevices, data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr(pErr)
	}

	for _, v := range list {
		v.client = c
	}

	return list, nil
}

func (c *Client) CreateDeviceType(name, description string) (*DeviceType, *Error) {
	devType := c.NewDeviceType()

	devType.Name = name
	devType.Description = description

	res, err := c.request(resourcenames.InsertDeviceType, map[string]interface{}{
		"deviceType": devType,
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, devType)
	if jsonErr != nil {
		return nil, newJSONErr(jsonErr)
	}

	return devType, nil
}

func (c *Client) GetDeviceType(deviceTypeId int) (*DeviceType, *Error) {
	devType := c.NewDeviceType()

	err := c.getModel(resourcenames.GetDeviceType, devType, map[string]interface{}{
		"deviceTypeId": deviceTypeId,
	})
	if err != nil {
		return nil, err
	}

	return devType, nil
}

// In case params is nil default values defined at DeviceHive take place
func (c *Client) ListDeviceTypes(params *ListParams) ([]*DeviceType, *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request(resourcenames.ListDeviceTypes, data)
	if err != nil {
		return nil, err
	}

	var list []*DeviceType
	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr(pErr)
	}
	for _, v := range list {
		v.client = c
	}

	return list, nil
}

// Gets information about DeviceHive server
func (c *Client) GetInfo() (*ServerInfo, *Error) {
	info := &ServerInfo{}

	err := c.getModel(resourcenames.ApiInfo, info, nil)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *Client) GetClusterInfo() (*ClusterInfo, *Error) {
	info := &ClusterInfo{}

	err := c.getModel(resourcenames.ClusterInfo, info, nil)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *Client) CreateNetwork(name, description string) (*Network, *Error) {
	ntwk := c.NewNetwork()

	ntwk.Name = name
	ntwk.Description = description

	res, err := c.request(resourcenames.InsertNetwork, map[string]interface{}{
		"network": ntwk,
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, ntwk)
	if jsonErr != nil {
		return nil, newJSONErr(jsonErr)
	}

	return ntwk, nil
}

func (c *Client) GetNetwork(networkId int) (*Network, *Error) {
	ntwk := c.NewNetwork()

	err := c.getModel(resourcenames.GetNetwork, ntwk, map[string]interface{}{
		"networkId": networkId,
	})
	if err != nil {
		return nil, err
	}

	return ntwk, nil
}

// In case params is nil default values defined at DeviceHive take place
func (c *Client) ListNetworks(params *ListParams) ([]*Network, *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request(resourcenames.ListNetworks, data)
	if err != nil {
		return nil, err
	}

	var list []*Network
	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr(pErr)
	}
	for _, v := range list {
		v.client = c
	}
	return list, nil
}

func (c *Client) GetProperty(name string) (*Configuration, *Error) {
	conf := &Configuration{}

	err := c.getModel(resourcenames.GetConfig, conf, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Client) SetProperty(name, value string) (entityVersion int, err *Error) {
	rawRes, err := c.request(resourcenames.PutConfig, map[string]interface{}{
		"name":  name,
		"value": value,
	})

	if err != nil {
		return -1, err
	}

	conf := &Configuration{}
	parseErr := json.Unmarshal(rawRes, conf)
	if parseErr != nil {
		return -1, newJSONErr(parseErr)
	}

	return conf.EntityVersion, nil
}

func (c *Client) DeleteProperty(name string) *Error {
	_, err := c.request(resourcenames.DeleteConfig, map[string]interface{}{
		"name": name,
	})

	return err
}

func (c *Client) CreateToken(userId int, expiration, refreshExpiration time.Time, actions, networkIds, deviceTypeIds []string) (accessToken, refreshToken string, err *Error) {
	data := map[string]interface{}{
		"userId": userId,
	}

	if actions != nil {
		data["actions"] = actions
	}
	if networkIds != nil {
		data["networkIds"] = networkIds
	}
	if deviceTypeIds != nil {
		data["deviceTypeIds"] = deviceTypeIds
	}
	if expiration.Unix() > 0 {
		data["expiration"] = &ISO8601Time{expiration}
	}
	if refreshExpiration.Unix() > 0 {
		data["refreshExpiration"] = &ISO8601Time{refreshExpiration}
	}

	rawRes, err := c.request(resourcenames.TokenCreate, map[string]interface{}{
		"payload": data,
	})

	if err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", "", newJSONErr(parseErr)
	}

	return tok.Access, tok.Refresh, nil
}

func (c *Client) RefreshToken() (accessToken string, err *Error) {
	accessToken, rawErr := c.transportAdapter.RefreshToken()
	if rawErr != nil {
		err = newError(rawErr)
	}

	return accessToken, err
}

func (c *Client) CreateUser(login, password string, role int, data map[string]interface{}, allDevTypesAvail bool) (*User, *Error) {
	usr := c.NewUser()
	usr.Login = login
	usr.Role = role
	usr.Data = data
	usr.AllDeviceTypesAvailable = allDevTypesAvail

	res, err := c.request(resourcenames.CreateUser, map[string]interface{}{
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

	jsonErr := json.Unmarshal(res, usr)
	if jsonErr != nil {
		return nil, newJSONErr(jsonErr)
	}

	return usr, nil
}

func (c *Client) GetUser(userId int) (*User, *Error) {
	usr := c.NewUser()

	err := c.getModel(resourcenames.GetUser, usr, map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (c *Client) GetCurrentUser() (*User, *Error) {
	usr := c.NewUser()

	err := c.getModel(resourcenames.GetCurrentUser, usr, nil)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (c *Client) ListUsers(params *ListParams) ([]*User, *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request(resourcenames.ListUsers, data)
	if err != nil {
		return nil, err
	}

	var list []*User
	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr(pErr)
	}
	for _, v := range list {
		v.client = c
	}
	return list, nil
}
