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

// ParseError decodes a transport error into the endpoint error type T.
// Generated callers pass a pointer to a zero value; decoded endpoint errors
// are returned by value, including when wrapped in HTTPError.
func ParseError[T error](err error, appError *T) (parsed error) {
	if httpErr, ok := err.(*dropbox.HTTPError); ok {
		err = httpErr.Err
		defer func() {
			httpErr.Err = parsed
			parsed = httpErr
		}()
	}
	sdkErr, ok := err.(dropbox.SDKInternalError)
	if !ok {
		return err
	}

	if sdkErr.StatusCode >= 500 && sdkErr.StatusCode <= 599 {
		return ServerError{
			StatusCode: sdkErr.StatusCode,
			APIError: dropbox.APIError{
				ErrorSummary: sdkErr.Content,
			},
		}
	}

	switch sdkErr.StatusCode {
	case http.StatusBadRequest:
		return BadRequest{
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
	case http.StatusConflict:
		if pErr := json.Unmarshal([]byte(sdkErr.Content), appError); pErr != nil {
			return pErr
		}

		return *appError
	}

	return err
}
