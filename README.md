# DeviceHive SDK for GO

## Installation

    go get -u github.com/devicehive/devicehive-go

## Usage
### Device creation

    import "github.com/devicehive/devicehive-go"

    func main() {
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
    }


### Command insert subscription

    import (
        "github.com/devicehive/devicehive-go"
        "fmt"
    )

    func main() {
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
    }

## Running tests
Integration tests of high-level DH client:

    go test github.com/devicehive/devicehive-go/integrationtest/dh -serverAddress ws://devicehive-api.com/ -accessToken your.accessToken -refreshToken your.accessToken -userId 123

Integration tests of low-level DH WS client (only ws:// URL is acceptable for this tests as server address):

    go test github.com/devicehive/devicehive-go/integrationtest/dh_wsclient -serverAddress ws://devicehive-api.com/ -accessToken your.accessToken -refreshToken your.accessToken -userId 123

Unit tests:

    go test github.com/devicehive/devicehive-go/test/...

