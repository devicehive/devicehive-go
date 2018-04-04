package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
	"strconv"
)

func TestConfiguration(t *testing.T) {
	is := is.New(t)

	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	name, val := "go-test" + strconv.FormatInt(time.Now().Unix(), 10), "go-sdk-test"

	conf, dhErr := client.ConfigurationPut(name, val)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(conf != nil)
	is.Equal(conf.Name, name)
	is.Equal(conf.Value, val)
	is.Equal(conf.EntityVersion, 0)

	conf, dhErr = client.ConfigurationGet(name)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(conf != nil)
	is.Equal(conf.Name, name)
	is.Equal(conf.Value, val)
	is.True(conf.EntityVersion == 0)

	dhErr = client.ConfigurationDelete(name)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}
}
