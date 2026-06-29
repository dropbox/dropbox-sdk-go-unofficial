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

package account

import (
	"encoding/json"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// DeleteProfilePhoto : Deletes the current user's profile photo.
	DeleteProfilePhoto(arg *DeleteProfilePhotoArg) (res *DeleteProfilePhotoResult, err error)
	// GetPhoto : This lovely endpoint gets the account photo of a given user.
	GetPhoto(arg *AccountPhotoGetArg) (res *AccountPhotoGetResult, content io.ReadCloser, err error)
	// SetProfilePhoto : Sets a user's profile photo.
	SetProfilePhoto(arg *SetProfilePhotoArg) (res *SetProfilePhotoResult, err error)
}

type apiImpl dropbox.Context

// DeleteProfilePhotoAPIError is an error-wrapper for the delete_profile_photo route
type DeleteProfilePhotoAPIError struct {
	dropbox.APIError
	EndpointError *DeleteProfilePhotoError `json:"error"`
}

func (dbx *apiImpl) DeleteProfilePhoto(arg *DeleteProfilePhotoArg) (res *DeleteProfilePhotoResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "account",
		Route:        "delete_profile_photo",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DeleteProfilePhotoAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return
	}

	_ = respBody
	return
}

// GetPhotoAPIError is an error-wrapper for the get_photo route
type GetPhotoAPIError struct {
	dropbox.APIError
	EndpointError *AccountPhotoGetError `json:"error"`
}

func (dbx *apiImpl) GetPhoto(arg *AccountPhotoGetArg) (res *AccountPhotoGetResult, content io.ReadCloser, err error) {
	req := dropbox.Request{
		Host:         "content",
		Namespace:    "account",
		Route:        "get_photo",
		Auth:         "user",
		Style:        "download",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr GetPhotoAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return
	}

	content = respBody
	return
}

// SetProfilePhotoAPIError is an error-wrapper for the set_profile_photo route
type SetProfilePhotoAPIError struct {
	dropbox.APIError
	EndpointError *SetProfilePhotoError `json:"error"`
}

func (dbx *apiImpl) SetProfilePhoto(arg *SetProfilePhotoArg) (res *SetProfilePhotoResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "account",
		Route:        "set_profile_photo",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr SetProfilePhotoAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return
	}

	_ = respBody
	return
}

// New returns a Client implementation for this namespace
func New(c dropbox.Config) Client {
	ctx := apiImpl(dropbox.NewContext(c))
	return &ctx
}
