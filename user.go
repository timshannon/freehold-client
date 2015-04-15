// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

// User is a user in a freehold instance
type User struct {
	Username string `json:"-"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	HomeApp  string `json:"homeApp,omitempty"`
	Admin    bool   `json:"admin,omitempty"`

	client *Client
}

// AllUsers retrieves all Users in the freehold instance
func (c *Client) AllUsers() ([]*User, error) {
	u := make(map[string]*User)
	err := c.doRequest("GET", "/v1/auth/user/", nil, &u)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(u))

	for k := range u {
		u[k].client = c
		u[k].Username = k
		users = append(users, u[k])
	}

	return users, nil
}

// GetUser retrieves a User in the freehold instance
func (c *Client) GetUser(username string) (*User, error) {
	u := &User{}
	err := c.doRequest("GET", "/v1/auth/user/", map[string]string{
		"user": username,
	}, u)
	if err != nil {
		return nil, err
	}

	u.client = c
	u.Username = username

	return u, nil
}

// NewUser creates a new user
func (c *Client) NewUser(username, password, name, homeApp string, isAdmin bool) (*User, error) {
	input := map[string]interface{}{
		"user":     username,
		"password": password,
		"name":     name,
		"homeApp":  homeApp,
		"admin":    isAdmin,
	}
	u := &User{}
	err := c.doRequest("POST", "/v1/auth/user/", input, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Delete deletes a user
func (u *User) Delete() error {
	return u.client.doRequest("DELETE", "/v1/auth/user/", map[string]string{
		"user": u.Username,
	}, nil)
}

// SetName sets the user's name
func (u *User) SetName(newName string) error {
	return u.client.doRequest("PUT", "/v1/auth/user/", map[string]string{
		"user": u.Username,
		"name": newName,
	}, nil)
}

// SetPassword sets the user's password
func (u *User) SetPassword(newPassword string) error {
	return u.client.doRequest("PUT", "/v1/auth/user/", map[string]string{
		"user":     u.Username,
		"password": newPassword,
	}, nil)
}

// SetHomeApp sets the user's home appliation
func (u *User) SetHomeApp(newHomeApp string) error {
	return u.client.doRequest("PUT", "/v1/auth/user/", map[string]string{
		"user":    u.Username,
		"homeApp": newHomeApp,
	}, nil)
}

// SetAdmin sets if the user is an admin or not
func (u *User) SetAdmin(isAdmin bool) error {
	return u.client.doRequest("PUT", "/v1/auth/user/", map[string]interface{}{
		"user":  u.Username,
		"admin": isAdmin,
	}, nil)
}
