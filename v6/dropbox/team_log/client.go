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

package team_log

import (
	"encoding/json"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// GetEvents : Retrieves team events. If the result's
	// `GetTeamEventsResult.has_more` field is true, call `getEventsContinue`
	// with the returned cursor to retrieve more entries. If end_time is not
	// specified in your request, you may use the returned cursor to poll
	// `getEventsContinue` for new events. Many attributes note 'may be missing
	// due to historical data gap'. Note that the file_operations category and &
	// analogous paper events are not available on all Dropbox Business `plans`
	// </business/plans-comparison>. Use `features/get_values`
	// </developers/documentation/http/teams#team-features-get_values> to check
	// for this feature. Permission : Team Auditing.
	GetEvents(arg *GetTeamEventsArg) (res *GetTeamEventsResult, err error)
	// GetEventsContinue : Once a cursor has been retrieved from `getEvents`,
	// use this to paginate through all events. Permission : Team Auditing.
	GetEventsContinue(arg *GetTeamEventsContinueArg) (res *GetTeamEventsResult, err error)
}

type apiImpl dropbox.Context

//GetEventsAPIError is an error-wrapper for the get_events route
type GetEventsAPIError struct {
	dropbox.APIError
	EndpointError *GetTeamEventsError `json:"error"`
}

func (dbx *apiImpl) GetEvents(arg *GetTeamEventsArg) (res *GetTeamEventsResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "team_log",
		Route:        "get_events",
		Auth:         "team",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr GetEventsAPIError
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

//GetEventsContinueAPIError is an error-wrapper for the get_events/continue route
type GetEventsContinueAPIError struct {
	dropbox.APIError
	EndpointError *GetTeamEventsContinueError `json:"error"`
}

func (dbx *apiImpl) GetEventsContinue(arg *GetTeamEventsContinueArg) (res *GetTeamEventsResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "team_log",
		Route:        "get_events/continue",
		Auth:         "team",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr GetEventsContinueAPIError
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
