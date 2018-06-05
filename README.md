# DeviceHive SDK for GO

Generally Golang SDK for DeviceHive consists of 2 APIs, feel free to choose yours:
- client (a.k.a. high-level client) — provides synchronous API and ORM-like access to DeviceHive models
- WS client (a.k.a. WS low-level client) — provides asynchronous API: just sends the request and returns an error only in case of request error,
all raw response data and response errors are written to appropriate channels which are created after WS connection is established
(for more details see [documentation at godoc](https://godoc.org/github.com/devicehive/devicehive-go))

## Installation

    go get -u github.com/devicehive/devicehive-go

## Documentation
Visit https://godoc.org/github.com/devicehive/devicehive-go for full API reference.

## Usage
### Connection

    import "github.com/devicehive/devicehive-go"

    func main() {
        client, err := devicehive_go.ConnectWithCreds("ws://devicehive-address.com/api/websocket", "login", "password")
        // OR
        client, err := devicehive_go.ConnectWithToken("ws://devicehive-address.com/api/websocket", "jwt.Access.Token", "jwt.Refresh.Token")
        if err != nil {
            fmt.Println(err)
            return
        }
    }

### Device creation

    import "github.com/devicehive/devicehive-go"

    func main() {
        client, err := devicehive_go.ConnectWithCreds("ws://devicehive-address.com/api/websocket", "login", "password")
        if err != nil {
            fmt.Println(err)
            return
        }

        deviceData := client.NewDevice()
        deviceData.Id = "my-device"
        device, err := client.PutDevice(*deviceData)
        if err != nil {
            fmt.Println(err)
            return
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
            fmt.Println(err)
            return
        }

        device, err := client.GetDevice("my-device")
        if err != nil {
            fmt.Println(err)
            return
        }

        subscription, err := device.SubscribeInsertCommands(nil, time.Time{})
        if err != nil {
            fmt.Println(err)
            return
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

