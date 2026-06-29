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

// Package riviera : has no documentation (yet)
package riviera

import (
	"encoding/json"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
)

// ApiStructuredTranscript : Structured transcript for APIv2
type ApiStructuredTranscript struct {
	// Segments : has no documentation (yet)
	Segments []*ApiTranscriptSegment `json:"segments,omitempty"`
	// TranscriptLocale : has no documentation (yet)
	TranscriptLocale string `json:"transcript_locale"`
}

// NewApiStructuredTranscript returns a new ApiStructuredTranscript instance
func NewApiStructuredTranscript() *ApiStructuredTranscript {
	s := new(ApiStructuredTranscript)
	s.TranscriptLocale = ""
	return s
}

// ApiTranscriptSegment : Transcript segment for APIv2
type ApiTranscriptSegment struct {
	// Text : has no documentation (yet)
	Text string `json:"text"`
	// StartTime : has no documentation (yet)
	StartTime float64 `json:"start_time"`
	// EndTime : has no documentation (yet)
	EndTime float64 `json:"end_time"`
}

// NewApiTranscriptSegment returns a new ApiTranscriptSegment instance
func NewApiTranscriptSegment() *ApiTranscriptSegment {
	s := new(ApiTranscriptSegment)
	s.Text = ""
	s.StartTime = 0.0
	s.EndTime = 0.0
	return s
}

// ContentApiV2Error : has no documentation (yet)
type ContentApiV2Error struct {
	dropbox.Tagged
	// ServerError : has no documentation (yet)
	ServerError string `json:"server_error,omitempty"`
	// UserError : has no documentation (yet)
	UserError string `json:"user_error,omitempty"`
	// MediaDurationError : has no documentation (yet)
	MediaDurationError *MediaDurationError `json:"media_duration_error,omitempty"`
}

// Valid tag values for ContentApiV2Error
const (
	ContentApiV2ErrorServerError                 = "server_error"
	ContentApiV2ErrorUserError                   = "user_error"
	ContentApiV2ErrorMediaDurationError          = "media_duration_error"
	ContentApiV2ErrorNoAudioError                = "no_audio_error"
	ContentApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	ContentApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	ContentApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	ContentApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a ContentApiV2Error instance
func (u *ContentApiV2Error) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// ServerError : has no documentation (yet)
		ServerError string `json:"server_error,omitempty"`
		// UserError : has no documentation (yet)
		UserError string `json:"user_error,omitempty"`
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "server_error":
		u.ServerError = w.ServerError

	case "user_error":
		u.UserError = w.UserError

	case "media_duration_error":
		if err = json.Unmarshal(body, &u.MediaDurationError); err != nil {
			return err
		}

	}
	return nil
}

// ErrorCode : has no documentation (yet)
type ErrorCode struct {
	dropbox.Tagged
}

// Valid tag values for ErrorCode
const (
	ErrorCodeUnknownError   = "unknown_error"
	ErrorCodeBadRequest     = "bad_request"
	ErrorCodeApiError       = "api_error"
	ErrorCodeAccessError    = "access_error"
	ErrorCodeRatelimitError = "ratelimit_error"
	ErrorCodeUnavailable    = "unavailable"
	ErrorCodeOther          = "other"
)

// FileIdOrUrl : has no documentation (yet)
type FileIdOrUrl struct {
	dropbox.Tagged
	// FileId : has no documentation (yet)
	FileId string `json:"file_id,omitempty"`
	// Url : has no documentation (yet)
	Url string `json:"url,omitempty"`
	// Path : has no documentation (yet)
	Path string `json:"path,omitempty"`
}

// Valid tag values for FileIdOrUrl
const (
	FileIdOrUrlFileId = "file_id"
	FileIdOrUrlUrl    = "url"
	FileIdOrUrlPath   = "path"
	FileIdOrUrlOther  = "other"
)

// UnmarshalJSON deserializes into a FileIdOrUrl instance
func (u *FileIdOrUrl) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// FileId : has no documentation (yet)
		FileId string `json:"file_id,omitempty"`
		// Url : has no documentation (yet)
		Url string `json:"url,omitempty"`
		// Path : has no documentation (yet)
		Path string `json:"path,omitempty"`
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "file_id":
		u.FileId = w.FileId

	case "url":
		u.Url = w.Url

	case "path":
		u.Path = w.Path

	}
	return nil
}

// GetMarkdownArgs : Arguments for the asynchronous `get_markdown_async` route.
// Exactly one of `file_id`, `path`, or `url` must be supplied via
// `file_id_or_url` to identify the document to convert to markdown.
type GetMarkdownArgs struct {
	// FileIdOrUrl : Identifier of the document to convert. Callers must set
	// exactly one of the oneof variants: - file_id: a Dropbox-issued file id
	// (format: "id:<id>") for a file the authenticated user has access to. -
	// path: an absolute Dropbox path, e.g. "/folder/report.docx". - url: either
	// a Dropbox shared link (www.dropbox.com) or an external HTTPS URL pointing
	// to a supported document file. - Dropbox shared links are resolved
	// internally using the caller's authenticated identity and the link's
	// visibility / download settings. They therefore require an authenticated
	// user context (anonymous `url` requests against Dropbox links are rejected
	// with an `ACCESS_ERROR`). Links protected by a password are rejected with
	// `shared_link_password_protected`; links with downloads disabled are
	// rejected with `link_download_disabled_error`. - External URLs are fetched
	// over HTTPS through the backend's egress proxy and must point at a
	// supported document file extension. The referenced file must be a document
	// in a supported format; requests against unsupported formats return
	// `unsupported_format_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
	// EnableOcr : Enable OCR for PDF documents. Processing is slower when
	// enabled.
	EnableOcr bool `json:"enable_ocr"`
	// EmbedImages : When true, embed images as base64 data URIs in the markdown
	// output. This can significantly increase output size.
	EmbedImages bool `json:"embed_images"`
}

// NewGetMarkdownArgs returns a new GetMarkdownArgs instance
func NewGetMarkdownArgs() *GetMarkdownArgs {
	s := new(GetMarkdownArgs)
	s.EnableOcr = false
	s.EmbedImages = false
	return s
}

// GetMarkdownAsyncCheckResult : Result type for EventBus async check
type GetMarkdownAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : has no documentation (yet)
	Complete *GetMarkdownResult `json:"complete,omitempty"`
	// Failed : has no documentation (yet)
	Failed *GetMarkdownAsyncError `json:"failed,omitempty"`
}

// Valid tag values for GetMarkdownAsyncCheckResult
const (
	GetMarkdownAsyncCheckResultInProgress = "in_progress"
	GetMarkdownAsyncCheckResultComplete   = "complete"
	GetMarkdownAsyncCheckResultFailed     = "failed"
	GetMarkdownAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetMarkdownAsyncCheckResult instance
func (u *GetMarkdownAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "complete":
		if err = json.Unmarshal(body, &u.Complete); err != nil {
			return err
		}

	case "failed":
		if err = json.Unmarshal(body, &u.Failed); err != nil {
			return err
		}

	}
	return nil
}

// GetMarkdownAsyncError : has no documentation (yet)
type GetMarkdownAsyncError struct {
	// ErrorCode : has no documentation (yet)
	ErrorCode *ErrorCode `json:"error_code"`
	// ErrorDetails : has no documentation (yet)
	ErrorDetails *MarkdownConversionApiV2Error `json:"error_details,omitempty"`
}

// NewGetMarkdownAsyncError returns a new GetMarkdownAsyncError instance
func NewGetMarkdownAsyncError() *GetMarkdownAsyncError {
	s := new(GetMarkdownAsyncError)
	s.ErrorCode = &ErrorCode{Tagged: dropbox.Tagged{Tag: "unknown_error"}}
	return s
}

// GetMarkdownResult : has no documentation (yet)
type GetMarkdownResult struct {
	// Markdown : The converted markdown content
	Markdown string `json:"markdown"`
}

// NewGetMarkdownResult returns a new GetMarkdownResult instance
func NewGetMarkdownResult() *GetMarkdownResult {
	s := new(GetMarkdownResult)
	s.Markdown = ""
	return s
}

// GetTranscriptArgs : Arguments for the asynchronous `get_transcript_async`
// route. Exactly one of `file_id`, `path`, or `url` must be supplied via
// `file_id_or_url` to identify the audio or video asset to transcribe.
type GetTranscriptArgs struct {
	// FileIdOrUrl : Identifier of the media asset to transcribe. Callers must
	// set exactly one of the oneof variants: - file_id: a Dropbox-issued file
	// id (format: "id:<id>") for a file the authenticated user has access to. -
	// path: an absolute Dropbox path, e.g. "/folder/recording.mp4". - url:
	// either a Dropbox shared link (www.dropbox.com) or an external HTTPS URL
	// pointing to a supported audio/video file. - Dropbox shared links are
	// resolved internally using the caller's authenticated identity and the
	// link's visibility / download settings. They therefore require an
	// authenticated user context (anonymous `url` requests against Dropbox
	// links are rejected with an `ACCESS_ERROR`). Links protected by a password
	// are rejected with `shared_link_password_protected`; links with downloads
	// disabled are rejected with `link_download_disabled_error`. - External
	// URLs are fetched over HTTPS through the backend's egress proxy and must
	// point at a supported audio/video file extension. The referenced asset
	// must be an audio or video file in a supported format; requests against
	// files with no audio track return a `no_audio_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
	// TimestampLevel : Granularity of the time offsets returned for each
	// transcript segment. Defaults to `SENTENCE. - SENTENCE: one segment per
	// spoken sentence (recommended). - WORD: one segment per word, useful for
	// fine-grained alignment such as captioning or highlight-as-you-listen
	// experiences.
	TimestampLevel *TimestampLevel `json:"timestamp_level"`
	// IncludedSpecialWords : Comma-delimited list of non-lexical filler words
	// to preserve in the transcript output, e.g. `"uh, ah, uhm"`. By default
	// these fillers are stripped. Unrecognized tokens are ignored. Leave empty
	// to use the default filtering behavior.
	IncludedSpecialWords string `json:"included_special_words"`
	// AudioLanguage : Optional BCP-47 language tag hinting the spoken language
	// of the source audio (e.g. "en-US", "ja-JP"). When empty, the service
	// auto-detects the language; supplying a hint improves accuracy and latency
	// for short or ambiguous clips. Unsupported languages fall back to
	// auto-detection.
	AudioLanguage string `json:"audio_language"`
}

// NewGetTranscriptArgs returns a new GetTranscriptArgs instance
func NewGetTranscriptArgs() *GetTranscriptArgs {
	s := new(GetTranscriptArgs)
	s.TimestampLevel = &TimestampLevel{Tagged: dropbox.Tagged{Tag: "unknown"}}
	s.IncludedSpecialWords = ""
	s.AudioLanguage = ""
	return s
}

// GetTranscriptAsyncCheckResult : Result type for EventBus async check - must
// end in "CheckResult"
type GetTranscriptAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : has no documentation (yet)
	Complete *GetTranscriptResult `json:"complete,omitempty"`
	// Failed : has no documentation (yet)
	Failed *GetTranscriptAsyncError `json:"failed,omitempty"`
}

// Valid tag values for GetTranscriptAsyncCheckResult
const (
	GetTranscriptAsyncCheckResultInProgress = "in_progress"
	GetTranscriptAsyncCheckResultComplete   = "complete"
	GetTranscriptAsyncCheckResultFailed     = "failed"
	GetTranscriptAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetTranscriptAsyncCheckResult instance
func (u *GetTranscriptAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "complete":
		if err = json.Unmarshal(body, &u.Complete); err != nil {
			return err
		}

	case "failed":
		if err = json.Unmarshal(body, &u.Failed); err != nil {
			return err
		}

	}
	return nil
}

// GetTranscriptAsyncError : has no documentation (yet)
type GetTranscriptAsyncError struct {
	// ErrorCode : has no documentation (yet)
	ErrorCode *ErrorCode `json:"error_code"`
	// ErrorDetails : has no documentation (yet)
	ErrorDetails *ContentApiV2Error `json:"error_details,omitempty"`
}

// NewGetTranscriptAsyncError returns a new GetTranscriptAsyncError instance
func NewGetTranscriptAsyncError() *GetTranscriptAsyncError {
	s := new(GetTranscriptAsyncError)
	s.ErrorCode = &ErrorCode{Tagged: dropbox.Tagged{Tag: "unknown_error"}}
	return s
}

// GetTranscriptResult : has no documentation (yet)
type GetTranscriptResult struct {
	// StructuredTranscript : The structured transcript produced for the
	// requested media asset, with per-segment text, start/end offsets (in
	// seconds from the beginning of the media), and the detected or
	// caller-supplied locale.
	StructuredTranscript *ApiStructuredTranscript `json:"structured_transcript,omitempty"`
}

// NewGetTranscriptResult returns a new GetTranscriptResult instance
func NewGetTranscriptResult() *GetTranscriptResult {
	s := new(GetTranscriptResult)
	return s
}

// MarkdownConversionApiV2Error : has no documentation (yet)
type MarkdownConversionApiV2Error struct {
	dropbox.Tagged
	// ServerError : has no documentation (yet)
	ServerError string `json:"server_error,omitempty"`
	// UserError : has no documentation (yet)
	UserError string `json:"user_error,omitempty"`
}

// Valid tag values for MarkdownConversionApiV2Error
const (
	MarkdownConversionApiV2ErrorServerError                 = "server_error"
	MarkdownConversionApiV2ErrorUserError                   = "user_error"
	MarkdownConversionApiV2ErrorUnsupportedFormatError      = "unsupported_format_error"
	MarkdownConversionApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	MarkdownConversionApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	MarkdownConversionApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	MarkdownConversionApiV2ErrorConversionFailureError      = "conversion_failure_error"
	MarkdownConversionApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a MarkdownConversionApiV2Error instance
func (u *MarkdownConversionApiV2Error) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// ServerError : has no documentation (yet)
		ServerError string `json:"server_error,omitempty"`
		// UserError : has no documentation (yet)
		UserError string `json:"user_error,omitempty"`
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "server_error":
		u.ServerError = w.ServerError

	case "user_error":
		u.UserError = w.UserError

	}
	return nil
}

// MediaDurationError : has no documentation (yet)
type MediaDurationError struct {
	// Limit : has no documentation (yet)
	Limit int32 `json:"limit"`
}

// NewMediaDurationError returns a new MediaDurationError instance
func NewMediaDurationError() *MediaDurationError {
	s := new(MediaDurationError)
	s.Limit = 0
	return s
}

// TimestampLevel : has no documentation (yet)
type TimestampLevel struct {
	dropbox.Tagged
}

// Valid tag values for TimestampLevel
const (
	TimestampLevelUnknown  = "unknown"
	TimestampLevelSentence = "sentence"
	TimestampLevelWord     = "word"
	TimestampLevelOther    = "other"
)
