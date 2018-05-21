package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
)

func TestDeviceType(t *testing.T) {
	is := is.New(t)

	deviceType, err := client.CreateDeviceType("go-test-dev", "go sdk test device type")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = deviceType.Remove()
		if err != nil {
			t.Fatal(err)
		}
	}()

	is.True(deviceType != nil)
	is.True(deviceType.Id != 0)

	deviceType.Description = "updated go sdk test network"
	err = deviceType.Save()
	if err != nil {
		t.Fatal(err)
	}

	sameDeviceType, err := client.GetDeviceType(deviceType.Id)
	if err != nil {
		t.Fatal(err)
	}

	is.True(sameDeviceType != nil)
	is.Equal(sameDeviceType.Name, "go-test-dev")

	list, err := client.ListDeviceTypes(&dh.ListParams{
		NamePattern: "go-%-dev",
	})
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(len(list), 1)
	is.Equal(list[0].Name, "go-test-dev")
}
