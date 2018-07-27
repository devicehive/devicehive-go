// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package requester

import "github.com/devicehive/devicehive-go/internal/resourcenames"

var httpResources = map[string][2]string{
	resourcenames.TokenCreate:            {"token/create", "POST"},
	resourcenames.TokenRefresh:           {"token/refresh", "POST"},
	resourcenames.TokenByCreds:           {"token", "POST"},
	resourcenames.ApiInfo:                {"info"},
	resourcenames.ClusterInfo:            {"info/config/cluster"},
	resourcenames.PutConfig:              {"configuration/{{.name}}", "PUT"},
	resourcenames.GetConfig:              {"configuration/{{.name}}"},
	resourcenames.DeleteConfig:           {"configuration/{{.name}}", "DELETE"},
	resourcenames.PutDevice:              {"device/{{.deviceId}}", "PUT"},
	resourcenames.GetDevice:              {"device/{{.deviceId}}"},
	resourcenames.ListDevices:            {"device"},
	resourcenames.DeleteDevice:           {"device/{{.deviceId}}", "DELETE"},
	resourcenames.InsertCommand:          {"device/{{.deviceId}}/command", "POST"},
	resourcenames.ListCommands:           {`device/{{.deviceId}}/command`},
	resourcenames.UpdateCommand:          {"device/{{.deviceId}}/command/{{.commandId}}", "PUT"},
	resourcenames.InsertNotification:     {"device/{{.deviceId}}/notification", "POST"},
	resourcenames.ListNotifications:      {"device/{{.deviceId}}/notification"},
	resourcenames.SubscribeCommands:      {"device/command/poll"},
	resourcenames.SubscribeNotifications: {"device/notification/poll"},
	resourcenames.InsertNetwork:          {"network", "POST"},
	resourcenames.DeleteNetwork:          {"network/{{.networkId}}", "DELETE"},
	resourcenames.UpdateNetwork:          {"network/{{.networkId}}", "PUT"},
	resourcenames.GetNetwork:             {"network/{{.networkId}}"},
	resourcenames.ListNetworks:           {"network"},
	resourcenames.InsertDeviceType:       {"devicetype", "POST"},
	resourcenames.UpdateDeviceType:       {"devicetype/{{.deviceTypeId}}", "PUT"},
	resourcenames.DeleteDeviceType:       {"devicetype/{{.deviceTypeId}}", "DELETE"},
	resourcenames.GetDeviceType:          {"devicetype/{{.deviceTypeId}}"},
	resourcenames.ListDeviceTypes:        {"devicetype"},
	resourcenames.CreateUser:             {"user", "POST"},
	resourcenames.DeleteUser:             {"user/{{.userId}}", "DELETE"},
	resourcenames.GetUser:                {"user/{{.userId}}"},
	resourcenames.GetCurrentUser:         {"user/current"},
	resourcenames.ListUsers:              {"user"},
	resourcenames.UpdateUser:             {"user/{{.userId}}", "PUT"},
	resourcenames.AssignNetwork:          {"user/{{.userId}}/network/{{.networkId}}", "PUT"},
	resourcenames.UnassignNetwork:        {"user/{{.userId}}/network/{{.networkId}}", "DELETE"},
	resourcenames.AssignDeviceType:       {"user/{{.userId}}/devicetype/{{.deviceTypeId}}", "PUT"},
	resourcenames.UnassignDeviceType:     {"user/{{.userId}}/devicetype/{{.deviceTypeId}}", "DELETE"},
	resourcenames.GetUserDeviceTypes:     {"user/{{.userId}}/devicetype"},
	resourcenames.AllowAllDeviceTypes:    {"user/{{.userId}}/devicetype/all", "PUT"},
	resourcenames.DisallowAllDeviceTypes: {"user/{{.userId}}/devicetype/all", "DELETE"},
}

var httpResourcesQueryParams = map[string][]string{
	resourcenames.DeleteNetwork:          {"force"},
	resourcenames.DeleteDeviceType:       {"force"},
	resourcenames.ListCommands:           {"start", "end", "command", "status", "sortField", "sortOrder", "take", "skip"},
	resourcenames.ListNotifications:      {"start", "end", "notification", "sortField", "sortOrder", "take", "skip"},
	resourcenames.SubscribeCommands:      {"deviceId", "networkIds", "deviceTypeIds", "timestamp", "waitTimeout", "names"},
	resourcenames.SubscribeNotifications: {"deviceId", "networkIds", "deviceTypeIds", "timestamp", "waitTimeout", "names"},
	resourcenames.ListDevices:            {"name", "namePattern", "networkId", "networkName", "sortField", "sortOrder", "take", "skip"},
	resourcenames.ListNetworks:           {"name", "namePattern", "sortField", "sortOrder", "take", "skip"},
	resourcenames.ListDeviceTypes:        {"name", "namePattern", "sortField", "sortOrder", "take", "skip"},
	resourcenames.ListUsers:              {"login", "loginPattern", "role", "status", "sortField", "sortOrder", "take", "skip"},
}
