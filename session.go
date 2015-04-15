// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

// Session is a freehold session, tracked by a cookie
type Session struct {
	ID        string `json:"id,omitempty"`
	Expires   string `json:"expires,omitempty"`
	CSRFToken string `json:"CSRFToken,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	Created   string `json:"created,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
}
