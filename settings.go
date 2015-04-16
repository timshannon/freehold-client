// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

// Setting is a value that changes how the freehold instance operates
type Setting struct {
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value,omitempty"`
}

// AllSettings retrieves all the current settings
// for the freehold instance
func (c *Client) AllSettings() (map[string]*Setting, error) {
	s := make(map[string]*Setting)

	err := c.doRequest("GET", "/v1/settings/", nil, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// GetSetting gets a specific setting
func (c *Client) GetSetting(settingName string) (*Setting, error) {
	s := &Setting{}

	err := c.doRequest("GET", "/v1/settings/", map[string]string{
		"setting": settingName,
	}, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// SetSetting sets the given setting's value
func (c *Client) SetSetting(settingName string, value interface{}) error {
	return c.doRequest("PUT", "/v1/settings/", map[string]interface{}{
		"setting": settingName,
		"value":   value,
	}, nil)
}

// DefaultSetting sets the given setting's value
func (c *Client) DefaultSetting(settingName string) error {
	return c.doRequest("DELETE", "/v1/settings/", map[string]string{
		"setting": settingName,
	}, nil)
}
