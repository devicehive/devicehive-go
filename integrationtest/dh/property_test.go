package dh_test

//import (
//	"testing"
//	"github.com/matryer/is"
//	"strconv"
//	"time"
//)
//
//func TestProperty(t *testing.T) {
//	is := is.New(t)
//
//	name := "go-test" + strconv.FormatInt(time.Now().Unix(), 10)
//	val := "go-sdk-test"
//
//	entityVersion, dhErr := client.SetProperty(name, val)
//	if dhErr != nil {
//		t.Fatalf("%s: %v", dhErr.Name(), dhErr)
//	}
//
//	is.True(entityVersion >= 0)
//
//	prop, dhErr := client.GetProperty(name)
//	if dhErr != nil {
//		t.Fatalf("%s: %v", dhErr.Name(), dhErr)
//	}
//
//	is.True(prop != nil)
//	is.Equal(prop.Name, name)
//	is.Equal(prop.Value, val)
//	is.True(prop.EntityVersion >= 0)
//
//	dhErr = client.DeleteProperty(name)
//	if dhErr != nil {
//		t.Fatalf("%s: %v", dhErr.Name(), dhErr)
//	}
//}
