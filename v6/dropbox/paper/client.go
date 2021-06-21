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

package paper

import (
	"encoding/json"
	"io"
	"log"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// DocsArchive : Marks the given Paper doc as archived. This action can be
	// performed or undone by anyone with edit permissions to the doc. Note that
	// this endpoint will continue to work for content created by users on the
	// older version of Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. This endpoint will be
	// retired in September 2020. Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for more information.
	// Deprecated:
	DocsArchive(arg *RefPaperDoc) (err error)
	// DocsCreate : Creates a new Paper doc with the provided content. Note that
	// this endpoint will continue to work for content created by users on the
	// older version of Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. This endpoint will be
	// retired in September 2020. Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for more information.
	// Deprecated:
	DocsCreate(arg *PaperDocCreateArgs, content io.Reader) (res *PaperDocCreateUpdateResult, err error)
	// DocsDownload : Exports and downloads Paper doc either as HTML or
	// markdown. Note that this endpoint will continue to work for content
	// created by users on the older version of Paper. To check which version of
	// Paper a user is on, use /users/features/get_values. If the paper_as_files
	// feature is enabled, then the user is running the new version of Paper.
	// Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsDownload(arg *PaperDocExport) (res *PaperDocExportResult, content io.ReadCloser, err error)
	// DocsFolderUsersList : Lists the users who are explicitly invited to the
	// Paper folder in which the Paper doc is contained. For private folders all
	// users (including owner) shared on the folder are listed and for team
	// folders all non-team users shared on the folder are returned. Note that
	// this endpoint will continue to work for content created by users on the
	// older version of Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. Refer to the `Paper
	// Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsFolderUsersList(arg *ListUsersOnFolderArgs) (res *ListUsersOnFolderResponse, err error)
	// DocsFolderUsersListContinue : Once a cursor has been retrieved from
	// `docsFolderUsersList`, use this to paginate through all users on the
	// Paper folder. Note that this endpoint will continue to work for content
	// created by users on the older version of Paper. To check which version of
	// Paper a user is on, use /users/features/get_values. If the paper_as_files
	// feature is enabled, then the user is running the new version of Paper.
	// Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsFolderUsersListContinue(arg *ListUsersOnFolderContinueArgs) (res *ListUsersOnFolderResponse, err error)
	// DocsGetFolderInfo : Retrieves folder information for the given Paper doc.
	// This includes:   - folder sharing policy; permissions for subfolders are
	// set by the top-level folder.   - full 'filepath', i.e. the list of
	// folders (both folderId and folderName) from     the root folder to the
	// folder directly containing the Paper doc.  If the Paper doc is not in any
	// folder (aka unfiled) the response will be empty. Note that this endpoint
	// will continue to work for content created by users on the older version
	// of Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. Refer to the `Paper
	// Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsGetFolderInfo(arg *RefPaperDoc) (res *FoldersContainingPaperDoc, err error)
	// DocsList : Return the list of all Paper docs according to the argument
	// specifications. To iterate over through the full pagination, pass the
	// cursor to `docsListContinue`. Note that this endpoint will continue to
	// work for content created by users on the older version of Paper. To check
	// which version of Paper a user is on, use /users/features/get_values. If
	// the paper_as_files feature is enabled, then the user is running the new
	// version of Paper. Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsList(arg *ListPaperDocsArgs) (res *ListPaperDocsResponse, err error)
	// DocsListContinue : Once a cursor has been retrieved from `docsList`, use
	// this to paginate through all Paper doc. Note that this endpoint will
	// continue to work for content created by users on the older version of
	// Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. Refer to the `Paper
	// Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsListContinue(arg *ListPaperDocsContinueArgs) (res *ListPaperDocsResponse, err error)
	// DocsPermanentlyDelete : Permanently deletes the given Paper doc. This
	// operation is final as the doc cannot be recovered. This action can be
	// performed only by the doc owner. Note that this endpoint will continue to
	// work for content created by users on the older version of Paper. To check
	// which version of Paper a user is on, use /users/features/get_values. If
	// the paper_as_files feature is enabled, then the user is running the new
	// version of Paper. Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsPermanentlyDelete(arg *RefPaperDoc) (err error)
	// DocsSharingPolicyGet : Gets the default sharing policy for the given
	// Paper doc. Note that this endpoint will continue to work for content
	// created by users on the older version of Paper. To check which version of
	// Paper a user is on, use /users/features/get_values. If the paper_as_files
	// feature is enabled, then the user is running the new version of Paper.
	// Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsSharingPolicyGet(arg *RefPaperDoc) (res *SharingPolicy, err error)
	// DocsSharingPolicySet : Sets the default sharing policy for the given
	// Paper doc. The default 'team_sharing_policy' can be changed only by
	// teams, omit this field for personal accounts. The 'public_sharing_policy'
	// policy can't be set to the value 'disabled' because this setting can be
	// changed only via the team admin console. Note that this endpoint will
	// continue to work for content created by users on the older version of
	// Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. Refer to the `Paper
	// Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsSharingPolicySet(arg *PaperDocSharingPolicy) (err error)
	// DocsUpdate : Updates an existing Paper doc with the provided content.
	// Note that this endpoint will continue to work for content created by
	// users on the older version of Paper. To check which version of Paper a
	// user is on, use /users/features/get_values. If the paper_as_files feature
	// is enabled, then the user is running the new version of Paper. This
	// endpoint will be retired in September 2020. Refer to the `Paper Migration
	// Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for more information.
	// Deprecated:
	DocsUpdate(arg *PaperDocUpdateArgs, content io.Reader) (res *PaperDocCreateUpdateResult, err error)
	// DocsUsersAdd : Allows an owner or editor to add users to a Paper doc or
	// change their permissions using their email address or Dropbox account ID.
	// The doc owner's permissions cannot be changed. Note that this endpoint
	// will continue to work for content created by users on the older version
	// of Paper. To check which version of Paper a user is on, use
	// /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. Refer to the `Paper
	// Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsUsersAdd(arg *AddPaperDocUser) (res []*AddPaperDocUserMemberResult, err error)
	// DocsUsersList : Lists all users who visited the Paper doc or users with
	// explicit access. This call excludes users who have been removed. The list
	// is sorted by the date of the visit or the share date. The list will
	// include both users, the explicitly shared ones as well as those who came
	// in using the Paper url link. Note that this endpoint will continue to
	// work for content created by users on the older version of Paper. To check
	// which version of Paper a user is on, use /users/features/get_values. If
	// the paper_as_files feature is enabled, then the user is running the new
	// version of Paper. Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsUsersList(arg *ListUsersOnPaperDocArgs) (res *ListUsersOnPaperDocResponse, err error)
	// DocsUsersListContinue : Once a cursor has been retrieved from
	// `docsUsersList`, use this to paginate through all users on the Paper doc.
	// Note that this endpoint will continue to work for content created by
	// users on the older version of Paper. To check which version of Paper a
	// user is on, use /users/features/get_values. If the paper_as_files feature
	// is enabled, then the user is running the new version of Paper. Refer to
	// the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsUsersListContinue(arg *ListUsersOnPaperDocContinueArgs) (res *ListUsersOnPaperDocResponse, err error)
	// DocsUsersRemove : Allows an owner or editor to remove users from a Paper
	// doc using their email address or Dropbox account ID. The doc owner cannot
	// be removed. Note that this endpoint will continue to work for content
	// created by users on the older version of Paper. To check which version of
	// Paper a user is on, use /users/features/get_values. If the paper_as_files
	// feature is enabled, then the user is running the new version of Paper.
	// Refer to the `Paper Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	DocsUsersRemove(arg *RemovePaperDocUser) (err error)
	// FoldersCreate : Create a new Paper folder with the provided info. Note
	// that this endpoint will continue to work for content created by users on
	// the older version of Paper. To check which version of Paper a user is on,
	// use /users/features/get_values. If the paper_as_files feature is enabled,
	// then the user is running the new version of Paper. Refer to the `Paper
	// Migration Guide`
	// <https://www.dropbox.com/lp/developers/reference/paper-migration-guide>
	// for migration information.
	// Deprecated:
	FoldersCreate(arg *PaperFolderCreateArg) (res *PaperFolderCreateResult, err error)
}

type apiImpl dropbox.Context

//DocsArchiveAPIError is an error-wrapper for the docs/archive route
type DocsArchiveAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsArchive(arg *RefPaperDoc) (err error) {
	log.Printf("WARNING: API `DocsArchive` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/archive",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsArchiveAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

//DocsCreateAPIError is an error-wrapper for the docs/create route
type DocsCreateAPIError struct {
	dropbox.APIError
	EndpointError *PaperDocCreateError `json:"error"`
}

func (dbx *apiImpl) DocsCreate(arg *PaperDocCreateArgs, content io.Reader) (res *PaperDocCreateUpdateResult, err error) {
	log.Printf("WARNING: API `DocsCreate` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/create",
		Auth:         "user",
		Style:        "upload",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, content)
	if err != nil {
		var appErr DocsCreateAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsDownloadAPIError is an error-wrapper for the docs/download route
type DocsDownloadAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsDownload(arg *PaperDocExport) (res *PaperDocExportResult, content io.ReadCloser, err error) {
	log.Printf("WARNING: API `DocsDownload` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/download",
		Auth:         "user",
		Style:        "download",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsDownloadAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return
	}

	content = respBody
	return
}

//DocsFolderUsersListAPIError is an error-wrapper for the docs/folder_users/list route
type DocsFolderUsersListAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsFolderUsersList(arg *ListUsersOnFolderArgs) (res *ListUsersOnFolderResponse, err error) {
	log.Printf("WARNING: API `DocsFolderUsersList` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/folder_users/list",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsFolderUsersListAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsFolderUsersListContinueAPIError is an error-wrapper for the docs/folder_users/list/continue route
type DocsFolderUsersListContinueAPIError struct {
	dropbox.APIError
	EndpointError *ListUsersCursorError `json:"error"`
}

func (dbx *apiImpl) DocsFolderUsersListContinue(arg *ListUsersOnFolderContinueArgs) (res *ListUsersOnFolderResponse, err error) {
	log.Printf("WARNING: API `DocsFolderUsersListContinue` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/folder_users/list/continue",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsFolderUsersListContinueAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsGetFolderInfoAPIError is an error-wrapper for the docs/get_folder_info route
type DocsGetFolderInfoAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsGetFolderInfo(arg *RefPaperDoc) (res *FoldersContainingPaperDoc, err error) {
	log.Printf("WARNING: API `DocsGetFolderInfo` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/get_folder_info",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsGetFolderInfoAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsListAPIError is an error-wrapper for the docs/list route
type DocsListAPIError struct {
	dropbox.APIError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) DocsList(arg *ListPaperDocsArgs) (res *ListPaperDocsResponse, err error) {
	log.Printf("WARNING: API `DocsList` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/list",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsListAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsListContinueAPIError is an error-wrapper for the docs/list/continue route
type DocsListContinueAPIError struct {
	dropbox.APIError
	EndpointError *ListDocsCursorError `json:"error"`
}

func (dbx *apiImpl) DocsListContinue(arg *ListPaperDocsContinueArgs) (res *ListPaperDocsResponse, err error) {
	log.Printf("WARNING: API `DocsListContinue` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/list/continue",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsListContinueAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsPermanentlyDeleteAPIError is an error-wrapper for the docs/permanently_delete route
type DocsPermanentlyDeleteAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsPermanentlyDelete(arg *RefPaperDoc) (err error) {
	log.Printf("WARNING: API `DocsPermanentlyDelete` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/permanently_delete",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsPermanentlyDeleteAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

//DocsSharingPolicyGetAPIError is an error-wrapper for the docs/sharing_policy/get route
type DocsSharingPolicyGetAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsSharingPolicyGet(arg *RefPaperDoc) (res *SharingPolicy, err error) {
	log.Printf("WARNING: API `DocsSharingPolicyGet` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/sharing_policy/get",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsSharingPolicyGetAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsSharingPolicySetAPIError is an error-wrapper for the docs/sharing_policy/set route
type DocsSharingPolicySetAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsSharingPolicySet(arg *PaperDocSharingPolicy) (err error) {
	log.Printf("WARNING: API `DocsSharingPolicySet` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/sharing_policy/set",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsSharingPolicySetAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

//DocsUpdateAPIError is an error-wrapper for the docs/update route
type DocsUpdateAPIError struct {
	dropbox.APIError
	EndpointError *PaperDocUpdateError `json:"error"`
}

func (dbx *apiImpl) DocsUpdate(arg *PaperDocUpdateArgs, content io.Reader) (res *PaperDocCreateUpdateResult, err error) {
	log.Printf("WARNING: API `DocsUpdate` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/update",
		Auth:         "user",
		Style:        "upload",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, content)
	if err != nil {
		var appErr DocsUpdateAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsUsersAddAPIError is an error-wrapper for the docs/users/add route
type DocsUsersAddAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsUsersAdd(arg *AddPaperDocUser) (res []*AddPaperDocUserMemberResult, err error) {
	log.Printf("WARNING: API `DocsUsersAdd` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/users/add",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsUsersAddAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsUsersListAPIError is an error-wrapper for the docs/users/list route
type DocsUsersListAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsUsersList(arg *ListUsersOnPaperDocArgs) (res *ListUsersOnPaperDocResponse, err error) {
	log.Printf("WARNING: API `DocsUsersList` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/users/list",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsUsersListAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsUsersListContinueAPIError is an error-wrapper for the docs/users/list/continue route
type DocsUsersListContinueAPIError struct {
	dropbox.APIError
	EndpointError *ListUsersCursorError `json:"error"`
}

func (dbx *apiImpl) DocsUsersListContinue(arg *ListUsersOnPaperDocContinueArgs) (res *ListUsersOnPaperDocResponse, err error) {
	log.Printf("WARNING: API `DocsUsersListContinue` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/users/list/continue",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsUsersListContinueAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

//DocsUsersRemoveAPIError is an error-wrapper for the docs/users/remove route
type DocsUsersRemoveAPIError struct {
	dropbox.APIError
	EndpointError *DocLookupError `json:"error"`
}

func (dbx *apiImpl) DocsUsersRemove(arg *RemovePaperDocUser) (err error) {
	log.Printf("WARNING: API `DocsUsersRemove` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "docs/users/remove",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr DocsUsersRemoveAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return
	}

	_ = resp
	_ = respBody
	return
}

//FoldersCreateAPIError is an error-wrapper for the folders/create route
type FoldersCreateAPIError struct {
	dropbox.APIError
	EndpointError *PaperFolderCreateError `json:"error"`
}

func (dbx *apiImpl) FoldersCreate(arg *PaperFolderCreateArg) (res *PaperFolderCreateResult, err error) {
	log.Printf("WARNING: API `FoldersCreate` is deprecated")

	req := dropbox.Request{
		Host:         "api",
		Namespace:    "paper",
		Route:        "folders/create",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, nil)
	if err != nil {
		var appErr FoldersCreateAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
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

// New returns a Client implementation for this namespace
func New(c dropbox.Config) Client {
	ctx := apiImpl(dropbox.NewContext(c))
	return &ctx
}
