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

package riviera

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/async"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// GetMarkdownAsync : Asynchronous document-to-markdown conversion for
	// supported file formats.
	GetMarkdownAsync(arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error)
	// GetMarkdownAsyncCheck : Returns the status or result of specified
	// get_markdown_async task.
	GetMarkdownAsyncCheck(arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error)
	// GetTranscriptAsync : Asynchronous transcript generation for audio and
	// video files.
	GetTranscriptAsync(arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error)
	// GetTranscriptAsyncCheck : Returns the status or result of specified
	// get_transcript_async task.
	GetTranscriptAsyncCheck(arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error)
}

// ContextClient interface describes all routes in this namespace with context support
type ContextClient interface {
	Client
	// GetMarkdownAsyncContext : Asynchronous document-to-markdown conversion for
	// supported file formats.
	GetMarkdownAsyncContext(ctx context.Context, arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error)
	// GetMarkdownAsyncCheckContext : Returns the status or result of specified
	// get_markdown_async task.
	GetMarkdownAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error)
	// GetTranscriptAsyncContext : Asynchronous transcript generation for audio and
	// video files.
	GetTranscriptAsyncContext(ctx context.Context, arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error)
	// GetTranscriptAsyncCheckContext : Returns the status or result of specified
	// get_transcript_async task.
	GetTranscriptAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error)
}

type apiImpl dropbox.Context

// GetMarkdownAsyncAPIError is an error-wrapper for the get_markdown_async route
type GetMarkdownAsyncAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// GetMarkdownAsyncContext : Asynchronous document-to-markdown conversion for
// supported file formats.
func (dbx *apiImpl) GetMarkdownAsyncContext(ctx context.Context, arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_markdown_async",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetMarkdownAsyncAPIError
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

func (dbx *apiImpl) GetMarkdownAsync(arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error) {
	return dbx.GetMarkdownAsyncContext(context.Background(), arg)
}

// GetMarkdownAsyncCheckAPIError is an error-wrapper for the get_markdown_async/check route
type GetMarkdownAsyncCheckAPIError struct {
	dropbox.APIError
	EndpointError *async.PollError `json:"error"`
}

// GetMarkdownAsyncCheckContext : Returns the status or result of specified
// get_markdown_async task.
func (dbx *apiImpl) GetMarkdownAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_markdown_async/check",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetMarkdownAsyncCheckAPIError
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

func (dbx *apiImpl) GetMarkdownAsyncCheck(arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error) {
	return dbx.GetMarkdownAsyncCheckContext(context.Background(), arg)
}

// GetTranscriptAsyncAPIError is an error-wrapper for the get_transcript_async route
type GetTranscriptAsyncAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// GetTranscriptAsyncContext : Asynchronous transcript generation for audio and
// video files.
func (dbx *apiImpl) GetTranscriptAsyncContext(ctx context.Context, arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_transcript_async",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetTranscriptAsyncAPIError
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

func (dbx *apiImpl) GetTranscriptAsync(arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error) {
	return dbx.GetTranscriptAsyncContext(context.Background(), arg)
}

// GetTranscriptAsyncCheckAPIError is an error-wrapper for the get_transcript_async/check route
type GetTranscriptAsyncCheckAPIError struct {
	dropbox.APIError
	EndpointError *async.PollError `json:"error"`
}

// GetTranscriptAsyncCheckContext : Returns the status or result of specified
// get_transcript_async task.
func (dbx *apiImpl) GetTranscriptAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_transcript_async/check",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetTranscriptAsyncCheckAPIError
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

func (dbx *apiImpl) GetTranscriptAsyncCheck(arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error) {
	return dbx.GetTranscriptAsyncCheckContext(context.Background(), arg)
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
