package dh_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}
