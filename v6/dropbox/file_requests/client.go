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
	"context"
	"encoding/json"
	"errors"
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
	List() (res *ListFileRequestsResult, err error)
	// List : Returns a list of file requests owned by this user. For apps with
	// the app folder permission, this will only return file requests with
	// destinations in the app folder.
	ListV2(arg *ListFileRequestsArg) (res *ListFileRequestsV2Result, err error)
	// ListContinue : Once a cursor has been retrieved from `list`, use this to
	// paginate through all file requests. The cursor must come from a previous
	// call to `list` or `listContinue`.
	ListContinue(arg *ListFileRequestsContinueArg) (res *ListFileRequestsV2Result, err error)
	// Update : Update a file request.
	Update(arg *UpdateFileRequestArgs) (res *FileRequest, err error)
}

// ContextClient interface describes all routes in this namespace with context support
type ContextClient interface {
	Client
	// CountContext : Returns the total number of file requests owned by this
	// user. Includes both open and closed file requests.
	CountContext(ctx context.Context) (res *CountFileRequestsResult, err error)
	// CreateContext : Creates a file request for this user.
	CreateContext(ctx context.Context, arg *CreateFileRequestArgs) (res *FileRequest, err error)
	// DeleteContext : Delete a batch of closed file requests.
	DeleteContext(ctx context.Context, arg *DeleteFileRequestArgs) (res *DeleteFileRequestsResult, err error)
	// DeleteAllClosedContext : Delete all closed file requests owned by this
	// user.
	DeleteAllClosedContext(ctx context.Context) (res *DeleteAllClosedFileRequestsResult, err error)
	// GetContext : Returns the specified file request.
	GetContext(ctx context.Context, arg *GetFileRequestArgs) (res *FileRequest, err error)
	// ListContext : Returns a list of file requests owned by this user. For
	// apps with the app folder permission, this will only return file requests
	// with destinations in the app folder.
	ListContext(ctx context.Context) (res *ListFileRequestsResult, err error)
	// ListV2Context : Returns a list of file requests owned by this user. For
	// apps with the app folder permission, this will only return file requests
	// with destinations in the app folder.
	ListV2Context(ctx context.Context, arg *ListFileRequestsArg) (res *ListFileRequestsV2Result, err error)
	// ListContinueContext : Once a cursor has been retrieved from `list`, use
	// this to paginate through all file requests. The cursor must come from a
	// previous call to `list` or `listContinue`.
	ListContinueContext(ctx context.Context, arg *ListFileRequestsContinueArg) (res *ListFileRequestsV2Result, err error)
	// UpdateContext : Update a file request.
	UpdateContext(ctx context.Context, arg *UpdateFileRequestArgs) (res *FileRequest, err error)
}

type apiImpl dropbox.Context

// CountAPIError is an error-wrapper for the count route
type CountAPIError struct {
	dropbox.APIError
	EndpointError *CountFileRequestsError `json:"error"`
}

// CountContext : Returns the total number of file requests owned by this user.
// Includes both open and closed file requests.
func (dbx *apiImpl) CountContext(ctx context.Context) (res *CountFileRequestsResult, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr CountAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) Count() (res *CountFileRequestsResult, err error) {
	return dbx.CountContext(context.Background())
}

// CreateAPIError is an error-wrapper for the create route
type CreateAPIError struct {
	dropbox.APIError
	EndpointError *CreateFileRequestError `json:"error"`
}

// CreateContext : Creates a file request for this user.
func (dbx *apiImpl) CreateContext(ctx context.Context, arg *CreateFileRequestArgs) (res *FileRequest, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr CreateAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) Create(arg *CreateFileRequestArgs) (res *FileRequest, err error) {
	return dbx.CreateContext(context.Background(), arg)
}

// DeleteAPIError is an error-wrapper for the delete route
type DeleteAPIError struct {
	dropbox.APIError
	EndpointError *DeleteFileRequestError `json:"error"`
}

// DeleteContext : Delete a batch of closed file requests.
func (dbx *apiImpl) DeleteContext(ctx context.Context, arg *DeleteFileRequestArgs) (res *DeleteFileRequestsResult, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr DeleteAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) Delete(arg *DeleteFileRequestArgs) (res *DeleteFileRequestsResult, err error) {
	return dbx.DeleteContext(context.Background(), arg)
}

// DeleteAllClosedAPIError is an error-wrapper for the delete_all_closed route
type DeleteAllClosedAPIError struct {
	dropbox.APIError
	EndpointError *DeleteAllClosedFileRequestsError `json:"error"`
}

// DeleteAllClosedContext : Delete all closed file requests owned by this user.
func (dbx *apiImpl) DeleteAllClosedContext(ctx context.Context) (res *DeleteAllClosedFileRequestsResult, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr DeleteAllClosedAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) DeleteAllClosed() (res *DeleteAllClosedFileRequestsResult, err error) {
	return dbx.DeleteAllClosedContext(context.Background())
}

// GetAPIError is an error-wrapper for the get route
type GetAPIError struct {
	dropbox.APIError
	EndpointError *GetFileRequestError `json:"error"`
}

// GetContext : Returns the specified file request.
func (dbx *apiImpl) GetContext(ctx context.Context, arg *GetFileRequestArgs) (res *FileRequest, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) Get(arg *GetFileRequestArgs) (res *FileRequest, err error) {
	return dbx.GetContext(context.Background(), arg)
}

// ListAPIError is an error-wrapper for the list route
type ListAPIError struct {
	dropbox.APIError
	EndpointError *ListFileRequestsError `json:"error"`
}

// ListContext : Returns a list of file requests owned by this user. For apps
// with the app folder permission, this will only return file requests with
// destinations in the app folder.
func (dbx *apiImpl) ListContext(ctx context.Context) (res *ListFileRequestsResult, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr ListAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) List() (res *ListFileRequestsResult, err error) {
	return dbx.ListContext(context.Background())
}

// ListV2APIError is an error-wrapper for the list_v2 route
type ListV2APIError struct {
	dropbox.APIError
	EndpointError *ListFileRequestsError `json:"error"`
}

// ListV2Context : Returns a list of file requests owned by this user. For apps
// with the app folder permission, this will only return file requests with
// destinations in the app folder.
func (dbx *apiImpl) ListV2Context(ctx context.Context, arg *ListFileRequestsArg) (res *ListFileRequestsV2Result, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr ListV2APIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) ListV2(arg *ListFileRequestsArg) (res *ListFileRequestsV2Result, err error) {
	return dbx.ListV2Context(context.Background(), arg)
}

// ListContinueAPIError is an error-wrapper for the list/continue route
type ListContinueAPIError struct {
	dropbox.APIError
	EndpointError *ListFileRequestsContinueError `json:"error"`
}

// ListContinueContext : Once a cursor has been retrieved from `list`, use this
// to paginate through all file requests. The cursor must come from a previous
// call to `list` or `listContinue`.
func (dbx *apiImpl) ListContinueContext(ctx context.Context, arg *ListFileRequestsContinueArg) (res *ListFileRequestsV2Result, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr ListContinueAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) ListContinue(arg *ListFileRequestsContinueArg) (res *ListFileRequestsV2Result, err error) {
	return dbx.ListContinueContext(context.Background(), arg)
}

// UpdateAPIError is an error-wrapper for the update route
type UpdateAPIError struct {
	dropbox.APIError
	EndpointError *UpdateFileRequestError `json:"error"`
}

// UpdateContext : Update a file request.
func (dbx *apiImpl) UpdateContext(ctx context.Context, arg *UpdateFileRequestArgs) (res *FileRequest, err error) {
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
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr UpdateAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) Update(arg *UpdateFileRequestArgs) (res *FileRequest, err error) {
	return dbx.UpdateContext(context.Background(), arg)
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
