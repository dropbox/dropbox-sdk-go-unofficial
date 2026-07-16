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
	"context"
	"errors"
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

// ContextClient interface describes all routes in this namespace with context support
type ContextClient interface {
	Client
	// DeleteManualContactsContext : Removes all manually added contacts. You'll
	// still keep contacts who are on your team or who you imported. New
	// contacts will be added when you share.
	DeleteManualContactsContext(ctx context.Context) (err error)
	// DeleteManualContactsBatchContext : Removes manually added contacts from
	// the given list.
	DeleteManualContactsBatchContext(ctx context.Context, arg *DeleteManualContactsArg) (err error)
}

type apiImpl dropbox.Context

// DeleteManualContactsAPIError is an error-wrapper for the delete_manual_contacts route
type DeleteManualContactsAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// DeleteManualContactsContext : Removes all manually added contacts. You'll
// still keep contacts who are on your team or who you imported. New contacts
// will be added when you share.
func (dbx *apiImpl) DeleteManualContactsContext(ctx context.Context) (err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr DeleteManualContactsAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

func (dbx *apiImpl) DeleteManualContacts() (err error) {
	return dbx.DeleteManualContactsContext(context.Background())
}

// DeleteManualContactsBatchAPIError is an error-wrapper for the delete_manual_contacts_batch route
type DeleteManualContactsBatchAPIError struct {
	dropbox.APIError
	EndpointError *DeleteManualContactsError `json:"error"`
}

// DeleteManualContactsBatchContext : Removes manually added contacts from the
// given list.
func (dbx *apiImpl) DeleteManualContactsBatchContext(ctx context.Context, arg *DeleteManualContactsArg) (err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr DeleteManualContactsBatchAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

func (dbx *apiImpl) DeleteManualContactsBatch(arg *DeleteManualContactsArg) (err error) {
	return dbx.DeleteManualContactsBatchContext(context.Background(), arg)
}

// NewContext returns a ContextClient implementation for this namespace
func NewContext(c dropbox.Config) ContextClient {
	ctx := apiImpl(dropbox.NewContext(c))
	return &ctx
}

// New returns a Client implementation for this namespace
func New(c dropbox.Config) Client {
	return NewContext(c)
}
