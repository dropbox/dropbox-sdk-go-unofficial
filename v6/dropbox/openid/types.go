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

// Package openid : has no documentation (yet)
package openid

import (
	"encoding/json"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
)

// OpenIdError : has no documentation (yet)
type OpenIdError struct {
	dropbox.Tagged
}

// Valid tag values for OpenIdError
const (
	OpenIdErrorIncorrectOpenidScopes = "incorrect_openid_scopes"
	OpenIdErrorOther                 = "other"
)

// UserInfoArgs : No Parameters
type UserInfoArgs struct {
}

// NewUserInfoArgs returns a new UserInfoArgs instance
func NewUserInfoArgs() *UserInfoArgs {
	s := new(UserInfoArgs)
	return s
}

// UserInfoError : has no documentation (yet)
type UserInfoError struct {
	dropbox.Tagged
	// OpenidError : has no documentation (yet)
	OpenidError *OpenIdError `json:"openid_error,omitempty"`
}

// Valid tag values for UserInfoError
const (
	UserInfoErrorOpenidError = "openid_error"
	UserInfoErrorOther       = "other"
)

// UnmarshalJSON deserializes into a UserInfoError instance
func (u *UserInfoError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// OpenidError : has no documentation (yet)
		OpenidError *OpenIdError `json:"openid_error,omitempty"`
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "openid_error":
		u.OpenidError = w.OpenidError

	}
	return nil
}

// UserInfoResult : has no documentation (yet)
type UserInfoResult struct {
	// FamilyName : Last name of user.
	FamilyName string `json:"family_name,omitempty"`
	// GivenName : First name of user.
	GivenName string `json:"given_name,omitempty"`
	// Email : Email address of user.
	Email string `json:"email,omitempty"`
	// EmailVerified : If user is email verified.
	EmailVerified bool `json:"email_verified,omitempty"`
	// Iss : Issuer of token (in this case Dropbox).
	Iss string `json:"iss"`
	// Sub : An identifier for the user. This is the Dropbox account_id, a
	// string value such as dbid:AAH4f99T0taONIb-OurWxbNQ6ywGRopQngc.
	Sub string `json:"sub"`
}

// NewUserInfoResult returns a new UserInfoResult instance
func NewUserInfoResult() *UserInfoResult {
	s := new(UserInfoResult)
	s.Iss = ""
	s.Sub = ""
	return s
}
