# DeviceHive SDK for GO

## Installation

    go get -u github.com/devicehive/devicehive-go

## Usage
### Device creation

    client, err := devicehive_go.ConnectWithCreds("ws://devicehive-address.com/api/websocket", "login", "password")
    if err != nil {
        panic(err)
    }

    deviceData := client.NewDevice()
    deviceData.Id = "my-device"
    device, err := client.PutDevice(deviceData)
    if err != nil {
        panic(err)
    }


### Command insert subscription

    client, err := devicehive_go.ConnectWithCreds("ws://devicehive-address.com/api/websocket", "login", "password")
    if err != nil {
        panic(err)
    }

    device, err := client.GetDevice("my-device")
    if err != nil {
        panic(err)
    }

    subscription, err := device.SubscribeInsertCommands(nil, time.Time{})
    if err != nil {
        panic(err)
    }

    for command := range subscription.CommandsChan {
        fmt.Println(command)
    }
