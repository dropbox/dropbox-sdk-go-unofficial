// Copyright (c) Dropbox, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package contacts

import (
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// DeleteManualContacts : Removes all manually added contacts. You'll still
	// keep contacts who are on your team or who you imported. New contacts will
	// be added when you share.
	DeleteManualContacts() (err error)
	// DeleteManualContactsBatch : Removes manually added contacts from the
	// given list.
	DeleteManualContactsBatch(arg *DeleteManualContactsArg) (err error)
}

type apiImpl dropbox.Context

//DeleteManualContactsAPIError is an error-wrapper for the delete_manual_contacts route
type DeleteManualContactsAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) DeleteManualContacts() (err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "contacts",
		Route:        "delete_manual_contacts",
		Auth:         "user",
		Style:        "rpc",
		Arg:          nil,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DeleteManualContactsAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

//DeleteManualContactsBatchAPIError is an error-wrapper for the delete_manual_contacts_batch route
type DeleteManualContactsBatchAPIError struct {
	dropbox.APIError
	EndpointError *DeleteManualContactsError `json:"error"`
}

func (dbx *apiImpl) DeleteManualContactsBatch(arg *DeleteManualContactsArg) (err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "contacts",
		Route:        "delete_manual_contacts_batch",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DeleteManualContactsBatchAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

// New returns a Client implementation for this namespace
func New(c dropbox.Config) Client {
	ctx := apiImpl(dropbox.NewContext(c))
	return &ctx
}
