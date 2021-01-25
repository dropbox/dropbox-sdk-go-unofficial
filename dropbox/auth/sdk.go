package auth

import (
	"encoding/json"
	"net/http"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
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

func ParseError(err error, appError interface{}) error {
	sdkErr, ok := err.(dropbox.SDKInternalError)
	if !ok {
		return err
	}

	switch sdkErr.StatusCode {
	case http.StatusBadRequest:
		return BadRequest{
			APIError: dropbox.APIError{
				ErrorSummary: sdkErr.Content,
			},
		}
	case http.StatusInternalServerError:
		return ServerError{
			APIError: dropbox.APIError{
				ErrorSummary: sdkErr.Content,
			},
		}
	case http.StatusUnauthorized:
		var apiError AuthAPIError
		if pErr := json.Unmarshal([]byte(sdkErr.Content), &apiError); pErr != nil {
			return pErr
		}

		return apiError
	case http.StatusForbidden:
		var apiError AccessAPIError
		if pErr := json.Unmarshal([]byte(sdkErr.Content), &apiError); pErr != nil {
			return pErr
		}

		return apiError
	case http.StatusTooManyRequests:
		var apiError RateLimitAPIError
		if pErr := json.Unmarshal([]byte(sdkErr.Content), &apiError); pErr != nil {
			return pErr
		}

		return apiError
	}

	return err
}
