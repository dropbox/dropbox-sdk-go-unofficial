/* DO NOT EDIT */
/* This file was generated from async.babel */

package dropbox

import "encoding/json"

// Result returned by methods that launch an asynchronous job. A method who may
// either launch an asynchronous job, or complete the request synchronously, can
// use this union by extending it, and adding a 'complete' field with the type
// of the synchronous response. See :type:`LaunchEmptyResult` for an example.
type LaunchResultBase struct {
	Tag string `json:".tag"`
	// This response indicates that the processing is asynchronous. The string is
	// an id that can be used to obtain the status of the asynchronous job.
	AsyncJobId string `json:"async_job_id,omitempty"`
}

func (u *LaunchResultBase) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// This response indicates that the processing is asynchronous. The string is
		// an id that can be used to obtain the status of the asynchronous job.
		AsyncJobId json.RawMessage `json:"async_job_id"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "async_job_id":
		{
			if len(w.AsyncJobId) == 0 {
				break
			}
			if err := json.Unmarshal(w.AsyncJobId, &u.AsyncJobId); err != nil {
				return err
			}
		}
	}
	return nil
}

// Result returned by methods that may either launch an asynchronous job or
// complete synchronously. Upon synchronous completion of the job, no additional
// information is returned.
type LaunchEmptyResult struct {
	Tag string `json:".tag"`
}

func (u *LaunchEmptyResult) UnmarshalJSON(body []byte) error {
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

// Arguments for methods that poll the status of an asynchronous job.
type PollArg struct {
	// Id of the asynchronous job. This is the value of a response returned from
	// the method that launched the job.
	AsyncJobId string `json:"async_job_id"`
}

func NewPollArg() *PollArg {
	s := new(PollArg)
	return s
}

// Result returned by methods that poll for the status of an asynchronous job.
// Unions that extend this union should add a 'complete' field with a type of
// the information returned upon job completion. See :type:`PollEmptyResult` for
// an example.
type PollResultBase struct {
	Tag string `json:".tag"`
}

func (u *PollResultBase) UnmarshalJSON(body []byte) error {
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

// Result returned by methods that poll for the status of an asynchronous job.
// Upon completion of the job, no additional information is returned.
type PollEmptyResult struct {
	Tag string `json:".tag"`
}

func (u *PollEmptyResult) UnmarshalJSON(body []byte) error {
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

// Error returned by methods for polling the status of asynchronous job.
type PollError struct {
	Tag string `json:".tag"`
}

func (u *PollError) UnmarshalJSON(body []byte) error {
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
