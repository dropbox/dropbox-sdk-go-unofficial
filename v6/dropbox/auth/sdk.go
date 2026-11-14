package auth

import (
	"encoding/json"
	"net/http"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
)

// AuthAPIError wraps AuthError
type AuthAPIError struct {
	dropbox.APIError
	AuthError *AuthError `json:"error"`
}

// AccessAPIError wraps AccessError
type AccessAPIError struct {
	dropbox.APIError
	AccessError *AccessError `json:"error"`
}

// RateLimitAPIError wraps RateLimitError
type RateLimitAPIError struct {
	dropbox.APIError
	RateLimitError *RateLimitError `json:"error"`
}

// Bad input parameter.
type BadRequest struct {
	dropbox.APIError
}

// An error occurred on the Dropbox servers. Check status.dropbox.com for announcements about
// Dropbox service issues.
type ServerError struct {
	dropbox.APIError
	StatusCode int
}

func ParseError(err error, appError error) error {
	sdkErr, ok := err.(dropbox.SDKInternalError)
	if !ok {
		return err
	}

	base := dropbox.APIError{ErrorSummary: sdkErr.Content}
	base.SetResponseMetadata(sdkErr.ResponseMetadata())

	if sdkErr.StatusCode >= 500 && sdkErr.StatusCode <= 599 {
		return ServerError{
			StatusCode: sdkErr.StatusCode,
			APIError:   base,
		}
	}

	switch sdkErr.StatusCode {
	case http.StatusBadRequest:
		return BadRequest{
			APIError: base,
		}
	case http.StatusUnauthorized:
		var apiError AuthAPIError
		if pErr := json.Unmarshal([]byte(sdkErr.Content), &apiError); pErr != nil {
			return dropbox.NewResponseError(pErr, sdkErr.ResponseMetadata())
		}

		apiError.SetResponseMetadata(sdkErr.ResponseMetadata())
		return apiError
	case http.StatusForbidden:
		var apiError AccessAPIError
		if pErr := json.Unmarshal([]byte(sdkErr.Content), &apiError); pErr != nil {
			return dropbox.NewResponseError(pErr, sdkErr.ResponseMetadata())
		}

		apiError.SetResponseMetadata(sdkErr.ResponseMetadata())
		return apiError
	case http.StatusTooManyRequests:
		var apiError RateLimitAPIError
		if pErr := json.Unmarshal([]byte(sdkErr.Content), &apiError); pErr != nil {
			return dropbox.NewResponseError(pErr, sdkErr.ResponseMetadata())
		}

		apiError.SetResponseMetadata(sdkErr.ResponseMetadata())
		return apiError
	case http.StatusConflict:
		if pErr := json.Unmarshal([]byte(sdkErr.Content), appError); pErr != nil {
			return dropbox.NewResponseError(pErr, sdkErr.ResponseMetadata())
		}

		// Generated callers copy appError to a value after ParseError returns.
		// Set the embedded APIError before that copy, without changing its type.
		if target, ok := appError.(interface {
			SetResponseMetadata(*dropbox.ResponseMetadata)
		}); ok {
			target.SetResponseMetadata(sdkErr.ResponseMetadata())
		}
		return appError
	}

	return err
}
