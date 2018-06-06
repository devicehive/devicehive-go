/*
Package devicehive-go provides access to DeviceHive API through WebSocket or REST.
Error handling is omitted to simplify examples:

	client, _ := devicehive_go.ConnectWithCreds("ws://devicehive-address.com/api/websocket", "login", "password")
	// or
	client, _ := devicehive_go.ConnectWithCreds("http://devicehive-address.com/api/rest", "login", "password")

	device, _ := client.PutDevice("my-device", "", nil, 0, 0, false)
	subscription, _ := device.SubscribeInsertCommands(nil, time.Time{})

	done := make(chan struct{})
	go func() {
		for command := range subscription.CommandsChan {
			fmt.Printf("Received command with id %d\n", command.Id)
			close(done)
		}
	}()

	command, _ := device.SendCommand("my-command", nil, 120, time.Time{}, "", nil)

	fmt.Printf("Command with id %d has been sent\n", command.Id)

	<-done

In addition there is an ability to connect with tokens.

	client, err := devicehive_go.ConnectWithToken("ws://devicehive-address.com/api/websocket", "some.JWT.accessToken", "some.JWT.refreshToken")

The client will be automatically reauthenticated by credentials or refresh token in case of access token expiration.

The SDK has an ability to send requests in non-blocking manner, writing each response and error to separate channels that you can read in a separate go routine. This API is called WebSocket low-level API.
WS low-level API usage example:

	wsclient, err := devicehive_go.WSConnect("ws://devicehive-address.com/api/websocket")
	...

	done := make(chan struct{})
	go func() {
		for {
			select {
			case d := <- wsclient.DataChan:
				res := make(map[string]interface{})
				action := ""
				status := ""
				json.Unmarshal(d, &res) // If message was written to DataChan it must be valid JSON

				if a, ok := res["action"].(string); ok {
					action = a
				}

				if s, ok := res["status"].(string); ok {
					status = s
				}

				if action == "authenticate" && status == "success" {
					wsclient.SubscribeCommands(nil)
				} else {
					fmt.Println(string(d))
				}
			case err := <- wsclient.ErrorChan:
				fmt.Println("Error", err)
			}
		}

		close(done)
	}()

	err = wsclient.Authenticate("some.JWT.accessToken")
	...

	<-done
	fmt.Println("Done")
*/
package devicehive_go
