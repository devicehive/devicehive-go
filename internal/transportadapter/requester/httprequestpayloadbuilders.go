// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package requester

import "github.com/devicehive/devicehive-go/internal/resourcenames"

var httpRequestPayloadBuilders = map[string]func(map[string]interface{}) interface{}{
	resourcenames.TokenCreate: func(data map[string]interface{}) interface{} {
		return data["payload"]
	},
	resourcenames.PutConfig: func(data map[string]interface{}) interface{} {
		return map[string]interface{}{
			"value": data["value"],
		}
	},
	resourcenames.DeleteConfig: func(data map[string]interface{}) interface{} {
		return nil
	},
	resourcenames.GetConfig: func(data map[string]interface{}) interface{} {
		return nil
	},
	resourcenames.PutDevice: func(data map[string]interface{}) interface{} {
		return data["device"]
	},
	resourcenames.GetDevice: func(data map[string]interface{}) interface{} {
		return nil
	},
	resourcenames.InsertCommand: func(data map[string]interface{}) interface{} {
		return data["command"]
	},
	resourcenames.ListCommands: func(data map[string]interface{}) interface{} {
		return nil
	},
	resourcenames.UpdateCommand: func(data map[string]interface{}) interface{} {
		return data["command"]
	},
	resourcenames.InsertNotification: func(data map[string]interface{}) interface{} {
		return data["notification"]
	},
	resourcenames.InsertNetwork: func(data map[string]interface{}) interface{} {
		return data["network"]
	},
	resourcenames.UpdateNetwork: func(data map[string]interface{}) interface{} {
		return data["network"]
	},
	resourcenames.InsertDeviceType: func(data map[string]interface{}) interface{} {
		return data["deviceType"]
	},
	resourcenames.UpdateDeviceType: func(data map[string]interface{}) interface{} {
		return data["deviceType"]
	},
	resourcenames.CreateUser: func(data map[string]interface{}) interface{} {
		return data["user"]
	},
	resourcenames.UpdateUser: func(data map[string]interface{}) interface{} {
		return data["user"]
	},
}
