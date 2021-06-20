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

package file_requests

import (
	"encoding/json"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// Count : Returns the total number of file requests owned by this user.
	// Includes both open and closed file requests.
	Count() (res *CountFileRequestsResult, err error)
	// Create : Creates a file request for this user.
	Create(arg *CreateFileRequestArgs) (res *FileRequest, err error)
	// Delete : Delete a batch of closed file requests.
	Delete(arg *DeleteFileRequestArgs) (res *DeleteFileRequestsResult, err error)
	// DeleteAllClosed : Delete all closed file requests owned by this user.
	DeleteAllClosed() (res *DeleteAllClosedFileRequestsResult, err error)
	// Get : Returns the specified file request.
	Get(arg *GetFileRequestArgs) (res *FileRequest, err error)
	// List : Returns a list of file requests owned by this user. For apps with
	// the app folder permission, this will only return file requests with
	// destinations in the app folder.
	ListV2(arg *ListFileRequestsArg) (res *ListFileRequestsV2Result, err error)
	// List : Returns a list of file requests owned by this user. For apps with
	// the app folder permission, this will only return file requests with
	// destinations in the app folder.
	List() (res *ListFileRequestsResult, err error)
	// ListContinue : Once a cursor has been retrieved from `list`, use this to
	// paginate through all file requests. The cursor must come from a previous
	// call to `list` or `listContinue`.
	ListContinue(arg *ListFileRequestsContinueArg) (res *ListFileRequestsV2Result, err error)
	// Update : Update a file request.
	Update(arg *UpdateFileRequestArgs) (res *FileRequest, err error)
}

type apiImpl dropbox.Context

//CountAPIError is an error-wrapper for the count route
type CountAPIError struct {
	dropbox.APIError
	EndpointError *CountFileRequestsError `json:"error"`
}

func (dbx *apiImpl) Count() (res *CountFileRequestsResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "count",
		Auth:         "user",
		Style:        "rpc",
		Arg:          nil,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr CountAPIError
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

//CreateAPIError is an error-wrapper for the create route
type CreateAPIError struct {
	dropbox.APIError
	EndpointError *CreateFileRequestError `json:"error"`
}

func (dbx *apiImpl) Create(arg *CreateFileRequestArgs) (res *FileRequest, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "create",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr CreateAPIError
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

//DeleteAPIError is an error-wrapper for the delete route
type DeleteAPIError struct {
	dropbox.APIError
	EndpointError *DeleteFileRequestError `json:"error"`
}

func (dbx *apiImpl) Delete(arg *DeleteFileRequestArgs) (res *DeleteFileRequestsResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "delete",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DeleteAPIError
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

//DeleteAllClosedAPIError is an error-wrapper for the delete_all_closed route
type DeleteAllClosedAPIError struct {
	dropbox.APIError
	EndpointError *DeleteAllClosedFileRequestsError `json:"error"`
}

func (dbx *apiImpl) DeleteAllClosed() (res *DeleteAllClosedFileRequestsResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "delete_all_closed",
		Auth:         "user",
		Style:        "rpc",
		Arg:          nil,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DeleteAllClosedAPIError
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

//GetAPIError is an error-wrapper for the get route
type GetAPIError struct {
	dropbox.APIError
	EndpointError *GetFileRequestError `json:"error"`
}

func (dbx *apiImpl) Get(arg *GetFileRequestArgs) (res *FileRequest, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "get",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr GetAPIError
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

//ListV2APIError is an error-wrapper for the list_v2 route
type ListV2APIError struct {
	dropbox.APIError
	EndpointError *ListFileRequestsError `json:"error"`
}

func (dbx *apiImpl) ListV2(arg *ListFileRequestsArg) (res *ListFileRequestsV2Result, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "list_v2",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr ListV2APIError
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

//ListAPIError is an error-wrapper for the list route
type ListAPIError struct {
	dropbox.APIError
	EndpointError *ListFileRequestsError `json:"error"`
}

func (dbx *apiImpl) List() (res *ListFileRequestsResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "list",
		Auth:         "user",
		Style:        "rpc",
		Arg:          nil,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr ListAPIError
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

//ListContinueAPIError is an error-wrapper for the list/continue route
type ListContinueAPIError struct {
	dropbox.APIError
	EndpointError *ListFileRequestsContinueError `json:"error"`
}

func (dbx *apiImpl) ListContinue(arg *ListFileRequestsContinueArg) (res *ListFileRequestsV2Result, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "list/continue",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr ListContinueAPIError
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

//UpdateAPIError is an error-wrapper for the update route
type UpdateAPIError struct {
	dropbox.APIError
	EndpointError *UpdateFileRequestError `json:"error"`
}

func (dbx *apiImpl) Update(arg *UpdateFileRequestArgs) (res *FileRequest, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "file_requests",
		Route:        "update",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr UpdateAPIError
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
