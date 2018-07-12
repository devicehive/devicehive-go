// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package apirequests

type WSSubscription struct {
	*PendingRequest
	SubscriptionResource string
	SubscriptionParams   *RequestParams
	SubscriptionId       string
}
