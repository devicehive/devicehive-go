/*
Package devicehive-go provides access to DeviceHive API through WebSocket or REST.

	client, err := devicehive_go.ConnectWithCreds("ws://devicehive-address.com/api/websocket", "login", "password")
	// or
	client, err := devicehive_go.ConnectWithCreds("http://devicehive-address.com/api/rest", "login", "password")
	...
	deviceData := client.NewDevice()
	deviceData.Id = "device-id"
	deviceData.NetworkId = 1
	deviceData.DeviceTypeId = 1

	device, err := client.PutDevice(deviceData)
	...
	subscription, err := device.SubscribeInsertCommands(nil, time.Time{})
	...
	go func() {
		for command := range subscription.CommandsChan {
			fmt.Println(command)
		}
	}()

	command, err := device.SendCommand("my-command", nil, 120, time.Time{}, "", nil)
	...

In addition there is an ability to connect with tokens.

	client, err := devicehive_go.ConnectWithToken("ws://devicehive-address.com/api/websocket", "some.JWT.accessToken", "some.JWT.refreshToken")

The client will be automatically reauthenticated by credentials or refresh token in case of access token expiration
*/
package devicehive_go
