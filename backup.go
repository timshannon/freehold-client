// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import "time"

// Backup is the structure of a freehold backup
type Backup struct {
	When       string   `json:"when"`
	File       string   `json:"file"`
	Who        string   `json:"who"`
	Datastores []string `json:"datastores"`

	whenTime time.Time
}

// GetBackups retrieves the previously generated backups
func (c *Client) GetBackups(from, to time.Time) ([]*Backup, error) {
	var b []*Backup

	fromFmt := from.Format(time.RFC3339)
	toFmt := ""
	if !to.IsZero() {
		toFmt = to.Format(time.RFC3339)
	}
	err := c.doRequest("GET", "/v1/backup/", map[string]string{
		"from": fromFmt,
		"to":   toFmt,
	}, &b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// NewBackup Generates a new freehold instance backup, and returns the path to
// the backup file
func (c *Client) NewBackup(optionalFile string, optionalDSList []string) (string, error) {
	result := ""
	input := make(map[string]interface{})

	input["file"] = optionalFile
	if len(optionalDSList) > 0 {
		input["datastores"] = optionalDSList
	}
	err := c.doRequest("POST", "/v1/backup/", input, &result)

	if err != nil {
		return "", err
	}
	return result, nil
}

// WhenTime is the parsed Time from the Backup
func (b *Backup) WhenTime() time.Time {
	if b.whenTime.IsZero() {
		tme, err := time.Parse(time.RFC3339, b.When)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has bad date!")
		}
		b.whenTime = tme
	}
	return b.whenTime
}
