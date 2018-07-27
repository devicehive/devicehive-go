// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package requester

import "github.com/devicehive/devicehive-go/internal/resourcenames"

var wsResources = map[string]string{
	resourcenames.Auth:                   "authenticate",
	resourcenames.TokenCreate:            "token/create",
	resourcenames.TokenRefresh:           "token/refresh",
	resourcenames.TokenByCreds:           "token",
	resourcenames.ApiInfo:                "server/info",
	resourcenames.ClusterInfo:            "cluster/info",
	resourcenames.PutConfig:              "configuration/put",
	resourcenames.GetConfig:              "configuration/get",
	resourcenames.DeleteConfig:           "configuration/delete",
	resourcenames.PutDevice:              "device/save",
	resourcenames.GetDevice:              "device/get",
	resourcenames.ListDevices:            "device/list",
	resourcenames.DeleteDevice:           "device/delete",
	resourcenames.InsertCommand:          "command/insert",
	resourcenames.ListCommands:           "command/list",
	resourcenames.UpdateCommand:          "command/update",
	resourcenames.InsertNotification:     "notification/insert",
	resourcenames.ListNotifications:      "notification/list",
	resourcenames.SubscribeCommands:      "command/subscribe",
	resourcenames.SubscribeNotifications: "notification/subscribe",
	resourcenames.InsertNetwork:          "network/insert",
	resourcenames.DeleteNetwork:          "network/delete",
	resourcenames.UpdateNetwork:          "network/update",
	resourcenames.GetNetwork:             "network/get",
	resourcenames.ListNetworks:           "network/list",
	resourcenames.InsertDeviceType:       "devicetype/insert",
	resourcenames.UpdateDeviceType:       "devicetype/update",
	resourcenames.DeleteDeviceType:       "devicetype/delete",
	resourcenames.GetDeviceType:          "devicetype/get",
	resourcenames.ListDeviceTypes:        "devicetype/list",
	resourcenames.CreateUser:             "user/insert",
	resourcenames.DeleteUser:             "user/delete",
	resourcenames.GetUser:                "user/get",
	resourcenames.GetCurrentUser:         "user/getCurrent",
	resourcenames.ListUsers:              "user/list",
	resourcenames.UpdateUser:             "user/update",
	resourcenames.AssignNetwork:          "user/assignNetwork",
	resourcenames.UnassignNetwork:        "user/unassignNetwork",
	resourcenames.AssignDeviceType:       "user/assignDeviceType",
	resourcenames.UnassignDeviceType:     "user/unassignDeviceType",
	resourcenames.GetUserDeviceTypes:     "user/getDeviceTypes",
	resourcenames.AllowAllDeviceTypes:    "user/allowAllDeviceTypes",
	resourcenames.DisallowAllDeviceTypes: "user/disallowAllDeviceTypes",
}
