// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

// Application is the structure of an Application Install
type Application struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Root        string `json:"root,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`

	client *Client
}

// AvailableApplication is an application file available for install to a freehold instance
type AvailableApplication struct {
	Application
	File string `json:"file,omitempty"`
}

// AllApplications retrieves all installed applications
func (c *Client) AllApplications() ([]*Application, error) {
	a := make(map[string]*Application)

	err := c.doRequest("GET", "/v1/application/", nil, &a)
	if err != nil {
		return nil, err
	}

	apps := make([]*Application, 0, len(a))

	for k := range a {
		a[k].ID = k
		a[k].client = c
		apps = append(apps, a[k])
	}

	return apps, nil
}

// GetApplication retrieves a specific Application
func (c *Client) GetApplication(appID string) (*Application, error) {
	a := &Application{}

	err := c.doRequest("GET", "/v1/application/", map[string]string{
		"id": appID,
	}, &a)
	if err != nil {
		return nil, err
	}
	a.client = c
	return a, nil
}

// AvailableApplications retrieves all the application files available for install
func (c *Client) AvailableApplications() ([]*AvailableApplication, error) {
	a := make(map[string]*AvailableApplication)

	err := c.doRequest("GET", "/v1/application/available", nil, &a)
	if err != nil {
		return nil, err
	}

	apps := make([]*AvailableApplication, 0, len(a))

	for k := range a {
		a[k].ID = k
		a[k].client = c
		apps = append(apps, a[k])
	}

	return apps, nil
}

// PostAvailableApplication posts a new available application for install from the passed in URL
func (c *Client) PostAvailableApplication(url string) (*AvailableApplication, error) {
	a := &AvailableApplication{}

	err := a.client.doRequest("POST", "/v1/application/available", map[string]string{
		"file": url,
	}, &a.File)

	if err != nil {
		return nil, err
	}
	a.client = c
	return a, nil
}

// Install installs the available application
func (a *AvailableApplication) Install() (*Application, error) {
	app := &Application{}

	err := a.client.doRequest("POST", "/v1/application/", map[string]string{
		"file": a.File,
	}, &app)
	if err != nil {
		return nil, err
	}

	app.client = a.client

	return app, nil
}

// Upgrade Upgrades a currently installed application
func (a *AvailableApplication) Upgrade() (*Application, error) {
	app := &Application{}

	err := a.client.doRequest("PUT", "/v1/application/", map[string]string{
		"file": a.File,
	}, &app)
	if err != nil {
		return nil, err
	}
	app.client = a.client

	return app, nil
}

// Uninstall removes the installed application from the freehold instance
func (a *Application) Uninstall() error {
	return a.client.doRequest("DELETE", "/v1/application/", map[string]string{
		"id": a.ID,
	}, nil)
}
