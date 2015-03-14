// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
)

//FHError is an error returned by a freehold instance
type FHError struct {
	url        string
	status     string
	statusCode int
	message    string
}

func isError(url string, statusCode int, response *jsend) error {
	if response == nil {
		status := "success"
		if statusCode >= 500 {
			status = "error"
		} else if statusCode >= 400 {
			status = "fail"
		}

		response = &jsend{
			Status:  status,
			Message: http.StatusText(statusCode),
		}
	}
	if response.Status == "success" {
		return nil
	}

	return &FHError{
		url:        url,
		status:     response.Status,
		statusCode: statusCode,
		message:    response.Message,
	}
}

func (e *FHError) Error() string {
	return fmt.Sprintf("Request %s failed with a status of %d.  Message: %s", e.url, e.statusCode, e.message)
}

// IsNotFound returns whether or not the error is a
// 404
func IsNotFound(err error) bool {
	switch e := err.(type) {
	case nil:
		return false
	case *FHError:
		return e.statusCode == http.StatusNotFound
	}
	return false
}
