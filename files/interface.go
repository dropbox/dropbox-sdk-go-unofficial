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

// This namespace contains endpoints and data types for basic file operations.
package files

import (
	"encoding/json"
	"io"
	"time"
)

type CommitInfo struct {
	// Path in the user's Dropbox to save the file.
	Path string `json:"path"`
	// Selects what to do if the file already exists.
	Mode *WriteMode `json:"mode"`
	// If there's a conflict, as determined by `mode`, have the Dropbox server try
	// to autorename the file to avoid conflict.
	Autorename bool `json:"autorename"`
	// The value to store as the `client_modified` timestamp. Dropbox automatically
	// records the time at which the file was written to the Dropbox servers. It
	// can also record an additional timestamp, provided by Dropbox desktop
	// clients, mobile clients, and API apps of when the file was actually created
	// or modified.
	ClientModified time.Time `json:"client_modified,omitempty"`
	// Normally, users are made aware of any file modifications in their Dropbox
	// account via notifications in the client software. If `True`, this tells the
	// clients that this modification shouldn't result in a user notification.
	Mute bool `json:"mute"`
}

func NewCommitInfo() *CommitInfo {
	s := new(CommitInfo)
	s.Mode = &WriteMode{Tag: "add"}
	s.Autorename = false
	s.Mute = false
	return s
}

type CreateFolderArg struct {
	// Path in the user's Dropbox to create.
	Path string `json:"path"`
}

func NewCreateFolderArg() *CreateFolderArg {
	s := new(CreateFolderArg)
	return s
}

type CreateFolderError struct {
	Tag  string      `json:".tag"`
	Path *WriteError `json:"path,omitempty"`
}

func (u *CreateFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type DeleteArg struct {
	// Path in the user's Dropbox to delete.
	Path string `json:"path"`
}

func NewDeleteArg() *DeleteArg {
	s := new(DeleteArg)
	return s
}

type DeleteError struct {
	Tag        string       `json:".tag"`
	PathLookup *LookupError `json:"path_lookup,omitempty"`
	PathWrite  *WriteError  `json:"path_write,omitempty"`
}

func (u *DeleteError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag        string          `json:".tag"`
		PathLookup json.RawMessage `json:"path_lookup"`
		PathWrite  json.RawMessage `json:"path_write"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path_lookup":
		{
			if len(w.PathLookup) == 0 {
				break
			}
			if err := json.Unmarshal(w.PathLookup, &u.PathLookup); err != nil {
				return err
			}
		}
	case "path_write":
		{
			if len(w.PathWrite) == 0 {
				break
			}
			if err := json.Unmarshal(w.PathWrite, &u.PathWrite); err != nil {
				return err
			}
		}
	}
	return nil
}

// Metadata for a file or folder.
type Metadata struct {
	Tag     string           `json:".tag"`
	File    *FileMetadata    `json:"file,omitempty"`
	Folder  *FolderMetadata  `json:"folder,omitempty"`
	Deleted *DeletedMetadata `json:"deleted,omitempty"`
}

func (u *Metadata) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag     string          `json:".tag"`
		File    json.RawMessage `json:"file"`
		Folder  json.RawMessage `json:"folder"`
		Deleted json.RawMessage `json:"deleted"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "file":
		{
			if err := json.Unmarshal(body, &u.File); err != nil {
				return err
			}
		}
	case "folder":
		{
			if err := json.Unmarshal(body, &u.Folder); err != nil {
				return err
			}
		}
	case "deleted":
		{
			if err := json.Unmarshal(body, &u.Deleted); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMetadata() *Metadata {
	s := new(Metadata)
	return s
}

// Indicates that there used to be a file or folder at this path, but it no
// longer exists.
type DeletedMetadata struct {
	// The last component of the path (including extension). This never contains a
	// slash.
	Name string `json:"name"`
	// The lowercased full path in the user's Dropbox. This always starts with a
	// slash.
	PathLower string `json:"path_lower"`
	// Deprecated. Please use :field:'FileSharingInfo.parent_shared_folder_id' or
	// :field:'FolderSharingInfo.parent_shared_folder_id' instead.
	ParentSharedFolderId string `json:"parent_shared_folder_id,omitempty"`
}

func NewDeletedMetadata() *DeletedMetadata {
	s := new(DeletedMetadata)
	return s
}

// Dimensions for a photo or video.
type Dimensions struct {
	// Height of the photo/video.
	Height uint64 `json:"height"`
	// Width of the photo/video.
	Width uint64 `json:"width"`
}

func NewDimensions() *Dimensions {
	s := new(Dimensions)
	return s
}

type DownloadArg struct {
	// The path of the file to download.
	Path string `json:"path"`
	// Deprecated. Please specify revision in :field:'path' instead
	Rev string `json:"rev,omitempty"`
}

func NewDownloadArg() *DownloadArg {
	s := new(DownloadArg)
	return s
}

type DownloadError struct {
	Tag  string       `json:".tag"`
	Path *LookupError `json:"path,omitempty"`
}

func (u *DownloadError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type FileMetadata struct {
	// The last component of the path (including extension). This never contains a
	// slash.
	Name string `json:"name"`
	// The lowercased full path in the user's Dropbox. This always starts with a
	// slash.
	PathLower string `json:"path_lower"`
	// For files, this is the modification time set by the desktop client when the
	// file was added to Dropbox. Since this time is not verified (the Dropbox
	// server stores whatever the desktop client sends up), this should only be
	// used for display purposes (such as sorting) and not, for example, to
	// determine if a file has changed or not.
	ClientModified time.Time `json:"client_modified"`
	// The last time the file was modified on Dropbox.
	ServerModified time.Time `json:"server_modified"`
	// A unique identifier for the current revision of a file. This field is the
	// same rev as elsewhere in the API and can be used to detect changes and avoid
	// conflicts.
	Rev string `json:"rev"`
	// The file size in bytes.
	Size uint64 `json:"size"`
	// Deprecated. Please use :field:'FileSharingInfo.parent_shared_folder_id' or
	// :field:'FolderSharingInfo.parent_shared_folder_id' instead.
	ParentSharedFolderId string `json:"parent_shared_folder_id,omitempty"`
	// A unique identifier for the file.
	Id string `json:"id,omitempty"`
	// Additional information if the file is a photo or video.
	MediaInfo *MediaInfo `json:"media_info,omitempty"`
	// Set if this file is contained in a shared folder.
	SharingInfo *FileSharingInfo `json:"sharing_info,omitempty"`
}

func NewFileMetadata() *FileMetadata {
	s := new(FileMetadata)
	return s
}

// Sharing info for a file or folder.
type SharingInfo struct {
	// True if the file or folder is inside a read-only shared folder.
	ReadOnly bool `json:"read_only"`
}

func NewSharingInfo() *SharingInfo {
	s := new(SharingInfo)
	return s
}

// Sharing info for a file which is contained by a shared folder.
type FileSharingInfo struct {
	// True if the file or folder is inside a read-only shared folder.
	ReadOnly bool `json:"read_only"`
	// ID of shared folder that holds this file.
	ParentSharedFolderId string `json:"parent_shared_folder_id"`
	// The last user who modified the file. This field will be null if the user's
	// account has been deleted.
	ModifiedBy string `json:"modified_by,omitempty"`
}

func NewFileSharingInfo() *FileSharingInfo {
	s := new(FileSharingInfo)
	return s
}

type FolderMetadata struct {
	// The last component of the path (including extension). This never contains a
	// slash.
	Name string `json:"name"`
	// The lowercased full path in the user's Dropbox. This always starts with a
	// slash.
	PathLower string `json:"path_lower"`
	// Deprecated. Please use :field:'FileSharingInfo.parent_shared_folder_id' or
	// :field:'FolderSharingInfo.parent_shared_folder_id' instead.
	ParentSharedFolderId string `json:"parent_shared_folder_id,omitempty"`
	// A unique identifier for the folder.
	Id string `json:"id,omitempty"`
	// Deprecated. Please use :field:'sharing_info' instead.
	SharedFolderId string `json:"shared_folder_id,omitempty"`
	// Set if the folder is contained in a shared folder or is a shared folder
	// mount point.
	SharingInfo *FolderSharingInfo `json:"sharing_info,omitempty"`
}

func NewFolderMetadata() *FolderMetadata {
	s := new(FolderMetadata)
	return s
}

// Sharing info for a folder which is contained in a shared folder or is a
// shared folder mount point.
type FolderSharingInfo struct {
	// True if the file or folder is inside a read-only shared folder.
	ReadOnly bool `json:"read_only"`
	// Set if the folder is contained by a shared folder.
	ParentSharedFolderId string `json:"parent_shared_folder_id,omitempty"`
	// If this folder is a shared folder mount point, the ID of the shared folder
	// mounted at this location.
	SharedFolderId string `json:"shared_folder_id,omitempty"`
}

func NewFolderSharingInfo() *FolderSharingInfo {
	s := new(FolderSharingInfo)
	return s
}

type GetMetadataArg struct {
	// The path of a file or folder on Dropbox
	Path string `json:"path"`
	// If true, :field:'FileMetadata.media_info' is set for photo and video.
	IncludeMediaInfo bool `json:"include_media_info"`
}

func NewGetMetadataArg() *GetMetadataArg {
	s := new(GetMetadataArg)
	s.IncludeMediaInfo = false
	return s
}

type GetMetadataError struct {
	Tag  string       `json:".tag"`
	Path *LookupError `json:"path,omitempty"`
}

func (u *GetMetadataError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

// GPS coordinates for a photo or video.
type GpsCoordinates struct {
	// Latitude of the GPS coordinates.
	Latitude float64 `json:"latitude"`
	// Longitude of the GPS coordinates.
	Longitude float64 `json:"longitude"`
}

func NewGpsCoordinates() *GpsCoordinates {
	s := new(GpsCoordinates)
	return s
}

type ListFolderArg struct {
	// The path to the folder you want to see the contents of.
	Path string `json:"path"`
	// If true, the list folder operation will be applied recursively to all
	// subfolders and the response will contain contents of all subfolders.
	Recursive bool `json:"recursive"`
	// If true, :field:'FileMetadata.media_info' is set for photo and video.
	IncludeMediaInfo bool `json:"include_media_info"`
	// If true, the results will include entries for files and folders that used to
	// exist but were deleted.
	IncludeDeleted bool `json:"include_deleted"`
}

func NewListFolderArg() *ListFolderArg {
	s := new(ListFolderArg)
	s.Recursive = false
	s.IncludeMediaInfo = false
	s.IncludeDeleted = false
	return s
}

type ListFolderContinueArg struct {
	// The cursor returned by your last call to `ListFolder` or
	// `ListFolderContinue`.
	Cursor string `json:"cursor"`
}

func NewListFolderContinueArg() *ListFolderContinueArg {
	s := new(ListFolderContinueArg)
	return s
}

type ListFolderContinueError struct {
	Tag  string       `json:".tag"`
	Path *LookupError `json:"path,omitempty"`
}

func (u *ListFolderContinueError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type ListFolderError struct {
	Tag  string       `json:".tag"`
	Path *LookupError `json:"path,omitempty"`
}

func (u *ListFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type ListFolderGetLatestCursorResult struct {
	// Pass the cursor into `ListFolderContinue` to see what's changed in the
	// folder since your previous query.
	Cursor string `json:"cursor"`
}

func NewListFolderGetLatestCursorResult() *ListFolderGetLatestCursorResult {
	s := new(ListFolderGetLatestCursorResult)
	return s
}

type ListFolderLongpollArg struct {
	// A cursor as returned by `ListFolder` or `ListFolderContinue`
	Cursor string `json:"cursor"`
	// A timeout in seconds. The request will block for at most this length of
	// time, plus up to 90 seconds of random jitter added to avoid the thundering
	// herd problem. Care should be taken when using this parameter, as some
	// network infrastructure does not support long timeouts.
	Timeout uint64 `json:"timeout"`
}

func NewListFolderLongpollArg() *ListFolderLongpollArg {
	s := new(ListFolderLongpollArg)
	s.Timeout = 30
	return s
}

type ListFolderLongpollError struct {
	Tag string `json:".tag"`
}

type ListFolderLongpollResult struct {
	// Indicates whether new changes are available. If true, call `ListFolder` to
	// retrieve the changes.
	Changes bool `json:"changes"`
	// If present, backoff for at least this many seconds before calling
	// `ListFolderLongpoll` again.
	Backoff uint64 `json:"backoff,omitempty"`
}

func NewListFolderLongpollResult() *ListFolderLongpollResult {
	s := new(ListFolderLongpollResult)
	return s
}

type ListFolderResult struct {
	// The files and (direct) subfolders in the folder.
	Entries []*Metadata `json:"entries"`
	// Pass the cursor into `ListFolderContinue` to see what's changed in the
	// folder since your previous query.
	Cursor string `json:"cursor"`
	// If true, then there are more entries available. Pass the cursor to
	// `ListFolderContinue` to retrieve the rest.
	HasMore bool `json:"has_more"`
}

func NewListFolderResult() *ListFolderResult {
	s := new(ListFolderResult)
	return s
}

type ListRevisionsArg struct {
	// The path to the file you want to see the revisions of.
	Path string `json:"path"`
	// The maximum number of revision entries returned.
	Limit uint64 `json:"limit"`
}

func NewListRevisionsArg() *ListRevisionsArg {
	s := new(ListRevisionsArg)
	s.Limit = 10
	return s
}

type ListRevisionsError struct {
	Tag  string       `json:".tag"`
	Path *LookupError `json:"path,omitempty"`
}

func (u *ListRevisionsError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type ListRevisionsResult struct {
	// If the file is deleted.
	IsDeleted bool `json:"is_deleted"`
	// The revisions for the file. Only non-delete revisions will show up here.
	Entries []*FileMetadata `json:"entries"`
}

func NewListRevisionsResult() *ListRevisionsResult {
	s := new(ListRevisionsResult)
	return s
}

type LookupError struct {
	Tag           string `json:".tag"`
	MalformedPath string `json:"malformed_path,omitempty"`
}

func (u *LookupError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag           string          `json:".tag"`
		MalformedPath json.RawMessage `json:"malformed_path,omitempty"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "malformed_path":
		{
			if len(w.MalformedPath) == 0 {
				break
			}
			if err := json.Unmarshal(w.MalformedPath, &u.MalformedPath); err != nil {
				return err
			}
		}
	}
	return nil
}

type MediaInfo struct {
	Tag string `json:".tag"`
	// The metadata for the photo/video.
	Metadata *MediaMetadata `json:"metadata,omitempty"`
}

func (u *MediaInfo) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The metadata for the photo/video.
		Metadata json.RawMessage `json:"metadata"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "metadata":
		{
			if err := json.Unmarshal(body, &u.Metadata); err != nil {
				return err
			}
		}
	}
	return nil
}

// Metadata for a photo or video.
type MediaMetadata struct {
	Tag   string         `json:".tag"`
	Photo *PhotoMetadata `json:"photo,omitempty"`
	Video *VideoMetadata `json:"video,omitempty"`
}

func (u *MediaMetadata) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag   string          `json:".tag"`
		Photo json.RawMessage `json:"photo"`
		Video json.RawMessage `json:"video"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "photo":
		{
			if err := json.Unmarshal(body, &u.Photo); err != nil {
				return err
			}
		}
	case "video":
		{
			if err := json.Unmarshal(body, &u.Video); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMediaMetadata() *MediaMetadata {
	s := new(MediaMetadata)
	return s
}

// Metadata for a photo.
type PhotoMetadata struct {
	// Dimension of the photo/video.
	Dimensions *Dimensions `json:"dimensions,omitempty"`
	// The GPS coordinate of the photo/video.
	Location *GpsCoordinates `json:"location,omitempty"`
	// The timestamp when the photo/video is taken.
	TimeTaken time.Time `json:"time_taken,omitempty"`
}

func NewPhotoMetadata() *PhotoMetadata {
	s := new(PhotoMetadata)
	return s
}

type PreviewArg struct {
	// The path of the file to preview.
	Path string `json:"path"`
	// Deprecated. Please specify revision in :field:'path' instead
	Rev string `json:"rev,omitempty"`
}

func NewPreviewArg() *PreviewArg {
	s := new(PreviewArg)
	return s
}

type PreviewError struct {
	Tag string `json:".tag"`
	// An error occurs when downloading metadata for the file.
	Path *LookupError `json:"path,omitempty"`
}

func (u *PreviewError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// An error occurs when downloading metadata for the file.
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type RelocationArg struct {
	// Path in the user's Dropbox to be copied or moved.
	FromPath string `json:"from_path"`
	// Path in the user's Dropbox that is the destination.
	ToPath string `json:"to_path"`
}

func NewRelocationArg() *RelocationArg {
	s := new(RelocationArg)
	return s
}

type RelocationError struct {
	Tag        string       `json:".tag"`
	FromLookup *LookupError `json:"from_lookup,omitempty"`
	FromWrite  *WriteError  `json:"from_write,omitempty"`
	To         *WriteError  `json:"to,omitempty"`
}

func (u *RelocationError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag        string          `json:".tag"`
		FromLookup json.RawMessage `json:"from_lookup"`
		FromWrite  json.RawMessage `json:"from_write"`
		To         json.RawMessage `json:"to"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "from_lookup":
		{
			if len(w.FromLookup) == 0 {
				break
			}
			if err := json.Unmarshal(w.FromLookup, &u.FromLookup); err != nil {
				return err
			}
		}
	case "from_write":
		{
			if len(w.FromWrite) == 0 {
				break
			}
			if err := json.Unmarshal(w.FromWrite, &u.FromWrite); err != nil {
				return err
			}
		}
	case "to":
		{
			if len(w.To) == 0 {
				break
			}
			if err := json.Unmarshal(w.To, &u.To); err != nil {
				return err
			}
		}
	}
	return nil
}

type RestoreArg struct {
	// The path to the file you want to restore.
	Path string `json:"path"`
	// The revision to restore for the file.
	Rev string `json:"rev"`
}

func NewRestoreArg() *RestoreArg {
	s := new(RestoreArg)
	return s
}

type RestoreError struct {
	Tag string `json:".tag"`
	// An error occurs when downloading metadata for the file.
	PathLookup *LookupError `json:"path_lookup,omitempty"`
	// An error occurs when trying to restore the file to that path.
	PathWrite *WriteError `json:"path_write,omitempty"`
}

func (u *RestoreError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// An error occurs when downloading metadata for the file.
		PathLookup json.RawMessage `json:"path_lookup"`
		// An error occurs when trying to restore the file to that path.
		PathWrite json.RawMessage `json:"path_write"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path_lookup":
		{
			if len(w.PathLookup) == 0 {
				break
			}
			if err := json.Unmarshal(w.PathLookup, &u.PathLookup); err != nil {
				return err
			}
		}
	case "path_write":
		{
			if len(w.PathWrite) == 0 {
				break
			}
			if err := json.Unmarshal(w.PathWrite, &u.PathWrite); err != nil {
				return err
			}
		}
	}
	return nil
}

type SearchArg struct {
	// The path in the user's Dropbox to search. Should probably be a folder.
	Path string `json:"path"`
	// The string to search for. The search string is split on spaces into multiple
	// tokens. For file name searching, the last token is used for prefix matching
	// (i.e. "bat c" matches "bat cave" but not "batman car").
	Query string `json:"query"`
	// The starting index within the search results (used for paging).
	Start uint64 `json:"start"`
	// The maximum number of search results to return.
	MaxResults uint64 `json:"max_results"`
	// The search mode (filename, filename_and_content, or deleted_filename). Note
	// that searching file content is only available for Dropbox Business accounts.
	Mode *SearchMode `json:"mode"`
}

func NewSearchArg() *SearchArg {
	s := new(SearchArg)
	s.Start = 0
	s.MaxResults = 100
	s.Mode = &SearchMode{Tag: "filename"}
	return s
}

type SearchError struct {
	Tag  string       `json:".tag"`
	Path *LookupError `json:"path,omitempty"`
}

func (u *SearchError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type SearchMatch struct {
	// The type of the match.
	MatchType *SearchMatchType `json:"match_type"`
	// The metadata for the matched file or folder.
	Metadata *Metadata `json:"metadata"`
}

func NewSearchMatch() *SearchMatch {
	s := new(SearchMatch)
	return s
}

// Indicates what type of match was found for a given item.
type SearchMatchType struct {
	Tag string `json:".tag"`
}

type SearchMode struct {
	Tag string `json:".tag"`
}

type SearchResult struct {
	// A list (possibly empty) of matches for the query.
	Matches []*SearchMatch `json:"matches"`
	// Used for paging. If true, indicates there is another page of results
	// available that can be fetched by calling `Search` again.
	More bool `json:"more"`
	// Used for paging. Value to set the start argument to when calling `Search` to
	// fetch the next page of results.
	Start uint64 `json:"start"`
}

func NewSearchResult() *SearchResult {
	s := new(SearchResult)
	return s
}

type ThumbnailArg struct {
	// The path to the image file you want to thumbnail.
	Path string `json:"path"`
	// The format for the thumbnail image, jpeg (default) or png. For  images that
	// are photos, jpeg should be preferred, while png is  better for screenshots
	// and digital arts.
	Format *ThumbnailFormat `json:"format"`
	// The size for the thumbnail image.
	Size *ThumbnailSize `json:"size"`
}

func NewThumbnailArg() *ThumbnailArg {
	s := new(ThumbnailArg)
	s.Format = &ThumbnailFormat{Tag: "jpeg"}
	s.Size = &ThumbnailSize{Tag: "w64h64"}
	return s
}

type ThumbnailError struct {
	Tag string `json:".tag"`
	// An error occurs when downloading metadata for the image.
	Path *LookupError `json:"path,omitempty"`
}

func (u *ThumbnailError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// An error occurs when downloading metadata for the image.
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type ThumbnailFormat struct {
	Tag string `json:".tag"`
}

type ThumbnailSize struct {
	Tag string `json:".tag"`
}

type UploadError struct {
	Tag string `json:".tag"`
	// Unable to save the uploaded contents to a file.
	Path *UploadWriteFailed `json:"path,omitempty"`
}

func (u *UploadError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// Unable to save the uploaded contents to a file.
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "path":
		{
			if err := json.Unmarshal(body, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type UploadSessionCursor struct {
	// The upload session ID (returned by `UploadSessionStart`).
	SessionId string `json:"session_id"`
	// The amount of data that has been uploaded so far. We use this to make sure
	// upload data isn't lost or duplicated in the event of a network error.
	Offset uint64 `json:"offset"`
}

func NewUploadSessionCursor() *UploadSessionCursor {
	s := new(UploadSessionCursor)
	return s
}

type UploadSessionFinishArg struct {
	// Contains the upload session ID and the offset.
	Cursor *UploadSessionCursor `json:"cursor"`
	// Contains the path and other optional modifiers for the commit.
	Commit *CommitInfo `json:"commit"`
}

func NewUploadSessionFinishArg() *UploadSessionFinishArg {
	s := new(UploadSessionFinishArg)
	return s
}

type UploadSessionFinishError struct {
	Tag string `json:".tag"`
	// The session arguments are incorrect; the value explains the reason.
	LookupFailed *UploadSessionLookupError `json:"lookup_failed,omitempty"`
	// Unable to save the uploaded contents to a file.
	Path *WriteError `json:"path,omitempty"`
}

func (u *UploadSessionFinishError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The session arguments are incorrect; the value explains the reason.
		LookupFailed json.RawMessage `json:"lookup_failed"`
		// Unable to save the uploaded contents to a file.
		Path json.RawMessage `json:"path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "lookup_failed":
		{
			if len(w.LookupFailed) == 0 {
				break
			}
			if err := json.Unmarshal(w.LookupFailed, &u.LookupFailed); err != nil {
				return err
			}
		}
	case "path":
		{
			if len(w.Path) == 0 {
				break
			}
			if err := json.Unmarshal(w.Path, &u.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

type UploadSessionLookupError struct {
	Tag string `json:".tag"`
	// The specified offset was incorrect. See the value for the correct offset.
	// (This error may occur when a previous request was received and processed
	// successfully but the client did not receive the response, e.g. due to a
	// network error.)
	IncorrectOffset *UploadSessionOffsetError `json:"incorrect_offset,omitempty"`
}

func (u *UploadSessionLookupError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The specified offset was incorrect. See the value for the correct offset.
		// (This error may occur when a previous request was received and processed
		// successfully but the client did not receive the response, e.g. due to a
		// network error.)
		IncorrectOffset json.RawMessage `json:"incorrect_offset"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "incorrect_offset":
		{
			if err := json.Unmarshal(body, &u.IncorrectOffset); err != nil {
				return err
			}
		}
	}
	return nil
}

type UploadSessionOffsetError struct {
	// The offset up to which data has been collected.
	CorrectOffset uint64 `json:"correct_offset"`
}

func NewUploadSessionOffsetError() *UploadSessionOffsetError {
	s := new(UploadSessionOffsetError)
	return s
}

type UploadSessionStartResult struct {
	// A unique identifier for the upload session. Pass this to
	// `UploadSessionAppend` and `UploadSessionFinish`.
	SessionId string `json:"session_id"`
}

func NewUploadSessionStartResult() *UploadSessionStartResult {
	s := new(UploadSessionStartResult)
	return s
}

type UploadWriteFailed struct {
	// The reason why the file couldn't be saved.
	Reason *WriteError `json:"reason"`
	// The upload session ID; this may be used to retry the commit.
	UploadSessionId string `json:"upload_session_id"`
}

func NewUploadWriteFailed() *UploadWriteFailed {
	s := new(UploadWriteFailed)
	return s
}

// Metadata for a video.
type VideoMetadata struct {
	// Dimension of the photo/video.
	Dimensions *Dimensions `json:"dimensions,omitempty"`
	// The GPS coordinate of the photo/video.
	Location *GpsCoordinates `json:"location,omitempty"`
	// The timestamp when the photo/video is taken.
	TimeTaken time.Time `json:"time_taken,omitempty"`
	// The duration of the video in milliseconds.
	Duration uint64 `json:"duration,omitempty"`
}

func NewVideoMetadata() *VideoMetadata {
	s := new(VideoMetadata)
	return s
}

type WriteConflictError struct {
	Tag string `json:".tag"`
}

type WriteError struct {
	Tag           string `json:".tag"`
	MalformedPath string `json:"malformed_path,omitempty"`
	// Couldn't write to the target path because there was something in the way.
	Conflict *WriteConflictError `json:"conflict,omitempty"`
}

func (u *WriteError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag           string          `json:".tag"`
		MalformedPath json.RawMessage `json:"malformed_path,omitempty"`
		// Couldn't write to the target path because there was something in the way.
		Conflict json.RawMessage `json:"conflict"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "malformed_path":
		{
			if len(w.MalformedPath) == 0 {
				break
			}
			if err := json.Unmarshal(w.MalformedPath, &u.MalformedPath); err != nil {
				return err
			}
		}
	case "conflict":
		{
			if len(w.Conflict) == 0 {
				break
			}
			if err := json.Unmarshal(w.Conflict, &u.Conflict); err != nil {
				return err
			}
		}
	}
	return nil
}

// Your intent when writing a file to some path. This is used to determine what
// constitutes a conflict and what the autorename strategy is. In some
// situations, the conflict behavior is identical: (a) If the target path
// doesn't contain anything, the file is always written; no conflict. (b) If the
// target path contains a folder, it's always a conflict. (c) If the target path
// contains a file with identical contents, nothing gets written; no conflict.
// The conflict checking differs in the case where there's a file at the target
// path with contents different from the contents you're trying to write.
type WriteMode struct {
	Tag string `json:".tag"`
	// Overwrite if the given "rev" matches the existing file's "rev". The
	// autorename strategy is to append the string "conflicted copy" to the file
	// name. For example, "document.txt" might become "document (conflicted
	// copy).txt" or "document (Panda's conflicted copy).txt".
	Update string `json:"update,omitempty"`
}

func (u *WriteMode) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// Overwrite if the given "rev" matches the existing file's "rev". The
		// autorename strategy is to append the string "conflicted copy" to the file
		// name. For example, "document.txt" might become "document (conflicted
		// copy).txt" or "document (Panda's conflicted copy).txt".
		Update json.RawMessage `json:"update"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "update":
		{
			if len(w.Update) == 0 {
				break
			}
			if err := json.Unmarshal(w.Update, &u.Update); err != nil {
				return err
			}
		}
	}
	return nil
}

type Files interface {
	// Copy a file or folder to a different location in the user's Dropbox. If the
	// source path is a folder all its contents will be copied.
	Copy(arg *RelocationArg) (res *Metadata, err error)
	// Create a folder at a given path.
	CreateFolder(arg *CreateFolderArg) (res *FolderMetadata, err error)
	// Delete the file or folder at a given path. If the path is a folder, all its
	// contents will be deleted too.
	Delete(arg *DeleteArg) (res *Metadata, err error)
	// Download a file from a user's Dropbox.
	Download(arg *DownloadArg) (res *FileMetadata, content io.ReadCloser, err error)
	// Returns the metadata for a file or folder.
	GetMetadata(arg *GetMetadataArg) (res *Metadata, err error)
	// Get a preview for a file. Currently previews are only generated for the
	// files with  the following extensions: .doc, .docx, .docm, .ppt, .pps, .ppsx,
	// .ppsm, .pptx, .pptm,  .xls, .xlsx, .xlsm, .rtf
	GetPreview(arg *PreviewArg) (res *FileMetadata, content io.ReadCloser, err error)
	// Get a thumbnail for an image. This method currently supports files with the
	// following file extensions: jpg, jpeg, png, tiff, tif, gif and bmp. Photos
	// that are larger than 20MB in size won't be converted to a thumbnail.
	GetThumbnail(arg *ThumbnailArg) (res *FileMetadata, content io.ReadCloser, err error)
	// Returns the contents of a folder.
	ListFolder(arg *ListFolderArg) (res *ListFolderResult, err error)
	// Once a cursor has been retrieved from `ListFolder`, use this to paginate
	// through all files and retrieve updates to the folder.
	ListFolderContinue(arg *ListFolderContinueArg) (res *ListFolderResult, err error)
	// A way to quickly get a cursor for the folder's state. Unlike `ListFolder`,
	// `ListFolderGetLatestCursor` doesn't return any entries. This endpoint is for
	// app which only needs to know about new files and modifications and doesn't
	// need to know about files that already exist in Dropbox.
	ListFolderGetLatestCursor(arg *ListFolderArg) (res *ListFolderGetLatestCursorResult, err error)
	// A longpoll endpoint to wait for changes on an account. In conjunction with
	// `ListFolder`, this call gives you a low-latency way to monitor an account
	// for file changes. The connection will block until there are changes
	// available or a timeout occurs. This endpoint is useful mostly for
	// client-side apps. If you're looking for server-side notifications, check out
	// our `webhooks documentation`
	// <https://www.dropbox.com/developers/reference/webhooks>.
	ListFolderLongpoll(arg *ListFolderLongpollArg) (res *ListFolderLongpollResult, err error)
	// Return revisions of a file
	ListRevisions(arg *ListRevisionsArg) (res *ListRevisionsResult, err error)
	// Move a file or folder to a different location in the user's Dropbox. If the
	// source path is a folder all its contents will be moved.
	Move(arg *RelocationArg) (res *Metadata, err error)
	// Permanently delete the file or folder at a given path (see
	// https://www.dropbox.com/en/help/40). Note: This endpoint is only available
	// for Dropbox Business apps.
	PermanentlyDelete(arg *DeleteArg) (err error)
	// Restore a file to a specific revision
	Restore(arg *RestoreArg) (res *FileMetadata, err error)
	// Searches for files and folders.
	Search(arg *SearchArg) (res *SearchResult, err error)
	// Create a new file with the contents provided in the request. Do not use this
	// to upload a file larger than 150 MB. Instead, create an upload session with
	// `UploadSessionStart`.
	Upload(arg *CommitInfo, content io.Reader) (res *FileMetadata, err error)
	// Append more data to an upload session. A single request should not upload
	// more than 150 MB of file contents.
	UploadSessionAppend(arg *UploadSessionCursor, content io.Reader) (err error)
	// Finish an upload session and save the uploaded data to the given file path.
	// A single request should not upload more than 150 MB of file contents.
	UploadSessionFinish(arg *UploadSessionFinishArg, content io.Reader) (res *FileMetadata, err error)
	// Upload sessions allow you to upload a single file using multiple requests.
	// This call starts a new upload session with the given data.  You can then use
	// `UploadSessionAppend` to add more data and `UploadSessionFinish` to save all
	// the data to a file in Dropbox. A single request should not upload more than
	// 150 MB of file contents.
	UploadSessionStart(content io.Reader) (res *UploadSessionStartResult, err error)
}
