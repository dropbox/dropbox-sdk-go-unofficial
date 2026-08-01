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

package riviera

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/async"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// GetMarkdownAsync : Asynchronous document-to-markdown conversion for
	// supported file formats. Supported formats: .binder, .docx, .html, .paper,
	// .papert, .pptx, .xlsx, .gsheet, .ods, .pdf. Unsupported formats return an
	// `unsupported_format_error`. Size limit: the source file must be at most
	// 50 MB. Larger files are rejected.
	GetMarkdownAsync(arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error)
	// GetMarkdownAsyncCheck : Returns the status or result of specified
	// get_markdown_async task.
	GetMarkdownAsyncCheck(arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error)
	// GetMetadataAsync : Asynchronous file metadata extraction for supported
	// file formats. The kind of metadata returned depends on the file type: -
	// Image (EXIF) formats: .3fr, .arw, .avif, .bmp, .cr2, .cr3, .crw, .dcr,
	// .dcs, .dng, .erf, .gif, .heic, .j2c, .j2k, .jp2, .jpc, .jpeg, .jpf, .jpg,
	// .jpg2, .jpm, .jpx, .kdc, .mef, .mos, .mrw, .nef, .nrw, .orf, .pef, .png,
	// .ppm, .r3d, .raf, .rw2, .rwl, .sr2, .tga, .tif, .tiff, .wbmp, .web,
	// .webp, .x3f. - Audio/video (media) formats: .aac, .aif, .aiff, .flac,
	// .m4a, .m4r, .mp3, .oga, .ogg, .wav, .wma, .3gp, .3gpp, .3gpp2, .asf,
	// .avi, .dv, .flv, .m2t, .m2ts, .m4v, .mkv, .mov, .mp4, .mpeg, .mpg, .mts,
	// .mxf, .oggtheora, .ogv, .rm, .ts, .vob, .webm, .wmv. - PDF format: .pdf.
	// - MS Office formats: .docx, .pptx, .xlsx. Unsupported formats return an
	// `unsupported_format_error`.
	GetMetadataAsync(arg *GetMetadataArgs) (res *async.LaunchResultBase, err error)
	// GetMetadataAsyncCheck : Returns the status or result of specified
	// get_metadata_async task.
	GetMetadataAsyncCheck(arg *async.PollArg) (res *GetMetadataAsyncCheckResult, err error)
	// GetTextAsync : Asynchronous plain-text extraction from documents.
	// Supported formats include: - Word processing: .doc, .docx, .docm, .rtf. -
	// Presentations: .ppt, .pptx, .pptm. - Spreadsheets: .xls, .xlsx, .xlsm. -
	// PDF: .pdf. - Dropbox document types: .paper, .papert, .binder, .gdoc,
	// .gsheet, .gslides. - Plain text / subtitles: .txt, .vtt. Unsupported
	// formats return an `unsupported_format_error`. For the `url` variant only
	// Dropbox shared links are supported; external URLs return
	// `unsupported_format_error`.
	GetTextAsync(arg *GetTextArgs) (res *async.LaunchResultBase, err error)
	// GetTextAsyncCheck : Returns the status or result of specified
	// get_text_async task.
	GetTextAsyncCheck(arg *async.PollArg) (res *GetTextAsyncCheckResult, err error)
	// GetTranscriptAsync : Asynchronous transcript generation for audio and
	// video files. Supported audio formats: .aac, .aif, .aiff, .flac, .m4a,
	// .m4r, .mp3, .oga, .ogg, .wav, .wma. Supported video formats: .3gp, .3gpp,
	// .3gpp2, .asf, .avi, .dv, .flv, .m2t, .m2ts, .m4v, .mkv, .mov, .mp4,
	// .mpeg, .mpg, .mts, .mxf, .oggtheora, .ogv, .rm, .ts, .vob, .webm, .wmv.
	// Unsupported formats return an `unsupported_format_error`. Size limits:
	// the source file must be at most 10 GB and its audio track at most 1 hour
	// in duration. Files exceeding these limits are rejected.
	GetTranscriptAsync(arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error)
	// GetTranscriptAsyncCheck : Returns the status or result of specified
	// get_transcript_async task.
	GetTranscriptAsyncCheck(arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error)
}

// ContextClient interface describes all routes in this namespace with context support
type ContextClient interface {
	Client
	// GetMarkdownAsyncContext : Asynchronous document-to-markdown conversion
	// for supported file formats. Supported formats: .binder, .docx, .html,
	// .paper, .papert, .pptx, .xlsx, .gsheet, .ods, .pdf. Unsupported formats
	// return an `unsupported_format_error`. Size limit: the source file must be
	// at most 50 MB. Larger files are rejected.
	GetMarkdownAsyncContext(ctx context.Context, arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error)
	// GetMarkdownAsyncCheckContext : Returns the status or result of specified
	// get_markdown_async task.
	GetMarkdownAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error)
	// GetMetadataAsyncContext : Asynchronous file metadata extraction for
	// supported file formats. The kind of metadata returned depends on the file
	// type: - Image (EXIF) formats: .3fr, .arw, .avif, .bmp, .cr2, .cr3, .crw,
	// .dcr, .dcs, .dng, .erf, .gif, .heic, .j2c, .j2k, .jp2, .jpc, .jpeg, .jpf,
	// .jpg, .jpg2, .jpm, .jpx, .kdc, .mef, .mos, .mrw, .nef, .nrw, .orf, .pef,
	// .png, .ppm, .r3d, .raf, .rw2, .rwl, .sr2, .tga, .tif, .tiff, .wbmp, .web,
	// .webp, .x3f. - Audio/video (media) formats: .aac, .aif, .aiff, .flac,
	// .m4a, .m4r, .mp3, .oga, .ogg, .wav, .wma, .3gp, .3gpp, .3gpp2, .asf,
	// .avi, .dv, .flv, .m2t, .m2ts, .m4v, .mkv, .mov, .mp4, .mpeg, .mpg, .mts,
	// .mxf, .oggtheora, .ogv, .rm, .ts, .vob, .webm, .wmv. - PDF format: .pdf.
	// - MS Office formats: .docx, .pptx, .xlsx. Unsupported formats return an
	// `unsupported_format_error`.
	GetMetadataAsyncContext(ctx context.Context, arg *GetMetadataArgs) (res *async.LaunchResultBase, err error)
	// GetMetadataAsyncCheckContext : Returns the status or result of specified
	// get_metadata_async task.
	GetMetadataAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetMetadataAsyncCheckResult, err error)
	// GetTextAsyncContext : Asynchronous plain-text extraction from documents.
	// Supported formats include: - Word processing: .doc, .docx, .docm, .rtf. -
	// Presentations: .ppt, .pptx, .pptm. - Spreadsheets: .xls, .xlsx, .xlsm. -
	// PDF: .pdf. - Dropbox document types: .paper, .papert, .binder, .gdoc,
	// .gsheet, .gslides. - Plain text / subtitles: .txt, .vtt. Unsupported
	// formats return an `unsupported_format_error`. For the `url` variant only
	// Dropbox shared links are supported; external URLs return
	// `unsupported_format_error`.
	GetTextAsyncContext(ctx context.Context, arg *GetTextArgs) (res *async.LaunchResultBase, err error)
	// GetTextAsyncCheckContext : Returns the status or result of specified
	// get_text_async task.
	GetTextAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetTextAsyncCheckResult, err error)
	// GetTranscriptAsyncContext : Asynchronous transcript generation for audio
	// and video files. Supported audio formats: .aac, .aif, .aiff, .flac, .m4a,
	// .m4r, .mp3, .oga, .ogg, .wav, .wma. Supported video formats: .3gp, .3gpp,
	// .3gpp2, .asf, .avi, .dv, .flv, .m2t, .m2ts, .m4v, .mkv, .mov, .mp4,
	// .mpeg, .mpg, .mts, .mxf, .oggtheora, .ogv, .rm, .ts, .vob, .webm, .wmv.
	// Unsupported formats return an `unsupported_format_error`. Size limits:
	// the source file must be at most 10 GB and its audio track at most 1 hour
	// in duration. Files exceeding these limits are rejected.
	GetTranscriptAsyncContext(ctx context.Context, arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error)
	// GetTranscriptAsyncCheckContext : Returns the status or result of
	// specified get_transcript_async task.
	GetTranscriptAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error)
}

type apiImpl dropbox.Context

// GetMarkdownAsyncAPIError is an error-wrapper for the get_markdown_async route
type GetMarkdownAsyncAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// GetMarkdownAsyncContext : Asynchronous document-to-markdown conversion for
// supported file formats. Supported formats: .binder, .docx, .html, .paper,
// .papert, .pptx, .xlsx, .gsheet, .ods, .pdf. Unsupported formats return an
// `unsupported_format_error`. Size limit: the source file must be at most 50
// MB. Larger files are rejected.
func (dbx *apiImpl) GetMarkdownAsyncContext(ctx context.Context, arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_markdown_async",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetMarkdownAsyncAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetMarkdownAsync(arg *GetMarkdownArgs) (res *async.LaunchResultBase, err error) {
	return dbx.GetMarkdownAsyncContext(context.Background(), arg)
}

// GetMarkdownAsyncCheckAPIError is an error-wrapper for the get_markdown_async/check route
type GetMarkdownAsyncCheckAPIError struct {
	dropbox.APIError
	EndpointError *async.PollError `json:"error"`
}

// GetMarkdownAsyncCheckContext : Returns the status or result of specified
// get_markdown_async task.
func (dbx *apiImpl) GetMarkdownAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_markdown_async/check",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetMarkdownAsyncCheckAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetMarkdownAsyncCheck(arg *async.PollArg) (res *GetMarkdownAsyncCheckResult, err error) {
	return dbx.GetMarkdownAsyncCheckContext(context.Background(), arg)
}

// GetMetadataAsyncAPIError is an error-wrapper for the get_metadata_async route
type GetMetadataAsyncAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// GetMetadataAsyncContext : Asynchronous file metadata extraction for supported
// file formats. The kind of metadata returned depends on the file type: - Image
// (EXIF) formats: .3fr, .arw, .avif, .bmp, .cr2, .cr3, .crw, .dcr, .dcs, .dng,
// .erf, .gif, .heic, .j2c, .j2k, .jp2, .jpc, .jpeg, .jpf, .jpg, .jpg2, .jpm,
// .jpx, .kdc, .mef, .mos, .mrw, .nef, .nrw, .orf, .pef, .png, .ppm, .r3d, .raf,
// .rw2, .rwl, .sr2, .tga, .tif, .tiff, .wbmp, .web, .webp, .x3f. - Audio/video
// (media) formats: .aac, .aif, .aiff, .flac, .m4a, .m4r, .mp3, .oga, .ogg,
// .wav, .wma, .3gp, .3gpp, .3gpp2, .asf, .avi, .dv, .flv, .m2t, .m2ts, .m4v,
// .mkv, .mov, .mp4, .mpeg, .mpg, .mts, .mxf, .oggtheora, .ogv, .rm, .ts, .vob,
// .webm, .wmv. - PDF format: .pdf. - MS Office formats: .docx, .pptx, .xlsx.
// Unsupported formats return an `unsupported_format_error`.
func (dbx *apiImpl) GetMetadataAsyncContext(ctx context.Context, arg *GetMetadataArgs) (res *async.LaunchResultBase, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_metadata_async",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetMetadataAsyncAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetMetadataAsync(arg *GetMetadataArgs) (res *async.LaunchResultBase, err error) {
	return dbx.GetMetadataAsyncContext(context.Background(), arg)
}

// GetMetadataAsyncCheckAPIError is an error-wrapper for the get_metadata_async/check route
type GetMetadataAsyncCheckAPIError struct {
	dropbox.APIError
	EndpointError *async.PollError `json:"error"`
}

// GetMetadataAsyncCheckContext : Returns the status or result of specified
// get_metadata_async task.
func (dbx *apiImpl) GetMetadataAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetMetadataAsyncCheckResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_metadata_async/check",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetMetadataAsyncCheckAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetMetadataAsyncCheck(arg *async.PollArg) (res *GetMetadataAsyncCheckResult, err error) {
	return dbx.GetMetadataAsyncCheckContext(context.Background(), arg)
}

// GetTextAsyncAPIError is an error-wrapper for the get_text_async route
type GetTextAsyncAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// GetTextAsyncContext : Asynchronous plain-text extraction from documents.
// Supported formats include: - Word processing: .doc, .docx, .docm, .rtf. -
// Presentations: .ppt, .pptx, .pptm. - Spreadsheets: .xls, .xlsx, .xlsm. - PDF:
// .pdf. - Dropbox document types: .paper, .papert, .binder, .gdoc, .gsheet,
// .gslides. - Plain text / subtitles: .txt, .vtt. Unsupported formats return an
// `unsupported_format_error`. For the `url` variant only Dropbox shared links
// are supported; external URLs return `unsupported_format_error`.
func (dbx *apiImpl) GetTextAsyncContext(ctx context.Context, arg *GetTextArgs) (res *async.LaunchResultBase, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_text_async",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetTextAsyncAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetTextAsync(arg *GetTextArgs) (res *async.LaunchResultBase, err error) {
	return dbx.GetTextAsyncContext(context.Background(), arg)
}

// GetTextAsyncCheckAPIError is an error-wrapper for the get_text_async/check route
type GetTextAsyncCheckAPIError struct {
	dropbox.APIError
	EndpointError *async.PollError `json:"error"`
}

// GetTextAsyncCheckContext : Returns the status or result of specified
// get_text_async task.
func (dbx *apiImpl) GetTextAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetTextAsyncCheckResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_text_async/check",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetTextAsyncCheckAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetTextAsyncCheck(arg *async.PollArg) (res *GetTextAsyncCheckResult, err error) {
	return dbx.GetTextAsyncCheckContext(context.Background(), arg)
}

// GetTranscriptAsyncAPIError is an error-wrapper for the get_transcript_async route
type GetTranscriptAsyncAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

// GetTranscriptAsyncContext : Asynchronous transcript generation for audio and
// video files. Supported audio formats: .aac, .aif, .aiff, .flac, .m4a, .m4r,
// .mp3, .oga, .ogg, .wav, .wma. Supported video formats: .3gp, .3gpp, .3gpp2,
// .asf, .avi, .dv, .flv, .m2t, .m2ts, .m4v, .mkv, .mov, .mp4, .mpeg, .mpg,
// .mts, .mxf, .oggtheora, .ogv, .rm, .ts, .vob, .webm, .wmv. Unsupported
// formats return an `unsupported_format_error`. Size limits: the source file
// must be at most 10 GB and its audio track at most 1 hour in duration. Files
// exceeding these limits are rejected.
func (dbx *apiImpl) GetTranscriptAsyncContext(ctx context.Context, arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_transcript_async",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetTranscriptAsyncAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetTranscriptAsync(arg *GetTranscriptArgs) (res *async.LaunchResultBase, err error) {
	return dbx.GetTranscriptAsyncContext(context.Background(), arg)
}

// GetTranscriptAsyncCheckAPIError is an error-wrapper for the get_transcript_async/check route
type GetTranscriptAsyncCheckAPIError struct {
	dropbox.APIError
	EndpointError *async.PollError `json:"error"`
}

// GetTranscriptAsyncCheckContext : Returns the status or result of specified
// get_transcript_async task.
func (dbx *apiImpl) GetTranscriptAsyncCheckContext(ctx context.Context, arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "riviera",
		Route:        "get_transcript_async/check",
		Auth:         "app, user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr GetTranscriptAsyncCheckAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
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

func (dbx *apiImpl) GetTranscriptAsyncCheck(arg *async.PollArg) (res *GetTranscriptAsyncCheckResult, err error) {
	return dbx.GetTranscriptAsyncCheckContext(context.Background(), arg)
}

// NewContext returns a ContextClient implementation for this namespace
func NewContext(c dropbox.Config) ContextClient {
	ctx := apiImpl(dropbox.NewContext(c))
	return &ctx
}

// New returns a Client implementation for this namespace
func New(c dropbox.Config) Client {
	return NewContext(c)
}
