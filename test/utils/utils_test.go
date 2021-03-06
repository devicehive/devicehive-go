// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/devicehive/devicehive-go/internal/utils"
	"github.com/matryer/is"
	"testing"
)

func TestCastInterfaceSliceToStringSlice(t *testing.T) {
	is := is.New(t)

	validSlice := []interface{}{"elem1", "elem2", "elem3"}
	invalidSlice := []interface{}{"elem1", 2, 3}

	strSlice, err := utils.ISliceToStrSlice(validSlice)

	is.NoErr(err)
	is.Equal(len(strSlice), 3)

	strSlice, err = utils.ISliceToStrSlice(invalidSlice)

	is.Equal(err.Error(), "element is not string: 2")
	is.True(strSlice == nil)
}
