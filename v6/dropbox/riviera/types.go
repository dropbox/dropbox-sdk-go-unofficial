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

// ApiExifGpsMetadata : GPS coordinates and related tags extracted from image
// EXIF data. Fields are populated on a best-effort basis and may be empty when
// absent from the source file.
type ApiExifGpsMetadata struct {
	// Latitude : Latitude / longitude in decimal degrees (positive = N/E,
	// negative = S/W).
	Latitude float32 `json:"latitude"`
	// Longitude : has no documentation (yet)
	Longitude float32 `json:"longitude"`
	// Altitude : Altitude in meters, as reported by the source (string to
	// preserve the original representation, which may include a reference
	// direction).
	Altitude string `json:"altitude"`
	// Timestamp : Timestamp / datestamp of the GPS fix, in the EXIF-provided
	// format.
	Timestamp string `json:"timestamp"`
	// Datestamp : has no documentation (yet)
	Datestamp string `json:"datestamp"`
}

// NewApiExifGpsMetadata returns a new ApiExifGpsMetadata instance
func NewApiExifGpsMetadata() *ApiExifGpsMetadata {
	s := new(ApiExifGpsMetadata)
	s.Latitude = 0.0
	s.Longitude = 0.0
	s.Altitude = ""
	s.Timestamp = ""
	s.Datestamp = ""
	return s
}

// ApiExifMetadata : Image EXIF metadata. Mirrors the useful subset of the
// internal `riviera.ExifMetadata` message. Fields are best-effort and may be
// empty.
type ApiExifMetadata struct {
	// ImageWidth : has no documentation (yet)
	ImageWidth uint32 `json:"image_width"`
	// ImageHeight : has no documentation (yet)
	ImageHeight uint32 `json:"image_height"`
	// CameraMake : has no documentation (yet)
	CameraMake string `json:"camera_make"`
	// CameraModel : has no documentation (yet)
	CameraModel string `json:"camera_model"`
	// LensModel : has no documentation (yet)
	LensModel string `json:"lens_model"`
	// DateTimeOriginal : Capture time in the EXIF-provided format (local time
	// of the camera).
	DateTimeOriginal string `json:"date_time_original"`
	// OffsetTimeOriginal : Timezone offset for `date_time_original`, e.g.
	// "+09:00".
	OffsetTimeOriginal string `json:"offset_time_original"`
	// Orientation : EXIF orientation value (1-8). See the EXIF spec; 1 is the
	// normal upright orientation.
	Orientation uint32 `json:"orientation"`
	// ExposureTime : fraction in string form, e.g. "1/250"
	ExposureTime string `json:"exposure_time"`
	// ApertureValue : has no documentation (yet)
	ApertureValue float64 `json:"aperture_value"`
	// IsoSpeed : has no documentation (yet)
	IsoSpeed uint32 `json:"iso_speed"`
	// FocalLength : e.g. "26.0 mm"
	FocalLength string `json:"focal_length"`
	// Megapixels : has no documentation (yet)
	Megapixels float64 `json:"megapixels"`
	// Artist : has no documentation (yet)
	Artist string `json:"artist"`
	// Copyright : has no documentation (yet)
	Copyright string `json:"copyright"`
	// GpsMetadata : has no documentation (yet)
	GpsMetadata *ApiExifGpsMetadata `json:"gps_metadata,omitempty"`
}

// NewApiExifMetadata returns a new ApiExifMetadata instance
func NewApiExifMetadata() *ApiExifMetadata {
	s := new(ApiExifMetadata)
	s.ImageWidth = 0
	s.ImageHeight = 0
	s.CameraMake = ""
	s.CameraModel = ""
	s.LensModel = ""
	s.DateTimeOriginal = ""
	s.OffsetTimeOriginal = ""
	s.Orientation = 0
	s.ExposureTime = ""
	s.ApertureValue = 0.0
	s.IsoSpeed = 0
	s.FocalLength = ""
	s.Megapixels = 0.0
	s.Artist = ""
	s.Copyright = ""
	return s
}

// ApiMediaMetadata : Audio/video container and per-stream metadata. Mirrors the
// useful subset of the internal `riviera.MediaMetadata` message.
type ApiMediaMetadata struct {
	// BitrateBps : has no documentation (yet)
	BitrateBps uint64 `json:"bitrate_bps"`
	// DurationS : has no documentation (yet)
	DurationS float64 `json:"duration_s"`
	// CreationTime : Container-level creation time, when present.
	CreationTime string `json:"creation_time"`
	// Streams : has no documentation (yet)
	Streams []*ApiMediaStream `json:"streams,omitempty"`
}

// NewApiMediaMetadata returns a new ApiMediaMetadata instance
func NewApiMediaMetadata() *ApiMediaMetadata {
	s := new(ApiMediaMetadata)
	s.BitrateBps = 0
	s.DurationS = 0.0
	s.CreationTime = ""
	return s
}

// ApiMediaStream : A single audio or video stream within a media file.
type ApiMediaStream struct {
	// Index : has no documentation (yet)
	Index uint32 `json:"index"`
	// CodecType : "audio", "video", etc.
	CodecType string `json:"codec_type"`
	// CodecName : has no documentation (yet)
	CodecName string `json:"codec_name"`
	// BitrateBps : has no documentation (yet)
	BitrateBps uint64 `json:"bitrate_bps"`
	// DurationS : has no documentation (yet)
	DurationS float64 `json:"duration_s"`
	// Width : Video-specific fields (zero / empty for audio streams).
	Width uint32 `json:"width"`
	// Height : has no documentation (yet)
	Height uint32 `json:"height"`
	// FramesPerSecond : has no documentation (yet)
	FramesPerSecond float64 `json:"frames_per_second"`
	// Rotation : has no documentation (yet)
	Rotation int32 `json:"rotation"`
	// DisplayAspectRatio : e.g. "16:9"
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	// Channels : Audio-specific fields (zero / empty for video streams).
	Channels uint32 `json:"channels"`
	// ChannelLayout : has no documentation (yet)
	ChannelLayout string `json:"channel_layout"`
	// SampleRateS : has no documentation (yet)
	SampleRateS uint64 `json:"sample_rate_s"`
	// LanguageIso639 : ISO 639 language code for the stream, when present.
	LanguageIso639 string `json:"language_iso_639"`
}

// NewApiMediaStream returns a new ApiMediaStream instance
func NewApiMediaStream() *ApiMediaStream {
	s := new(ApiMediaStream)
	s.Index = 0
	s.CodecType = ""
	s.CodecName = ""
	s.BitrateBps = 0
	s.DurationS = 0.0
	s.Width = 0
	s.Height = 0
	s.FramesPerSecond = 0.0
	s.Rotation = 0
	s.DisplayAspectRatio = ""
	s.Channels = 0
	s.ChannelLayout = ""
	s.SampleRateS = 0
	s.LanguageIso639 = ""
	return s
}

// ApiOfficeMetadata : MS Office document metadata. Mirrors the internal
// `riviera.OfficeMetadata` message. Some fields apply only to specific document
// types (e.g. `slides` for PowerPoint, `words`/`pages` for Word).
type ApiOfficeMetadata struct {
	// FileType : has no documentation (yet)
	FileType *OfficeFileType `json:"file_type"`
	// Creator : has no documentation (yet)
	Creator string `json:"creator"`
	// Company : has no documentation (yet)
	Company string `json:"company"`
	// Title : has no documentation (yet)
	Title string `json:"title"`
	// Subject : has no documentation (yet)
	Subject string `json:"subject"`
	// Keywords : has no documentation (yet)
	Keywords string `json:"keywords"`
	// Description : has no documentation (yet)
	Description string `json:"description"`
	// TotalEditTimeMinutes : has no documentation (yet)
	TotalEditTimeMinutes uint32 `json:"total_edit_time_minutes"`
	// Pages : Word only.
	Pages uint32 `json:"pages"`
	// Words : has no documentation (yet)
	Words uint32 `json:"words"`
	// Slides : PowerPoint only.
	Slides uint32 `json:"slides"`
	// RevisionNumber : has no documentation (yet)
	RevisionNumber string `json:"revision_number"`
}

// NewApiOfficeMetadata returns a new ApiOfficeMetadata instance
func NewApiOfficeMetadata() *ApiOfficeMetadata {
	s := new(ApiOfficeMetadata)
	s.FileType = &OfficeFileType{Tagged: dropbox.Tagged{Tag: "office_filetype_unknown"}}
	s.Creator = ""
	s.Company = ""
	s.Title = ""
	s.Subject = ""
	s.Keywords = ""
	s.Description = ""
	s.TotalEditTimeMinutes = 0
	s.Pages = 0
	s.Words = 0
	s.Slides = 0
	s.RevisionNumber = ""
	return s
}

// ApiPdfMetadata : PDF document metadata.
type ApiPdfMetadata struct {
	// Pages : has no documentation (yet)
	Pages uint32 `json:"pages"`
	// Width : Width / height of the first page, in PDF points.
	Width uint32 `json:"width"`
	// Height : has no documentation (yet)
	Height uint32 `json:"height"`
}

// NewApiPdfMetadata returns a new ApiPdfMetadata instance
func NewApiPdfMetadata() *ApiPdfMetadata {
	s := new(ApiPdfMetadata)
	s.Pages = 0
	s.Width = 0
	s.Height = 0
	return s
}

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

// ContentApiV2Error : Reason a transcript job failed. Returned in the `failed`
// variant of `GetTranscriptAsyncCheckResult`. This is a semantic error union:
// the HTTP status of the poll request itself is unaffected (a poll that
// surfaces a failed job is still a normal successful poll response). Callers
// should branch on the variant.
type ContentApiV2Error struct {
	dropbox.Tagged
	// ServerError : An unexpected, typically transient, server-side failure.
	// The string is a human-readable message; retrying with backoff may
	// succeed.
	ServerError string `json:"server_error,omitempty"`
	// UserError : The request could not be processed as supplied (a problem
	// with the caller's input). The string is a human-readable message;
	// retrying the same request will not help.
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
	ContentApiV2ErrorNotFoundError               = "not_found_error"
	ContentApiV2ErrorIsAFolderError              = "is_a_folder_error"
	ContentApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a ContentApiV2Error instance
func (u *ContentApiV2Error) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// ServerError : An unexpected, typically transient, server-side
		// failure. The string is a human-readable message; retrying with
		// backoff may succeed.
		ServerError string `json:"server_error,omitempty"`
		// UserError : The request could not be processed as supplied (a problem
		// with the caller's input). The string is a human-readable message;
		// retrying the same request will not help.
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

// FileIdOrUrl : has no documentation (yet)
type FileIdOrUrl struct {
	dropbox.Tagged
	// FileId : A Dropbox-issued file id (format: "id:<id>") for a file the
	// authenticated user has access to.
	FileId string `json:"file_id,omitempty"`
	// Url : Either a Dropbox shared link (www.dropbox.com) or an external HTTP
	// or HTTPS URL pointing to a supported file. - Dropbox shared links are
	// resolved internally using the caller's authenticated identity and the
	// link's visibility / download settings. They therefore require an
	// authenticated user context (anonymous `url` requests against Dropbox
	// links are rejected with an `access_error`). Links protected by a password
	// are rejected with `shared_link_password_protected`; links with downloads
	// disabled are rejected with `link_download_disabled_error`. - External
	// URLs are fetched through the backend's egress proxy and must point at a
	// supported file extension.
	Url string `json:"url,omitempty"`
	// Path : An absolute Dropbox path, e.g. "/folder/example.pdf".
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
		// FileId : A Dropbox-issued file id (format: "id:<id>") for a file the
		// authenticated user has access to.
		FileId string `json:"file_id,omitempty"`
		// Url : Either a Dropbox shared link (www.dropbox.com) or an external
		// HTTP or HTTPS URL pointing to a supported file. - Dropbox shared
		// links are resolved internally using the caller's authenticated
		// identity and the link's visibility / download settings. They
		// therefore require an authenticated user context (anonymous `url`
		// requests against Dropbox links are rejected with an `access_error`).
		// Links protected by a password are rejected with
		// `shared_link_password_protected`; links with downloads disabled are
		// rejected with `link_download_disabled_error`. - External URLs are
		// fetched through the backend's egress proxy and must point at a
		// supported file extension.
		Url string `json:"url,omitempty"`
		// Path : An absolute Dropbox path, e.g. "/folder/example.pdf".
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
	// exactly one of the `FileIdOrUrl` variants. The referenced file must be a
	// document in a supported format (see the route description for the list);
	// requests against unsupported formats return `unsupported_format_error`.
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
	Failed *MarkdownConversionApiV2Error `json:"failed,omitempty"`
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
		// Failed : has no documentation (yet)
		Failed *MarkdownConversionApiV2Error `json:"failed,omitempty"`
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
		u.Failed = w.Failed

	}
	return nil
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

// GetMetadataArgs : Arguments for the asynchronous `get_metadata_async` route.
// Exactly one of `file_id`, `path`, or `url` must be supplied via
// `file_id_or_url` to identify the file whose metadata should be extracted.
type GetMetadataArgs struct {
	// FileIdOrUrl : Identifier of the file to extract metadata from. Callers
	// must set exactly one of the `FileIdOrUrl` variants. The kind of metadata
	// returned is determined by the file type: image files return EXIF
	// metadata, audio/video files return media metadata, PDFs return PDF
	// metadata, and MS Office documents (docx, pptx, xlsx) return Office
	// metadata. See the route description for the supported formats. Requests
	// against unsupported formats return `unsupported_format_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
}

// NewGetMetadataArgs returns a new GetMetadataArgs instance
func NewGetMetadataArgs() *GetMetadataArgs {
	s := new(GetMetadataArgs)
	return s
}

// GetMetadataAsyncCheckResult : Result type for EventBus async check - must end
// in "CheckResult"
type GetMetadataAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : has no documentation (yet)
	Complete *GetMetadataResult `json:"complete,omitempty"`
	// Failed : has no documentation (yet)
	Failed *MetadataExtractionApiV2Error `json:"failed,omitempty"`
}

// Valid tag values for GetMetadataAsyncCheckResult
const (
	GetMetadataAsyncCheckResultInProgress = "in_progress"
	GetMetadataAsyncCheckResultComplete   = "complete"
	GetMetadataAsyncCheckResultFailed     = "failed"
	GetMetadataAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetMetadataAsyncCheckResult instance
func (u *GetMetadataAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// Failed : has no documentation (yet)
		Failed *MetadataExtractionApiV2Error `json:"failed,omitempty"`
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
		u.Failed = w.Failed

	}
	return nil
}

// GetMetadataResult : has no documentation (yet)
type GetMetadataResult struct {
	// MetadataType : The kind of metadata that was extracted for the requested
	// file. Callers should read the matching field of the `metadata` oneof.
	MetadataType *MetadataType `json:"metadata_type"`
	// Metadata : has no documentation (yet)
	Metadata *MetadataUnion `json:"metadata,omitempty"`
}

// NewGetMetadataResult returns a new GetMetadataResult instance
func NewGetMetadataResult() *GetMetadataResult {
	s := new(GetMetadataResult)
	s.MetadataType = &MetadataType{Tagged: dropbox.Tagged{Tag: "metadata_type_unknown"}}
	return s
}

// GetTextArgs : Arguments for the asynchronous `get_text_async` route. Exactly
// one of `file_id`, `path`, or `url` must be supplied via `file_id_or_url` to
// identify the document whose plain-text content should be extracted.
type GetTextArgs struct {
	// FileIdOrUrl : Identifier of the document to extract text from. Callers
	// must set exactly one of the `FileIdOrUrl` variants. Text extraction is
	// supported for common document formats (Word, PowerPoint, Excel, PDF, RTF,
	// and Dropbox document types); see the route description for the supported
	// formats. Requests against unsupported formats return
	// `unsupported_format_error`. NOTE: for the `url` variant, only Dropbox
	// shared links (www.dropbox.com) are supported. External (non-Dropbox) URLs
	// are not supported and return `unsupported_format_error`; import the file
	// into Dropbox and reference it by `file_id` or `path` instead.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
}

// NewGetTextArgs returns a new GetTextArgs instance
func NewGetTextArgs() *GetTextArgs {
	s := new(GetTextArgs)
	return s
}

// GetTextAsyncCheckResult : Result type for EventBus async check - must end in
// "CheckResult"
type GetTextAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : has no documentation (yet)
	Complete *GetTextResult `json:"complete,omitempty"`
	// Failed : has no documentation (yet)
	Failed *TextExtractionApiV2Error `json:"failed,omitempty"`
}

// Valid tag values for GetTextAsyncCheckResult
const (
	GetTextAsyncCheckResultInProgress = "in_progress"
	GetTextAsyncCheckResultComplete   = "complete"
	GetTextAsyncCheckResultFailed     = "failed"
	GetTextAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetTextAsyncCheckResult instance
func (u *GetTextAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// Failed : has no documentation (yet)
		Failed *TextExtractionApiV2Error `json:"failed,omitempty"`
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
		u.Failed = w.Failed

	}
	return nil
}

// GetTextResult : has no documentation (yet)
type GetTextResult struct {
	// Text : The plain-text content extracted from the document. For multi-page
	// documents the text is concatenated in document order. May be empty when
	// no text is detected in the source.
	Text string `json:"text"`
}

// NewGetTextResult returns a new GetTextResult instance
func NewGetTextResult() *GetTextResult {
	s := new(GetTextResult)
	s.Text = ""
	return s
}

// GetTranscriptArgs : Arguments for the asynchronous `get_transcript_async`
// route. Exactly one of `file_id`, `path`, or `url` must be supplied via
// `file_id_or_url` to identify the audio or video asset to transcribe.
type GetTranscriptArgs struct {
	// FileIdOrUrl : Identifier of the media asset to transcribe. Callers must
	// set exactly one of the `FileIdOrUrl` variants. The referenced asset must
	// be an audio or video file in a supported format (see the route
	// description for the list); requests against files with no audio track
	// return a `no_audio_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
	// TimestampLevel : Granularity of the time offsets returned for each
	// transcript segment. Defaults to `SENTENCE` when the field is omitted. -
	// SENTENCE: one segment per spoken sentence (recommended). - WORD: one
	// segment per word, useful for fine-grained alignment such as captioning or
	// highlight-as-you-listen experiences.
	TimestampLevel *TimestampLevel `json:"timestamp_level"`
	// IncludedSpecialWords : Comma-delimited list of non-lexical filler words
	// to preserve in the transcript output, e.g. `"uh, ah, uhm"`. By default
	// these fillers are stripped. Unrecognized tokens are ignored. Leave empty
	// to use the default filtering behavior.
	IncludedSpecialWords string `json:"included_special_words"`
	// AudioLanguage : Optional ISO 639-1 two-letter language code hinting the
	// spoken language of the source audio (e.g. "en", "ja"). When empty, the
	// service auto-detects the language; supplying a hint improves accuracy and
	// latency for short or ambiguous clips. Unsupported languages fall back to
	// auto-detection.
	AudioLanguage string `json:"audio_language"`
}

// NewGetTranscriptArgs returns a new GetTranscriptArgs instance
func NewGetTranscriptArgs() *GetTranscriptArgs {
	s := new(GetTranscriptArgs)
	s.TimestampLevel = &TimestampLevel{Tagged: dropbox.Tagged{Tag: "sentence"}}
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
	Failed *ContentApiV2Error `json:"failed,omitempty"`
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
		// Failed : has no documentation (yet)
		Failed *ContentApiV2Error `json:"failed,omitempty"`
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
		u.Failed = w.Failed

	}
	return nil
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

// MarkdownConversionApiV2Error : Reason a markdown conversion job failed.
// Returned in the `failed` variant of `GetMarkdownAsyncCheckResult`. This is a
// semantic error union: the HTTP status of the poll request itself is
// unaffected (a poll that surfaces a failed job is still a normal successful
// poll response). Callers should branch on the variant.
type MarkdownConversionApiV2Error struct {
	dropbox.Tagged
	// ServerError : An unexpected, typically transient, server-side failure.
	// The string is a human-readable message; retrying with backoff may
	// succeed.
	ServerError string `json:"server_error,omitempty"`
	// UserError : The request could not be processed as supplied (a problem
	// with the caller's input). The string is a human-readable message;
	// retrying the same request will not help.
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
	MarkdownConversionApiV2ErrorNotFoundError               = "not_found_error"
	MarkdownConversionApiV2ErrorIsAFolderError              = "is_a_folder_error"
	MarkdownConversionApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a MarkdownConversionApiV2Error instance
func (u *MarkdownConversionApiV2Error) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// ServerError : An unexpected, typically transient, server-side
		// failure. The string is a human-readable message; retrying with
		// backoff may succeed.
		ServerError string `json:"server_error,omitempty"`
		// UserError : The request could not be processed as supplied (a problem
		// with the caller's input). The string is a human-readable message;
		// retrying the same request will not help.
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

// MetadataExtractionApiV2Error : Reason a metadata extraction job failed.
// Returned in the `failed` variant of `GetMetadataAsyncCheckResult`. This is a
// semantic error union: the HTTP status of the poll request itself is
// unaffected (a poll that surfaces a failed job is still a normal successful
// poll response). Callers should branch on the variant.
type MetadataExtractionApiV2Error struct {
	dropbox.Tagged
	// ServerError : An unexpected, typically transient, server-side failure.
	// The string is a human-readable message; retrying with backoff may
	// succeed.
	ServerError string `json:"server_error,omitempty"`
	// UserError : The request could not be processed as supplied (a problem
	// with the caller's input). The string is a human-readable message;
	// retrying the same request will not help.
	UserError string `json:"user_error,omitempty"`
}

// Valid tag values for MetadataExtractionApiV2Error
const (
	MetadataExtractionApiV2ErrorServerError                 = "server_error"
	MetadataExtractionApiV2ErrorUserError                   = "user_error"
	MetadataExtractionApiV2ErrorUnsupportedFormatError      = "unsupported_format_error"
	MetadataExtractionApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	MetadataExtractionApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	MetadataExtractionApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	MetadataExtractionApiV2ErrorConversionFailureError      = "conversion_failure_error"
	MetadataExtractionApiV2ErrorNotFoundError               = "not_found_error"
	MetadataExtractionApiV2ErrorIsAFolderError              = "is_a_folder_error"
	MetadataExtractionApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a MetadataExtractionApiV2Error instance
func (u *MetadataExtractionApiV2Error) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// ServerError : An unexpected, typically transient, server-side
		// failure. The string is a human-readable message; retrying with
		// backoff may succeed.
		ServerError string `json:"server_error,omitempty"`
		// UserError : The request could not be processed as supplied (a problem
		// with the caller's input). The string is a human-readable message;
		// retrying the same request will not help.
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

// MetadataType : Which metadata variant is populated in a `GetMetadataResult`,
// derived from the file type.
type MetadataType struct {
	dropbox.Tagged
}

// Valid tag values for MetadataType
const (
	MetadataTypeMetadataTypeUnknown = "metadata_type_unknown"
	MetadataTypeMetadataTypeExif    = "metadata_type_exif"
	MetadataTypeMetadataTypeMedia   = "metadata_type_media"
	MetadataTypeMetadataTypePdf     = "metadata_type_pdf"
	MetadataTypeMetadataTypeOffice  = "metadata_type_office"
	MetadataTypeOther               = "other"
)

// OfficeFileType : The kind of MS Office document that produced an
// `ApiOfficeMetadata` result.
type OfficeFileType struct {
	dropbox.Tagged
}

// Valid tag values for OfficeFileType
const (
	OfficeFileTypeOfficeFiletypeUnknown    = "office_filetype_unknown"
	OfficeFileTypeOfficeFiletypeWord       = "office_filetype_word"
	OfficeFileTypeOfficeFiletypePowerpoint = "office_filetype_powerpoint"
	OfficeFileTypeOfficeFiletypeExcel      = "office_filetype_excel"
	OfficeFileTypeOther                    = "other"
)

// TextExtractionApiV2Error : Reason a text extraction job failed. Returned in
// the `failed` variant of `GetTextAsyncCheckResult`. This is a semantic error
// union: the HTTP status of the poll request itself is unaffected (a poll that
// surfaces a failed job is still a normal successful poll response). Callers
// should branch on the variant.
type TextExtractionApiV2Error struct {
	dropbox.Tagged
	// ServerError : An unexpected, typically transient, server-side failure.
	// The string is a human-readable message; retrying with backoff may
	// succeed.
	ServerError string `json:"server_error,omitempty"`
	// UserError : The request could not be processed as supplied (a problem
	// with the caller's input). The string is a human-readable message;
	// retrying the same request will not help.
	UserError string `json:"user_error,omitempty"`
}

// Valid tag values for TextExtractionApiV2Error
const (
	TextExtractionApiV2ErrorServerError                 = "server_error"
	TextExtractionApiV2ErrorUserError                   = "user_error"
	TextExtractionApiV2ErrorUnsupportedFormatError      = "unsupported_format_error"
	TextExtractionApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	TextExtractionApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	TextExtractionApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	TextExtractionApiV2ErrorConversionFailureError      = "conversion_failure_error"
	TextExtractionApiV2ErrorNotFoundError               = "not_found_error"
	TextExtractionApiV2ErrorIsAFolderError              = "is_a_folder_error"
	TextExtractionApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a TextExtractionApiV2Error instance
func (u *TextExtractionApiV2Error) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// ServerError : An unexpected, typically transient, server-side
		// failure. The string is a human-readable message; retrying with
		// backoff may succeed.
		ServerError string `json:"server_error,omitempty"`
		// UserError : The request could not be processed as supplied (a problem
		// with the caller's input). The string is a human-readable message;
		// retrying the same request will not help.
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

// TimestampLevel : has no documentation (yet)
type TimestampLevel struct {
	dropbox.Tagged
}

// Valid tag values for TimestampLevel
const (
	TimestampLevelSentence = "sentence"
	TimestampLevelWord     = "word"
	TimestampLevelOther    = "other"
)

// MetadataUnion : Exactly one variant is populated, corresponding to
// `metadata_type`.
type MetadataUnion struct {
	dropbox.Tagged
	// Exif : has no documentation (yet)
	Exif *ApiExifMetadata `json:"exif,omitempty"`
	// Media : has no documentation (yet)
	Media *ApiMediaMetadata `json:"media,omitempty"`
	// Pdf : has no documentation (yet)
	Pdf *ApiPdfMetadata `json:"pdf,omitempty"`
	// Office : has no documentation (yet)
	Office *ApiOfficeMetadata `json:"office,omitempty"`
}

// Valid tag values for MetadataUnion
const (
	MetadataUnionExif   = "exif"
	MetadataUnionMedia  = "media"
	MetadataUnionPdf    = "pdf"
	MetadataUnionOffice = "office"
	MetadataUnionOther  = "other"
)

// UnmarshalJSON deserializes into a MetadataUnion instance
func (u *MetadataUnion) UnmarshalJSON(body []byte) error {
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
	case "exif":
		if err = json.Unmarshal(body, &u.Exif); err != nil {
			return err
		}

	case "media":
		if err = json.Unmarshal(body, &u.Media); err != nil {
			return err
		}

	case "pdf":
		if err = json.Unmarshal(body, &u.Pdf); err != nil {
			return err
		}

	case "office":
		if err = json.Unmarshal(body, &u.Office); err != nil {
			return err
		}

	}
	return nil
}
