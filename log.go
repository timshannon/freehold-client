// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import "time"

// Log is the storage stucture for a log entry
type Log struct {
	When string `json:"when"`
	Type string `json:"type"`
	Log  string `json:"log"`

	whenTime time.Time
}

// LogIter is used for iterating through freehold logs
type LogIter struct {
	Iter
	Type string `json:"type,omitempty"`
}

// GetLogs retrieves the logs that match the passed in Log Iterator
func (c *Client) GetLogs(iter *LogIter) ([]*Log, error) {
	var l []*Log
	err := c.doRequest("GET", "/v1/log/", iter, &l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// WhenTime is the parsed Time from the Log
func (l *Log) WhenTime() time.Time {
	if l.whenTime.IsZero() {
		tme, err := time.Parse(time.RFC3339, l.When)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has bad date!")
		}
		l.whenTime = tme
	}
	return l.whenTime
}
