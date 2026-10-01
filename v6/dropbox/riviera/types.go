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
	// Latitude : Latitude in decimal degrees (positive = north, negative =
	// south).
	Latitude float32 `json:"latitude"`
	// Longitude : Longitude in decimal degrees (positive = east, negative =
	// west).
	Longitude float32 `json:"longitude"`
	// Altitude : Altitude in meters, as reported by the source (string to
	// preserve the original representation, which may include a reference
	// direction).
	Altitude string `json:"altitude"`
	// Timestamp : Time of the GPS fix, in the EXIF-provided format.
	Timestamp string `json:"timestamp"`
	// Datestamp : Date of the GPS fix, in the EXIF-provided format.
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

// ApiExifMetadata : Image EXIF metadata. Fields are populated on a best-effort
// basis and may be empty when absent from the source file.
type ApiExifMetadata struct {
	// ImageWidth : Width of the image, in pixels.
	ImageWidth uint32 `json:"image_width"`
	// ImageHeight : Height of the image, in pixels.
	ImageHeight uint32 `json:"image_height"`
	// CameraMake : Manufacturer of the device that captured the image, e.g.
	// "Apple".
	CameraMake string `json:"camera_make"`
	// CameraModel : Model of the device that captured the image, e.g. "iPhone
	// 15 Pro".
	CameraModel string `json:"camera_model"`
	// LensModel : Model of the lens the image was captured with, when the
	// source records it.
	LensModel string `json:"lens_model"`
	// DateTimeOriginal : Capture time in the EXIF-provided format (local time
	// of the camera).
	DateTimeOriginal string `json:"date_time_original"`
	// OffsetTimeOriginal : Timezone offset for
	// `ApiExifMetadata.date_time_original`, e.g. "+09:00".
	OffsetTimeOriginal string `json:"offset_time_original"`
	// Orientation : EXIF orientation value (1-8). See the EXIF spec; 1 is the
	// normal upright orientation.
	Orientation uint32 `json:"orientation"`
	// ExposureTime : Exposure time the image was captured with, as a
	// fractional-second string, e.g. "1/250".
	ExposureTime string `json:"exposure_time"`
	// ApertureValue : Aperture the image was captured at, as reported by the
	// EXIF aperture tag.
	ApertureValue float64 `json:"aperture_value"`
	// IsoSpeed : ISO sensitivity the image was captured at.
	IsoSpeed uint32 `json:"iso_speed"`
	// FocalLength : Focal length the image was captured at, including the unit,
	// e.g. "26.0 mm".
	FocalLength string `json:"focal_length"`
	// Megapixels : Total pixel count of the image, in megapixels.
	Megapixels float64 `json:"megapixels"`
	// Artist : Creator credited in the EXIF artist tag.
	Artist string `json:"artist"`
	// Copyright : Copyright notice from the EXIF copyright tag.
	Copyright string `json:"copyright"`
	// GpsMetadata : Location tags from the image, when the source recorded a
	// location.
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

// ApiKeyframe : A single extracted scene-change keyframe.
type ApiKeyframe struct {
	// Timestamp : Presentation timestamp of the keyframe, in seconds from the
	// start of the video.
	Timestamp float64 `json:"timestamp"`
	// SceneScore : Scene-change score that triggered this keyframe, in the
	// range [0.0, 1.0]. Higher values indicate a more pronounced scene change
	// relative to the preceding frame. The first keyframe of a video is always
	// reported as 1.0: the start of a video is a scene boundary by definition,
	// so that score is not a measured frame-to-frame comparison.
	SceneScore float64 `json:"scene_score"`
	// ImageBase64 : The extracted frame as a base64-encoded JPEG image. Empty
	// when the request set `include_images = false`.
	ImageBase64 string `json:"image_base64"`
}

// NewApiKeyframe returns a new ApiKeyframe instance
func NewApiKeyframe() *ApiKeyframe {
	s := new(ApiKeyframe)
	s.Timestamp = 0.0
	s.SceneScore = 0.0
	s.ImageBase64 = ""
	return s
}

// ApiMediaMetadata : Audio/video container and per-stream metadata. Fields are
// populated on a best-effort basis and may be empty when absent from the source
// file.
type ApiMediaMetadata struct {
	// BitrateBps : Overall bitrate of the container, in bits per second.
	BitrateBps uint64 `json:"bitrate_bps"`
	// DurationS : Duration of the media, in seconds.
	DurationS float64 `json:"duration_s"`
	// CreationTime : Container-level creation time, when present.
	CreationTime string `json:"creation_time"`
	// Streams : The audio and video streams the container holds, in container
	// order.
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
	// Index : Zero-based index of the stream within the container.
	Index uint32 `json:"index"`
	// CodecType : Kind of media the stream carries, e.g. "audio" or "video".
	CodecType string `json:"codec_type"`
	// CodecName : Name of the codec the stream is encoded with, e.g. "h264" or
	// "aac".
	CodecName string `json:"codec_name"`
	// BitrateBps : Bitrate of this stream, in bits per second.
	BitrateBps uint64 `json:"bitrate_bps"`
	// DurationS : Duration of this stream, in seconds.
	DurationS float64 `json:"duration_s"`
	// Width : Width of the video frame, in pixels. Zero for audio streams.
	Width uint32 `json:"width"`
	// Height : Height of the video frame, in pixels. Zero for audio streams.
	Height uint32 `json:"height"`
	// FramesPerSecond : Frame rate of the stream, in frames per second. Zero
	// for audio streams.
	FramesPerSecond float64 `json:"frames_per_second"`
	// Rotation : Rotation to apply on playback, in degrees, as recorded in the
	// stream metadata. Zero for audio streams and for video that needs no
	// rotation.
	Rotation int32 `json:"rotation"`
	// DisplayAspectRatio : Aspect ratio the video should be displayed at, as a
	// "width:height" string, e.g. "16:9". Empty for audio streams.
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	// Channels : Number of audio channels in the stream. Zero for video
	// streams.
	Channels uint32 `json:"channels"`
	// ChannelLayout : Layout of the audio channels, e.g. "stereo". Empty for
	// video streams.
	ChannelLayout string `json:"channel_layout"`
	// SampleRateS : Sample rate of the audio stream, in samples per second.
	// Zero for video streams.
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

// ApiOfficeMetadata : MS Office document metadata. Some fields apply only to
// specific document types (e.g. `ApiOfficeMetadata.slides` for PowerPoint,
// `ApiOfficeMetadata.words` and `ApiOfficeMetadata.pages` for Word).
type ApiOfficeMetadata struct {
	// FileType : Which kind of Office document this metadata was extracted
	// from.
	FileType *OfficeFileType `json:"file_type"`
	// Creator : Author recorded in the document properties.
	Creator string `json:"creator"`
	// Company : Company recorded in the document properties.
	Company string `json:"company"`
	// Title : Title recorded in the document properties.
	Title string `json:"title"`
	// Subject : Subject recorded in the document properties.
	Subject string `json:"subject"`
	// Keywords : Keywords recorded in the document properties, in the
	// document's own formatting (typically a single comma- or space-separated
	// string).
	Keywords string `json:"keywords"`
	// Description : Description recorded in the document properties.
	Description string `json:"description"`
	// TotalEditTimeMinutes : Total editing time recorded in the document
	// properties, in minutes.
	TotalEditTimeMinutes uint32 `json:"total_edit_time_minutes"`
	// Pages : Page count recorded in the document properties. Word documents
	// only; zero for PowerPoint and Excel.
	Pages uint32 `json:"pages"`
	// Words : Word count recorded in the document properties. Word documents
	// only; zero for PowerPoint and Excel.
	Words uint32 `json:"words"`
	// Slides : Slide count recorded in the document properties. PowerPoint
	// documents only; zero for Word and Excel.
	Slides uint32 `json:"slides"`
	// RevisionNumber : Revision number recorded in the document properties.
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
	// Pages : Number of pages in the document.
	Pages uint32 `json:"pages"`
	// Width : Width of the first page, in PDF points.
	Width uint32 `json:"width"`
	// Height : Height of the first page, in PDF points.
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

// ApiStructuredTranscript : A transcript, split into segments.
type ApiStructuredTranscript struct {
	// Segments : The segments of the transcript, in playback order.
	Segments []*ApiTranscriptSegment `json:"segments,omitempty"`
	// TranscriptLocale : The language of the transcript, as an ISO 639-1 code
	// (e.g. "en"). This is the language detected in the audio, or the one
	// supplied in `GetTranscriptArgs.audio_language`.
	TranscriptLocale string `json:"transcript_locale"`
}

// NewApiStructuredTranscript returns a new ApiStructuredTranscript instance
func NewApiStructuredTranscript() *ApiStructuredTranscript {
	s := new(ApiStructuredTranscript)
	s.TranscriptLocale = ""
	return s
}

// ApiTranscriptSegment : A contiguous span of transcribed speech. The span
// covered by a segment depends on the requested `TimestampLevel`.
type ApiTranscriptSegment struct {
	// Text : The transcribed text of this segment.
	Text string `json:"text"`
	// StartTime : Offset of the start of this segment, in seconds from the
	// beginning of the media.
	StartTime float64 `json:"start_time"`
	// EndTime : Offset of the end of this segment, in seconds from the
	// beginning of the media.
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

// ContentApiV2Error : Reason a transcript job failed. Returned in the
// `GetTranscriptAsyncCheckResult.failed` variant. This is a semantic error
// union: the HTTP status of the poll request itself is unaffected (a poll that
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
	// MediaDurationError : The audio to transcribe is longer than the supported
	// maximum.
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

// DownloadTransformOutputArgs : Arguments for `download_transform_output`.
type DownloadTransformOutputArgs struct {
	// OutputHandle : The `output_handle` from a `complete`
	// `get_transform_async/check` result.
	OutputHandle string `json:"output_handle"`
}

// NewDownloadTransformOutputArgs returns a new DownloadTransformOutputArgs instance
func NewDownloadTransformOutputArgs(OutputHandle string) *DownloadTransformOutputArgs {
	s := new(DownloadTransformOutputArgs)
	s.OutputHandle = OutputHandle
	return s
}

// DownloadTransformOutputResult : Describes the bytes in the response body.
// Returned in the `Dropbox-API-Result` header.
type DownloadTransformOutputResult struct {
	// Size : Size of the output in bytes.
	Size uint64 `json:"size"`
	// Format : Format of the output, as a short lowercase format name such as
	// "pdf", "html", "jpeg", or "png".
	Format string `json:"format"`
	// MimeType : MIME type corresponding to `format`.
	MimeType string `json:"mime_type"`
}

// NewDownloadTransformOutputResult returns a new DownloadTransformOutputResult instance
func NewDownloadTransformOutputResult() *DownloadTransformOutputResult {
	s := new(DownloadTransformOutputResult)
	s.Size = 0
	s.Format = ""
	s.MimeType = ""
	return s
}

// FileIdOrUrl : has no documentation (yet)
type FileIdOrUrl struct {
	dropbox.Tagged
	// FileId : A Dropbox-issued file ID for a file the authenticated user has
	// access to, e.g. "id:a4ayc_80_OEAAAAAAAAAYa".
	FileId string `json:"file_id,omitempty"`
	// Url : Either a Dropbox shared link (www.dropbox.com) or an
	// internet-accessible URL pointing to a supported file. - Dropbox shared
	// links are resolved internally using the caller's authenticated identity
	// and the link's visibility / download settings. They therefore require an
	// authenticated user context; requests made with app auth alone are
	// rejected. Password-protected links and links with downloads disabled are
	// rejected as well. - Other URLs are fetched by Dropbox's servers, so they
	// must be reachable from the public internet -- not only from the calling
	// application's network -- and must point at a supported file extension.
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
		// FileId : A Dropbox-issued file ID for a file the authenticated user
		// has access to, e.g. "id:a4ayc_80_OEAAAAAAAAAYa".
		FileId string `json:"file_id,omitempty"`
		// Url : Either a Dropbox shared link (www.dropbox.com) or an
		// internet-accessible URL pointing to a supported file. - Dropbox
		// shared links are resolved internally using the caller's authenticated
		// identity and the link's visibility / download settings. They
		// therefore require an authenticated user context; requests made with
		// app auth alone are rejected. Password-protected links and links with
		// downloads disabled are rejected as well. - Other URLs are fetched by
		// Dropbox's servers, so they must be reachable from the public internet
		// -- not only from the calling application's network -- and must point
		// at a supported file extension.
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

// GetKeyframesArgs : Arguments for the asynchronous `get_keyframes_async`
// route. Exactly one of `file_id`, `path`, or `url` must be supplied via
// `file_id_or_url` to identify the video whose scene-change keyframes should be
// extracted.
type GetKeyframesArgs struct {
	// FileIdOrUrl : Identifier of the video file to extract keyframes from.
	// Callers must set exactly one of the `FileIdOrUrl` variants. Keyframe
	// extraction is supported for video files only; see the route description
	// for the supported formats. Requests against unsupported formats return
	// `unsupported_format_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
	// SceneChangeThreshold : Sensitivity of scene-change detection. A keyframe
	// is emitted whenever the frame-to-frame scene score crosses this
	// threshold, so a LOWER value yields MORE keyframes. Valid range is (0.0,
	// 1.0]. When omitted (0.0) the service uses a default of 0.3, which is a
	// good starting point for most videos.
	SceneChangeThreshold float64 `json:"scene_change_threshold"`
	// IncludeImages : When true, each returned keyframe includes the JPEG image
	// bytes, base64-encoded, in `ApiKeyframe.image_base64`. When false, the
	// response contains only per-keyframe metadata (timestamp and scene score)
	// and `image_base64` is left empty -- useful when you only need the scene
	// boundaries and want a small response. NOTE: because the field defaults to
	// false in proto3, callers who want images must set this explicitly to
	// true.
	IncludeImages bool `json:"include_images"`
}

// NewGetKeyframesArgs returns a new GetKeyframesArgs instance
func NewGetKeyframesArgs() *GetKeyframesArgs {
	s := new(GetKeyframesArgs)
	s.SceneChangeThreshold = 0.0
	s.IncludeImages = false
	return s
}

// GetKeyframesAsyncCheckResult : Result type for EventBus async check - must
// end in "CheckResult"
type GetKeyframesAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : has no documentation (yet)
	Complete *GetKeyframesResult `json:"complete,omitempty"`
	// Failed : has no documentation (yet)
	Failed *KeyframesExtractionApiV2Error `json:"failed,omitempty"`
}

// Valid tag values for GetKeyframesAsyncCheckResult
const (
	GetKeyframesAsyncCheckResultInProgress = "in_progress"
	GetKeyframesAsyncCheckResultComplete   = "complete"
	GetKeyframesAsyncCheckResultFailed     = "failed"
	GetKeyframesAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetKeyframesAsyncCheckResult instance
func (u *GetKeyframesAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// Failed : has no documentation (yet)
		Failed *KeyframesExtractionApiV2Error `json:"failed,omitempty"`
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

// GetKeyframesResult : has no documentation (yet)
type GetKeyframesResult struct {
	// Frames : The extracted keyframes, ordered by `timestamp`. May be empty
	// when no scene changes are detected in the source.
	Frames []*ApiKeyframe `json:"frames,omitempty"`
}

// NewGetKeyframesResult returns a new GetKeyframesResult instance
func NewGetKeyframesResult() *GetKeyframesResult {
	s := new(GetKeyframesResult)
	return s
}

// GetMarkdownArgs : Arguments for the asynchronous `getMarkdownAsync` route.
// Exactly one of `FileIdOrUrl.file_id`, `FileIdOrUrl.path`, or
// `FileIdOrUrl.url` must be supplied via `GetMarkdownArgs.file_id_or_url` to
// identify the document to convert to markdown.
type GetMarkdownArgs struct {
	// FileIdOrUrl : Identifier of the document to convert. Callers must set
	// exactly one of the `FileIdOrUrl` variants. The referenced file must be a
	// document in a supported format (see the route description for the list);
	// requests against unsupported formats fail with
	// `MarkdownConversionApiV2Error.user_error`.
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

// GetMarkdownAsyncCheckResult : Status of a markdown conversion job started by
// `getMarkdownAsync`, as returned by `getMarkdownAsyncCheck`.
type GetMarkdownAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : The job finished successfully.
	Complete *GetMarkdownResult `json:"complete,omitempty"`
	// Failed : The job finished unsuccessfully.
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
		// Failed : The job finished unsuccessfully.
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
	// Markdown : The markdown the source document was converted to.
	Markdown string `json:"markdown"`
}

// NewGetMarkdownResult returns a new GetMarkdownResult instance
func NewGetMarkdownResult() *GetMarkdownResult {
	s := new(GetMarkdownResult)
	s.Markdown = ""
	return s
}

// GetMetadataArgs : Arguments for the asynchronous `getMetadataAsync` route.
// Exactly one of `FileIdOrUrl.file_id`, `FileIdOrUrl.path`, or
// `FileIdOrUrl.url` must be supplied via `GetMetadataArgs.file_id_or_url` to
// identify the file whose metadata should be extracted.
type GetMetadataArgs struct {
	// FileIdOrUrl : Identifier of the file to extract metadata from. Callers
	// must set exactly one of the `FileIdOrUrl` variants. The kind of metadata
	// returned is determined by the file type: image files return EXIF
	// metadata, audio/video files return media metadata, PDFs return PDF
	// metadata, and MS Office documents (docx, pptx, xlsx) return Office
	// metadata. See the route description for the supported formats. Requests
	// against unsupported formats fail with
	// `MetadataExtractionApiV2Error.user_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
}

// NewGetMetadataArgs returns a new GetMetadataArgs instance
func NewGetMetadataArgs() *GetMetadataArgs {
	s := new(GetMetadataArgs)
	return s
}

// GetMetadataAsyncCheckResult : Status of a metadata extraction job started by
// `getMetadataAsync`, as returned by `getMetadataAsyncCheck`.
type GetMetadataAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : The job finished successfully.
	Complete *GetMetadataResult `json:"complete,omitempty"`
	// Failed : The job finished unsuccessfully.
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
		// Failed : The job finished unsuccessfully.
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
	// file. Callers should read the matching variant of
	// `GetMetadataResult.metadata`.
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

// GetOcrArgs : Arguments for the asynchronous `get_ocr_async` route. Exactly
// one of `file_id`, `path`, or `url` must be supplied via `file_id_or_url` to
// identify the image or PDF whose text should be extracted via OCR (optical
// character recognition).
type GetOcrArgs struct {
	// FileIdOrUrl : Identifier of the file to run OCR on. Callers must set
	// exactly one of the `FileIdOrUrl` variants. OCR is supported for image
	// files and PDFs, including scanned / non-text PDFs; see the route
	// description for the supported formats. Requests against unsupported
	// formats return `unsupported_format_error`. NOTE: for the `url` variant,
	// only Dropbox shared links (www.dropbox.com) are supported. External
	// (non-Dropbox) URLs are not supported and return
	// `unsupported_format_error`; import the file into Dropbox and reference it
	// by `file_id` or `path` instead.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
}

// NewGetOcrArgs returns a new GetOcrArgs instance
func NewGetOcrArgs() *GetOcrArgs {
	s := new(GetOcrArgs)
	return s
}

// GetOcrAsyncCheckResult : Result type for EventBus async check - must end in
// "CheckResult"
type GetOcrAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : has no documentation (yet)
	Complete *GetOcrResult `json:"complete,omitempty"`
	// Failed : has no documentation (yet)
	Failed *OcrExtractionApiV2Error `json:"failed,omitempty"`
}

// Valid tag values for GetOcrAsyncCheckResult
const (
	GetOcrAsyncCheckResultInProgress = "in_progress"
	GetOcrAsyncCheckResultComplete   = "complete"
	GetOcrAsyncCheckResultFailed     = "failed"
	GetOcrAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetOcrAsyncCheckResult instance
func (u *GetOcrAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// Failed : has no documentation (yet)
		Failed *OcrExtractionApiV2Error `json:"failed,omitempty"`
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

// GetOcrResult : has no documentation (yet)
type GetOcrResult struct {
	// Text : The plain-text content extracted from the file via OCR. Words
	// within a line are separated by a single space, lines are
	// newline-separated in reading order, and for multi-page PDFs pages are
	// separated by a blank line in page order. May be empty when no text is
	// detected in the source.
	Text string `json:"text"`
	// Hocr : The same content as hOCR: HTML that carries the position of every
	// recognized word. Each page is a `<section>` holding `<p class="line">`
	// elements with one `<span>` per word, and each element carries `data-x`,
	// `data-y`, `data-width`, and `data-height` attributes in pixels relative
	// to the upright page (whose dimensions are on the `<section>`). Use this
	// when you need word coordinates -- to highlight matches over a page image,
	// for example; use `text` when you just need the words.
	Hocr string `json:"hocr"`
}

// NewGetOcrResult returns a new GetOcrResult instance
func NewGetOcrResult() *GetOcrResult {
	s := new(GetOcrResult)
	s.Text = ""
	s.Hocr = ""
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

// GetTranscriptArgs : Arguments for the asynchronous `getTranscriptAsync`
// route. Exactly one of `FileIdOrUrl.file_id`, `FileIdOrUrl.path`, or
// `FileIdOrUrl.url` must be supplied via `GetTranscriptArgs.file_id_or_url` to
// identify the audio or video asset to transcribe.
type GetTranscriptArgs struct {
	// FileIdOrUrl : Identifier of the media asset to transcribe. Callers must
	// set exactly one of the `FileIdOrUrl` variants. The referenced asset must
	// be an audio or video file in a supported format (see the route
	// description for the list); requests against files with no audio track
	// fail with `ContentApiV2Error.no_audio_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
	// TimestampLevel : Granularity of the time offsets returned for each
	// transcript segment. Defaults to `TimestampLevel.sentence` when the field
	// is omitted.
	TimestampLevel *TimestampLevel `json:"timestamp_level"`
	// IncludedSpecialWords : Comma-delimited list of non-lexical filler words
	// to preserve in the transcript output, e.g. `"uh, ah, uhm"`. By default
	// these fillers are stripped. Unrecognized tokens are ignored. Leave empty
	// to use the default filtering behavior.
	IncludedSpecialWords string `json:"included_special_words"`
	// AudioLanguage : Hint for the spoken language of the source audio, as an
	// ISO 639-1 code (e.g. "en", "ja"). When empty, the service auto-detects
	// the language; supplying a hint improves accuracy and latency for short or
	// ambiguous clips. Languages the service does not support fall back to
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

// GetTranscriptAsyncCheckResult : Status of a transcript job started by
// `getTranscriptAsync`, as returned by `getTranscriptAsyncCheck`.
type GetTranscriptAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : The job finished successfully.
	Complete *GetTranscriptResult `json:"complete,omitempty"`
	// Failed : The job finished unsuccessfully.
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
		// Failed : The job finished unsuccessfully.
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
	// StructuredTranscript : The transcript produced for the requested media
	// asset.
	StructuredTranscript *ApiStructuredTranscript `json:"structured_transcript,omitempty"`
}

// NewGetTranscriptResult returns a new GetTranscriptResult instance
func NewGetTranscriptResult() *GetTranscriptResult {
	s := new(GetTranscriptResult)
	return s
}

// GetTransformAsyncCheckResult : Result type for EventBus async check - must
// end in "CheckResult"
type GetTransformAsyncCheckResult struct {
	dropbox.Tagged
	// Complete : The job finished successfully.
	Complete *TransformOutput `json:"complete,omitempty"`
	// Failed : The job finished unsuccessfully.
	Failed *TransformApiV2Error `json:"failed,omitempty"`
}

// Valid tag values for GetTransformAsyncCheckResult
const (
	GetTransformAsyncCheckResultInProgress = "in_progress"
	GetTransformAsyncCheckResultComplete   = "complete"
	GetTransformAsyncCheckResultFailed     = "failed"
	GetTransformAsyncCheckResultOther      = "other"
)

// UnmarshalJSON deserializes into a GetTransformAsyncCheckResult instance
func (u *GetTransformAsyncCheckResult) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// Failed : The job finished unsuccessfully.
		Failed *TransformApiV2Error `json:"failed,omitempty"`
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

// ImageOptions : Options for `TransformType.image` and
// `TransformType.image_pdf`. Supplying this message with any other transform
// type fails with `invalid_options_error`.
type ImageOptions struct {
	// PageNumber : For multi-page sources (PDFs, presentations, documents), the
	// 1-based page to render. Each request renders one page; to render a whole
	// document, issue one request per page. Defaults to the first page when
	// omitted.
	PageNumber uint32 `json:"page_number"`
	// ScalePercent : Scale the rendered image to this percentage of its natural
	// size. Must be in (0, 100] -- the pipeline does not upscale, so values
	// above 100 are rejected with `invalid_options_error`. Defaults to 100 (no
	// scaling) when omitted.
	ScalePercent uint32 `json:"scale_percent"`
}

// NewImageOptions returns a new ImageOptions instance
func NewImageOptions() *ImageOptions {
	s := new(ImageOptions)
	s.PageNumber = 1
	s.ScalePercent = 100
	return s
}

// KeyframesExtractionApiV2Error : Reason a keyframe extraction job failed.
// Returned in the `failed` variant of `GetKeyframesAsyncCheckResult`. This is a
// semantic error union: the HTTP status of the poll request itself is
// unaffected (a poll that surfaces a failed job is still a normal successful
// poll response). Callers should branch on the variant.
type KeyframesExtractionApiV2Error struct {
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

// Valid tag values for KeyframesExtractionApiV2Error
const (
	KeyframesExtractionApiV2ErrorServerError                 = "server_error"
	KeyframesExtractionApiV2ErrorUserError                   = "user_error"
	KeyframesExtractionApiV2ErrorUnsupportedFormatError      = "unsupported_format_error"
	KeyframesExtractionApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	KeyframesExtractionApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	KeyframesExtractionApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	KeyframesExtractionApiV2ErrorConversionFailureError      = "conversion_failure_error"
	KeyframesExtractionApiV2ErrorNotFoundError               = "not_found_error"
	KeyframesExtractionApiV2ErrorIsAFolderError              = "is_a_folder_error"
	KeyframesExtractionApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a KeyframesExtractionApiV2Error instance
func (u *KeyframesExtractionApiV2Error) UnmarshalJSON(body []byte) error {
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

// MarkdownConversionApiV2Error : Reason a markdown conversion job failed.
// Returned in the `GetMarkdownAsyncCheckResult.failed` variant. This is a
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
	// with the caller's input) -- for example an unsupported file format or a
	// file over the size limit. The string is a human-readable message;
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
		// with the caller's input) -- for example an unsupported file format or
		// a file over the size limit. The string is a human-readable message;
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
	// Limit : The maximum supported duration, in seconds, of the audio to
	// transcribe.
	Limit int32 `json:"limit"`
}

// NewMediaDurationError returns a new MediaDurationError instance
func NewMediaDurationError() *MediaDurationError {
	s := new(MediaDurationError)
	s.Limit = 0
	return s
}

// MetadataExtractionApiV2Error : Reason a metadata extraction job failed.
// Returned in the `GetMetadataAsyncCheckResult.failed` variant. This is a
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
	// with the caller's input) -- for example an unsupported file format or a
	// file over the size limit for its metadata kind. The string is a
	// human-readable message; retrying the same request will not help.
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
		// with the caller's input) -- for example an unsupported file format or
		// a file over the size limit for its metadata kind. The string is a
		// human-readable message; retrying the same request will not help.
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

// OcrExtractionApiV2Error : Reason an OCR extraction job failed. Returned in
// the `failed` variant of `GetOcrAsyncCheckResult`. This is a semantic error
// union: the HTTP status of the poll request itself is unaffected (a poll that
// surfaces a failed job is still a normal successful poll response). Callers
// should branch on the variant.
type OcrExtractionApiV2Error struct {
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

// Valid tag values for OcrExtractionApiV2Error
const (
	OcrExtractionApiV2ErrorServerError                 = "server_error"
	OcrExtractionApiV2ErrorUserError                   = "user_error"
	OcrExtractionApiV2ErrorUnsupportedFormatError      = "unsupported_format_error"
	OcrExtractionApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	OcrExtractionApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	OcrExtractionApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	OcrExtractionApiV2ErrorConversionFailureError      = "conversion_failure_error"
	OcrExtractionApiV2ErrorNotFoundError               = "not_found_error"
	OcrExtractionApiV2ErrorIsAFolderError              = "is_a_folder_error"
	OcrExtractionApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a OcrExtractionApiV2Error instance
func (u *OcrExtractionApiV2Error) UnmarshalJSON(body []byte) error {
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

// ThumbnailFormat : The encoding of the produced image. These match
// `files/get_thumbnail`'s formats: JPEG is the better choice for photographs,
// PNG for screenshots, line art, and anything with sharp text edges or
// transparency.
type ThumbnailFormat struct {
	dropbox.Tagged
}

// Valid tag values for ThumbnailFormat
const (
	ThumbnailFormatJpeg  = "jpeg"
	ThumbnailFormatPng   = "png"
	ThumbnailFormatWebp  = "webp"
	ThumbnailFormatOther = "other"
)

// ThumbnailMode : How to resize and crop the source to reach the requested
// `ThumbnailSize`. These match `files/get_thumbnail`'s modes.
type ThumbnailMode struct {
	dropbox.Tagged
}

// Valid tag values for ThumbnailMode
const (
	ThumbnailModeStrict        = "strict"
	ThumbnailModeBestfit       = "bestfit"
	ThumbnailModeFitoneBestfit = "fitone_bestfit"
	ThumbnailModeOriginal      = "original"
	ThumbnailModeOther         = "other"
)

// ThumbnailOptions : Options for `TransformType.thumbnail`. Supplying this
// message with any other transform type fails with `invalid_options_error`.
type ThumbnailOptions struct {
	// Size : The size bucket to produce. Defaults to `w64h64` when omitted.
	Size *ThumbnailSize `json:"size"`
	// Mode : How to fit the source into `size`. Defaults to `strict` when
	// omitted.
	Mode *ThumbnailMode `json:"mode"`
	// Format : The output encoding. Defaults to `jpeg` when omitted.
	Format *ThumbnailFormat `json:"format"`
}

// NewThumbnailOptions returns a new ThumbnailOptions instance
func NewThumbnailOptions() *ThumbnailOptions {
	s := new(ThumbnailOptions)
	s.Size = &ThumbnailSize{Tagged: dropbox.Tagged{Tag: "w64h64"}}
	s.Mode = &ThumbnailMode{Tagged: dropbox.Tagged{Tag: "strict"}}
	s.Format = &ThumbnailFormat{Tagged: dropbox.Tagged{Tag: "jpeg"}}
	return s
}

// ThumbnailSize : The size of the thumbnail to produce. These are the same
// named size buckets `files/get_thumbnail` supports, with the same meanings;
// arbitrary pixel dimensions are not accepted.
type ThumbnailSize struct {
	dropbox.Tagged
}

// Valid tag values for ThumbnailSize
const (
	ThumbnailSizeW32h32     = "w32h32"
	ThumbnailSizeW64h64     = "w64h64"
	ThumbnailSizeW128h128   = "w128h128"
	ThumbnailSizeW256h256   = "w256h256"
	ThumbnailSizeW480h320   = "w480h320"
	ThumbnailSizeW640h480   = "w640h480"
	ThumbnailSizeW960h640   = "w960h640"
	ThumbnailSizeW1024h768  = "w1024h768"
	ThumbnailSizeW2048h1536 = "w2048h1536"
	ThumbnailSizeOther      = "other"
)

// TimestampLevel : Granularity of the time offsets returned for each transcript
// segment.
type TimestampLevel struct {
	dropbox.Tagged
}

// Valid tag values for TimestampLevel
const (
	TimestampLevelSentence = "sentence"
	TimestampLevelWord     = "word"
	TimestampLevelOther    = "other"
)

// TransformApiV2Error : Reason a transform job failed. Returned in the `failed`
// variant of `GetTransformAsyncCheckResult`, and by
// `download_transform_output`. This is a semantic error union: the HTTP status
// of the poll request itself is unaffected (a poll that surfaces a failed job
// is still a normal successful poll response). Callers should branch on the
// variant.
type TransformApiV2Error struct {
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

// Valid tag values for TransformApiV2Error
const (
	TransformApiV2ErrorServerError                 = "server_error"
	TransformApiV2ErrorUserError                   = "user_error"
	TransformApiV2ErrorUnsupportedFormatError      = "unsupported_format_error"
	TransformApiV2ErrorLinkDownloadDisabledError   = "link_download_disabled_error"
	TransformApiV2ErrorSharedLinkPasswordProtected = "shared_link_password_protected"
	TransformApiV2ErrorLimitExceededError          = "limit_exceeded_error"
	TransformApiV2ErrorConversionFailureError      = "conversion_failure_error"
	TransformApiV2ErrorNotFoundError               = "not_found_error"
	TransformApiV2ErrorIsAFolderError              = "is_a_folder_error"
	TransformApiV2ErrorInvalidOptionsError         = "invalid_options_error"
	TransformApiV2ErrorExpiredHandleError          = "expired_handle_error"
	TransformApiV2ErrorOther                       = "other"
)

// UnmarshalJSON deserializes into a TransformApiV2Error instance
func (u *TransformApiV2Error) UnmarshalJSON(body []byte) error {
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

// TransformArgs : Arguments for the asynchronous `get_transform_async` route.
// Exactly one of `file_id`, `path`, or `url` must be supplied via
// `file_id_or_url` to identify the source file, and exactly one variant of
// `transform_type` must be set to say what to produce from it. At most one
// options message may be set, and it must be the one belonging to the requested
// `transform_type`. Options that belong to a different transform type are
// rejected with `invalid_options_error` rather than ignored, so that a request
// whose parameters were misassembled fails visibly instead of quietly producing
// the wrong output.
type TransformArgs struct {
	// FileIdOrUrl : Identifier of the source file to transform. Callers must
	// set exactly one of the `FileIdOrUrl` variants. The referenced file must
	// be in a format the requested `transform_type` supports; see the route
	// description for the per-transform format lists. Requests against
	// unsupported formats fail with `unsupported_format_error`.
	FileIdOrUrl *FileIdOrUrl `json:"file_id_or_url,omitempty"`
	// TransformType : What to produce from the source file. Required.
	TransformType *TransformType `json:"transform_type"`
	// Thumbnail : Options for `TransformType.thumbnail`.
	Thumbnail *ThumbnailOptions `json:"thumbnail,omitempty"`
	// Image : Options for `TransformType.image` and `TransformType.image_pdf`.
	Image *ImageOptions `json:"image,omitempty"`
	// VideoFrame : Options for `TransformType.video_frame`.
	VideoFrame *VideoFrameOptions `json:"video_frame,omitempty"`
}

// NewTransformArgs returns a new TransformArgs instance
func NewTransformArgs(TransformType *TransformType) *TransformArgs {
	s := new(TransformArgs)
	s.TransformType = TransformType
	return s
}

// TransformOutput : A completed transform: a handle for retrieving the produced
// bytes, plus enough metadata to decide whether to retrieve them. The bytes
// themselves are deliberately not carried here. A completed async result is
// persisted, so it is bounded by a row-size limit well below the size of a
// typical converted document -- an inline payload would fail for exactly the
// large documents this route exists to convert. The transform therefore
// completes by caching its output and handing back `output_handle`, which
// `download_transform_output` exchanges for the bytes.
type TransformOutput struct {
	// OutputHandle : Opaque, single-purpose handle for the produced bytes. Pass
	// it to `download_transform_output` to retrieve them. The handle is scoped
	// to the account that created it and cannot be used to read anything other
	// than the output of this transform. It is not a URL and carries no meaning
	// for callers beyond being passed back verbatim.
	OutputHandle string `json:"output_handle"`
	// Size : Size of the produced output in bytes.
	Size uint64 `json:"size"`
	// Format : Format of the produced output, as a short lowercase format name
	// such as "pdf", "html", "jpeg", or "png". This reflects what was actually
	// produced, which for some sources differs from what was requested.
	Format string `json:"format"`
	// MimeType : MIME type corresponding to `format`, for callers that need a
	// Content-Type to hand to a downstream consumer.
	MimeType string `json:"mime_type"`
	// ExpiresTs : Unix timestamp, in seconds, after which `output_handle` is no
	// longer accepted. Retrieve the bytes before this point; after it,
	// `download_transform_output` reports `expired_handle_error` and the
	// transform has to be requested again.
	ExpiresTs uint64 `json:"expires_ts"`
}

// NewTransformOutput returns a new TransformOutput instance
func NewTransformOutput() *TransformOutput {
	s := new(TransformOutput)
	s.OutputHandle = ""
	s.Size = 0
	s.Format = ""
	s.MimeType = ""
	s.ExpiresTs = 0
	return s
}

// TransformType : Which derived file to produce from the source file. Unlike
// the other Riviera content routes, which each expose one capability,
// `get_transform_async` is a single route over many conversions: the caller
// names the source file and the output it wants, and the service picks the
// conversion pipeline. Exactly one variant must be set; a request with none set
// fails with `invalid_options_error`.
type TransformType struct {
	dropbox.Tagged
}

// Valid tag values for TransformType
const (
	TransformTypePdf        = "pdf"
	TransformTypeHtml       = "html"
	TransformTypeImage      = "image"
	TransformTypeThumbnail  = "thumbnail"
	TransformTypeImagePdf   = "image_pdf"
	TransformTypeVideoFrame = "video_frame"
	TransformTypeOther      = "other"
)

// VideoFrameOptions : Options for `TransformType.video_frame`. Supplying this
// message with any other transform type fails with `invalid_options_error`.
type VideoFrameOptions struct {
	// OffsetInSeconds : Offset into the video, in seconds, of the frame to
	// extract. Should be within the video's duration. Defaults to 0 (the first
	// frame) when omitted. Only the lower bound is enforced, because the upper
	// bound is the source's duration, which is not known until the video is
	// opened. An offset past the end is not rejected.
	OffsetInSeconds float64 `json:"offset_in_seconds"`
	// ScalePercent : Scale the extracted frame to this percentage of the
	// video's natural frame size. Must be in (0, 100]; the pipeline does not
	// upscale. Defaults to 100 (no scaling) when omitted.
	ScalePercent uint32 `json:"scale_percent"`
}

// NewVideoFrameOptions returns a new VideoFrameOptions instance
func NewVideoFrameOptions() *VideoFrameOptions {
	s := new(VideoFrameOptions)
	s.OffsetInSeconds = 0.0
	s.ScalePercent = 100
	return s
}

// MetadataUnion : The extracted metadata. Exactly one variant is populated,
// corresponding to `GetMetadataResult.metadata_type`.
type MetadataUnion struct {
	dropbox.Tagged
	// Exif : EXIF metadata, for image files.
	Exif *ApiExifMetadata `json:"exif,omitempty"`
	// Media : Container and per-stream metadata, for audio and video files.
	Media *ApiMediaMetadata `json:"media,omitempty"`
	// Pdf : Document metadata, for PDFs.
	Pdf *ApiPdfMetadata `json:"pdf,omitempty"`
	// Office : Document metadata, for MS Office files.
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
