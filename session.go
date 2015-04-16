// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import "time"

// Session is a freehold session, tracked by a cookie
type Session struct {
	ID        string `json:"id,omitempty"`
	Expires   string `json:"expires,omitempty"`
	CSRFToken string `json:"CSRFToken,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	Created   string `json:"created,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`

	expiresTime time.Time
	createdTime time.Time
	client      *Client
}

// AllSessions retrieves all sessions for the given user
// who made the client connection
func (c *Client) AllSessions() ([]*Session, error) {
	var s []*Session
	err := c.doRequest("GET", "/v1/auth/session/", nil, &s)
	if err != nil {
		return nil, err
	}

	for i := range s {
		s[i].client = c
	}

	return s, nil
}

// Delete deletes the current session from the freehold instance
// making it invalid for all future uses
func (s *Session) Delete() error {
	return s.client.doRequest("DELETE", "/v1/auth/session/",
		map[string]string{
			"id": s.ID,
		}, nil)
}

// ExpiresTime is the parsed Time from the Session's JSON string response
func (s *Session) ExpiresTime() time.Time {
	if s.expiresTime.IsZero() {
		tme, err := time.Parse(time.RFC3339, s.Expires)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has bad Expired date!")
		}
		s.expiresTime = tme
	}
	return s.expiresTime
}

// CreatedTime is the parsed Time from the Session's JSON string response
func (s *Session) CreatedTime() time.Time {
	if s.createdTime.IsZero() {
		tme, err := time.Parse(time.RFC3339, s.Created)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has bad date!")
		}
		s.createdTime = tme
	}
	return s.createdTime
}
