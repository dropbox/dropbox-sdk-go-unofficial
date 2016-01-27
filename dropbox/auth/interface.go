/* DO NOT EDIT */
/* This file was generated from auth.babel */

package auth

import "encoding/json"

// Errors occurred during authentication.
type AuthError struct {
	Tag string `json:".tag"`
}

func (u *AuthError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	}
	return nil
}
