package devicehive_go

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
	"time"
)

type MainClientInterface interface {
	SubscribeNotifications(*SubscribeParams) (*NotificationSubscription, *Error)
	request(string, map[string]interface{}) ([]byte, *Error)
	SubscribeCommands(*SubscribeParams) (*CommandSubscription, *Error)
	NewCommand() *command
	NewNetwork() *network
	NewUser() *user
	NewDeviceType() *deviceType
	NewDevice() *device
}

type Client struct {
	transport                 transport.Transporter
	transportAdapter          transportadapter.TransportAdapter
	refreshToken              string
	login                     string
	password                  string
	PollingWaitTimeoutSeconds int
}

func (c *Client) NewDevice() *device {
	return &device{client: c}
}

func (c *Client) NewDeviceType() *deviceType {
	return &deviceType{client: c}
}

func (c *Client) NewUser() *user {
	return &user{client: c}
}

func (c *Client) NewNetwork() *network {
	return &network{client: c}
}

func (c *Client) NewCommand() *command {
	return &command{client: c}
}

func (c *Client) SubscribeNotifications(params *SubscribeParams) (subs *NotificationSubscription, err *Error) {
	tspChan, subsId, err := c.subscribe("subscribeNotifications", params)
	if err != nil || tspChan == nil {
		return nil, err
	}

	subs = newNotificationSubscription(subsId, tspChan, c)

	return subs, nil
}

func (c *Client) SubscribeCommands(params *SubscribeParams) (subs *CommandSubscription, err *Error) {
	tspChan, subsId, err := c.subscribe("subscribeCommands", params)
	if err != nil || tspChan == nil {
		return nil, err
	}

	subs = newCommandSubscription(subsId, tspChan, c)

	return subs, nil
}

func (c *Client) authenticate(token string) (result bool, err *Error) {
	result, rawErr := c.transportAdapter.Authenticate(token, Timeout)

	if rawErr != nil {
		return false, newError(rawErr)
	}

	return true, nil
}

func (c *Client) subscribe(resourceName string, params *SubscribeParams) (tspChan chan []byte, subscriptionId string, err *Error) {
	if params == nil {
		params = &SubscribeParams{}
	}

	params.WaitTimeout = c.PollingWaitTimeoutSeconds

	data, jsonErr := params.Map()
	if jsonErr != nil {
		return nil, "", &Error{name: InvalidRequestErr, reason: jsonErr.Error()}
	}

	tspChan, subscriptionId, rawErr := c.transportAdapter.Subscribe(resourceName, c.PollingWaitTimeoutSeconds, data)
	if rawErr != nil {
		return nil, "", newTransportErr(rawErr)
	}

	return tspChan, subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	err := c.transportAdapter.Unsubscribe(resourceName, subscriptionId, Timeout)
	if err != nil {
		return newError(err)
	}

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	resBytes, rawErr := c.transportAdapter.Request(resourceName, data, Timeout)

	if rawErr != nil && rawErr.Error() == "401 token expired" {
		resBytes, err = c.refreshRetry(resourceName, data)
		if err != nil {
			return nil, err
		}
	} else {
		err = newError(rawErr)
	}

	return resBytes, err
}

func (c *Client) refreshRetry(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	accessToken, err := c.RefreshToken()
	if err != nil {
		return nil, err
	}

	res, err := c.authenticate(accessToken)
	if !res || err != nil {
		return nil, err
	}

	resBytes, rawErr := c.transportAdapter.Request(resourceName, data, Timeout)
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
		return newJSONErr()
	}

	return nil
}

func (c *Client) GetDevice(deviceId string) (device *device, err *Error) {
	d := c.NewDevice()

	err = c.getModel("getDevice", d, map[string]interface{}{
		"deviceId": deviceId,
	})
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) PutDevice(device device) (*device, *Error) {
	if device.Name == "" {
		device.Name = device.Id
	}

	_, err := c.request("putDevice", map[string]interface{}{
		"deviceId": device.Id,
		"device":   device,
	})

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (c *Client) ListDevices(params *ListParams) (list []*device, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request("listDevices", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	for _, v := range list {
		v.client = c
	}

	return list, nil
}

func (c *Client) CreateDeviceType(name, description string) (devType *deviceType, err *Error) {
	devType = c.NewDeviceType()

	devType.Name = name
	devType.Description = description

	res, err := c.request("insertDeviceType", map[string]interface{}{
		"deviceType": devType,
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, devType)
	if jsonErr != nil {
		return nil, newJSONErr()
	}

	return devType, nil
}

func (c *Client) GetDeviceType(deviceTypeId int) (devType *deviceType, err *Error) {
	devType = c.NewDeviceType()

	err = c.getModel("getDeviceType", devType, map[string]interface{}{
		"deviceTypeId": deviceTypeId,
	})
	if err != nil {
		return nil, err
	}

	return devType, nil
}

func (c *Client) ListDeviceTypes(params *ListParams) (list []*deviceType, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request("listDeviceTypes", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}
	for _, v := range list {
		v.client = c
	}

	return list, nil
}

func (c *Client) GetInfo() (info *ServerInfo, err *Error) {
	info = &ServerInfo{}

	err = c.getModel("apiInfo", info, nil)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *Client) GetClusterInfo() (info *ClusterInfo, err *Error) {
	info = &ClusterInfo{}

	err = c.getModel("apiInfoCluster", info, nil)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *Client) CreateNetwork(name, description string) (ntwk *network, err *Error) {
	ntwk = c.NewNetwork()

	ntwk.Name = name
	ntwk.Description = description

	res, err := c.request("insertNetwork", map[string]interface{}{
		"network": ntwk,
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, ntwk)
	if jsonErr != nil {
		return nil, newJSONErr()
	}

	return ntwk, nil
}

func (c *Client) GetNetwork(networkId int) (ntwk *network, err *Error) {
	ntwk = c.NewNetwork()

	err = c.getModel("getNetwork", ntwk, map[string]interface{}{
		"networkId": networkId,
	})
	if err != nil {
		return nil, err
	}

	return ntwk, nil
}

func (c *Client) ListNetworks(params *ListParams) (list []*network, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request("listNetworks", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}
	for _, v := range list {
		v.client = c
	}
	return list, nil
}

func (c *Client) GetProperty(name string) (conf *Configuration, err *Error) {
	conf = &Configuration{}

	err = c.getModel("getConfig", conf, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Client) SetProperty(name, value string) (entityVersion int, err *Error) {
	rawRes, err := c.request("putConfig", map[string]interface{}{
		"name":  name,
		"value": value,
	})

	if err != nil {
		return -1, err
	}

	conf := &Configuration{}
	parseErr := json.Unmarshal(rawRes, conf)
	if parseErr != nil {
		return -1, newJSONErr()
	}

	return conf.EntityVersion, nil
}

func (c *Client) DeleteProperty(name string) *Error {
	_, err := c.request("deleteConfig", map[string]interface{}{
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
		data["expiration"] = expiration.UTC().Format(timestampLayout)
	}
	if refreshExpiration.Unix() > 0 {
		data["refreshExpiration"] = refreshExpiration.UTC().Format(timestampLayout)
	}

	return c.tokenRequest("tokenCreate", map[string]interface{}{
		"payload": data,
	})
}

func (c *Client) RefreshToken() (accessToken string, err *Error) {
	if c.refreshToken == "" {
		accessToken, _, err = c.tokensByCreds(c.login, c.password)
		return accessToken, err
	}

	return c.accessTokenByRefresh(c.refreshToken)
}

func (c *Client) accessTokenByRefresh(refreshToken string) (accessToken string, err *Error) {
	rawRes, err := c.request("tokenRefresh", map[string]interface{}{
		"refreshToken": c.refreshToken,
	})

	if err != nil {
		return "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", newJSONErr()
	}

	return tok.Access, nil
}

func (c *Client) tokensByCreds(login, pass string) (accessToken, refreshToken string, err *Error) {
	return c.tokenRequest("tokenByCreds", map[string]interface{}{
		"login":    login,
		"password": pass,
	})
}

func (c *Client) tokenRequest(resourceName string, data map[string]interface{}) (accessToken, refreshToken string, err *Error) {
	rawRes, err := c.request(resourceName, data)

	if err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", "", newJSONErr()
	}

	return tok.Access, tok.Refresh, nil
}

func (c *Client) CreateUser(login, password string, role int, data map[string]interface{}, allDevTypesAvail bool) (*user, *Error) {
	usr := c.NewUser()
	usr.Login = login
	usr.Role = role
	usr.Data = data
	usr.AllDeviceTypesAvailable = allDevTypesAvail

	res, err := c.request("createUser", map[string]interface{}{
		"user": map[string]interface{}{
			"login":                   login,
			"role":                    role,
			"status":                  UserStatusActive,
			"password":                password,
			"data":                    data,
			"allDeviceTypesAvailable": allDevTypesAvail,
		},
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, usr)
	if jsonErr != nil {
		return nil, newJSONErr()
	}

	return usr, nil
}

func (c *Client) GetUser(userId int) (usr *user, err *Error) {
	usr = c.NewUser()

	err = c.getModel("getUser", usr, map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (c *Client) GetCurrentUser() (usr *user, err *Error) {
	usr = c.NewUser()

	err = c.getModel("getCurrentUser", usr, nil)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (c *Client) ListUsers(params *ListParams) (list []*user, err *Error) {
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
	for _, v := range list {
		v.client = c
	}
	return list, nil
}
