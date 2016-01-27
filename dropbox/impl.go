/* DO NOT EDIT */

package dropbox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dropbox/dropbox-sdk-go/dropbox/apierror"
	"github.com/dropbox/dropbox-sdk-go/dropbox/async"
	"github.com/dropbox/dropbox-sdk-go/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go/dropbox/sharing"
	"github.com/dropbox/dropbox-sdk-go/dropbox/team"
	"github.com/dropbox/dropbox-sdk-go/dropbox/users"
)

type Api interface {
	files.Files
	sharing.Sharing
	team.Team
	users.Users
}

type GetMetadataWrapper struct {
	apierror.ApiError
	EndpointError *files.GetMetadataError `json:"error"`
}

func (dbx *apiImpl) GetMetadata(arg *files.GetMetadataArg) (res *files.Metadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/get_metadata", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetMetadataWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFolderLongpollWrapper struct {
	apierror.ApiError
	EndpointError *files.ListFolderLongpollError `json:"error"`
}

func (dbx *apiImpl) ListFolderLongpoll(arg *files.ListFolderLongpollArg) (res *files.ListFolderLongpollResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://notify.dropboxapi.com/2/files/list_folder/longpoll", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFolderLongpollWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFolderWrapper struct {
	apierror.ApiError
	EndpointError *files.ListFolderError `json:"error"`
}

func (dbx *apiImpl) ListFolder(arg *files.ListFolderArg) (res *files.ListFolderResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/list_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFolderContinueWrapper struct {
	apierror.ApiError
	EndpointError *files.ListFolderContinueError `json:"error"`
}

func (dbx *apiImpl) ListFolderContinue(arg *files.ListFolderContinueArg) (res *files.ListFolderResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/list_folder/continue", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFolderContinueWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFolderGetLatestCursorWrapper struct {
	apierror.ApiError
	EndpointError *files.ListFolderError `json:"error"`
}

func (dbx *apiImpl) ListFolderGetLatestCursor(arg *files.ListFolderArg) (res *files.ListFolderGetLatestCursorResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/list_folder/get_latest_cursor", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFolderGetLatestCursorWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type DownloadWrapper struct {
	apierror.ApiError
	EndpointError *files.DownloadError `json:"error"`
}

func (dbx *apiImpl) Download(arg *files.DownloadArg) (res *files.FileMetadata, content io.ReadCloser, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/download", nil)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	body := []byte(resp.Header.Get("Dropbox-API-Result"))
	content = resp.Body
	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap DownloadWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UploadSessionStartWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) UploadSessionStart(content io.Reader) (res *files.UploadSessionStartResult, err error) {
	cli := dbx.client

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload_session/start", nil)
	if err != nil {
		return
	}

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UploadSessionStartWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UploadSessionAppendWrapper struct {
	apierror.ApiError
	EndpointError *files.UploadSessionLookupError `json:"error"`
}

func (dbx *apiImpl) UploadSessionAppend(arg *files.UploadSessionCursor, content io.Reader) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload_session/append", content)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UploadSessionAppendWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UploadSessionFinishWrapper struct {
	apierror.ApiError
	EndpointError *files.UploadSessionFinishError `json:"error"`
}

func (dbx *apiImpl) UploadSessionFinish(arg *files.UploadSessionFinishArg, content io.Reader) (res *files.FileMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload_session/finish", content)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UploadSessionFinishWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UploadWrapper struct {
	apierror.ApiError
	EndpointError *files.UploadError `json:"error"`
}

func (dbx *apiImpl) Upload(arg *files.CommitInfo, content io.Reader) (res *files.FileMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", content)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UploadWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type SearchWrapper struct {
	apierror.ApiError
	EndpointError *files.SearchError `json:"error"`
}

func (dbx *apiImpl) Search(arg *files.SearchArg) (res *files.SearchResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/search", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap SearchWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type CreateFolderWrapper struct {
	apierror.ApiError
	EndpointError *files.CreateFolderError `json:"error"`
}

func (dbx *apiImpl) CreateFolder(arg *files.CreateFolderArg) (res *files.FolderMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/create_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap CreateFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type DeleteWrapper struct {
	apierror.ApiError
	EndpointError *files.DeleteError `json:"error"`
}

func (dbx *apiImpl) Delete(arg *files.DeleteArg) (res *files.Metadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/delete", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap DeleteWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type PermanentlyDeleteWrapper struct {
	apierror.ApiError
	EndpointError *files.DeleteError `json:"error"`
}

func (dbx *apiImpl) PermanentlyDelete(arg *files.DeleteArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/permanently_delete", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap PermanentlyDeleteWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type CopyWrapper struct {
	apierror.ApiError
	EndpointError *files.RelocationError `json:"error"`
}

func (dbx *apiImpl) Copy(arg *files.RelocationArg) (res *files.Metadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/copy", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap CopyWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MoveWrapper struct {
	apierror.ApiError
	EndpointError *files.RelocationError `json:"error"`
}

func (dbx *apiImpl) Move(arg *files.RelocationArg) (res *files.Metadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/move", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MoveWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetThumbnailWrapper struct {
	apierror.ApiError
	EndpointError *files.ThumbnailError `json:"error"`
}

func (dbx *apiImpl) GetThumbnail(arg *files.ThumbnailArg) (res *files.FileMetadata, content io.ReadCloser, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/get_thumbnail", nil)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	body := []byte(resp.Header.Get("Dropbox-API-Result"))
	content = resp.Body
	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetThumbnailWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetPreviewWrapper struct {
	apierror.ApiError
	EndpointError *files.PreviewError `json:"error"`
}

func (dbx *apiImpl) GetPreview(arg *files.PreviewArg) (res *files.FileMetadata, content io.ReadCloser, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/get_preview", nil)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	body := []byte(resp.Header.Get("Dropbox-API-Result"))
	content = resp.Body
	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetPreviewWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListRevisionsWrapper struct {
	apierror.ApiError
	EndpointError *files.ListRevisionsError `json:"error"`
}

func (dbx *apiImpl) ListRevisions(arg *files.ListRevisionsArg) (res *files.ListRevisionsResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/list_revisions", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListRevisionsWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type RestoreWrapper struct {
	apierror.ApiError
	EndpointError *files.RestoreError `json:"error"`
}

func (dbx *apiImpl) Restore(arg *files.RestoreArg) (res *files.FileMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/restore", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap RestoreWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetSharedLinkMetadataWrapper struct {
	apierror.ApiError
	EndpointError *sharing.SharedLinkError `json:"error"`
}

func (dbx *apiImpl) GetSharedLinkMetadata(arg *sharing.GetSharedLinkMetadataArg) (res *sharing.SharedLinkMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/get_shared_link_metadata", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetSharedLinkMetadataWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListSharedLinksWrapper struct {
	apierror.ApiError
	EndpointError *sharing.ListSharedLinksError `json:"error"`
}

func (dbx *apiImpl) ListSharedLinks(arg *sharing.ListSharedLinksArg) (res *sharing.ListSharedLinksResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/list_shared_links", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListSharedLinksWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ModifySharedLinkSettingsWrapper struct {
	apierror.ApiError
	EndpointError *sharing.ModifySharedLinkSettingsError `json:"error"`
}

func (dbx *apiImpl) ModifySharedLinkSettings(arg *sharing.ModifySharedLinkSettingsArgs) (res *sharing.SharedLinkMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/modify_shared_link_settings", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ModifySharedLinkSettingsWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type CreateSharedLinkWithSettingsWrapper struct {
	apierror.ApiError
	EndpointError *sharing.CreateSharedLinkWithSettingsError `json:"error"`
}

func (dbx *apiImpl) CreateSharedLinkWithSettings(arg *sharing.CreateSharedLinkWithSettingsArg) (res *sharing.SharedLinkMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap CreateSharedLinkWithSettingsWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type RevokeSharedLinkWrapper struct {
	apierror.ApiError
	EndpointError *sharing.RevokeSharedLinkError `json:"error"`
}

func (dbx *apiImpl) RevokeSharedLink(arg *sharing.RevokeSharedLinkArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/revoke_shared_link", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap RevokeSharedLinkWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetSharedLinkFileWrapper struct {
	apierror.ApiError
	EndpointError *sharing.GetSharedLinkFileError `json:"error"`
}

func (dbx *apiImpl) GetSharedLinkFile(arg *sharing.GetSharedLinkMetadataArg) (res *sharing.SharedLinkMetadata, content io.ReadCloser, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/sharing/get_shared_link_file", nil)
	if err != nil {
		return
	}

	req.Header.Set("Dropbox-API-Arg", string(b))
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	body := []byte(resp.Header.Get("Dropbox-API-Result"))
	content = resp.Body
	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetSharedLinkFileWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetSharedLinksWrapper struct {
	apierror.ApiError
	EndpointError *sharing.GetSharedLinksError `json:"error"`
}

func (dbx *apiImpl) GetSharedLinks(arg *sharing.GetSharedLinksArg) (res *sharing.GetSharedLinksResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/get_shared_links", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetSharedLinksWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type CreateSharedLinkWrapper struct {
	apierror.ApiError
	EndpointError *sharing.CreateSharedLinkError `json:"error"`
}

func (dbx *apiImpl) CreateSharedLink(arg *sharing.CreateSharedLinkArg) (res *sharing.PathLinkMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/create_shared_link", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap CreateSharedLinkWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFoldersWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) ListFolders() (res *sharing.ListFoldersResult, err error) {
	cli := dbx.client

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/list_folders", nil)
	if err != nil {
		return
	}

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFoldersWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFoldersContinueWrapper struct {
	apierror.ApiError
	EndpointError *sharing.ListFoldersContinueError `json:"error"`
}

func (dbx *apiImpl) ListFoldersContinue(arg *sharing.ListFoldersContinueArg) (res *sharing.ListFoldersResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/list_folders/continue", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFoldersContinueWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetFolderMetadataWrapper struct {
	apierror.ApiError
	EndpointError *sharing.SharedFolderAccessError `json:"error"`
}

func (dbx *apiImpl) GetFolderMetadata(arg *sharing.GetMetadataArgs) (res *sharing.SharedFolderMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/get_folder_metadata", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetFolderMetadataWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFolderMembersWrapper struct {
	apierror.ApiError
	EndpointError *sharing.SharedFolderAccessError `json:"error"`
}

func (dbx *apiImpl) ListFolderMembers(arg *sharing.ListFolderMembersArgs) (res *sharing.SharedFolderMembers, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/list_folder_members", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFolderMembersWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ListFolderMembersContinueWrapper struct {
	apierror.ApiError
	EndpointError *sharing.ListFolderMembersContinueError `json:"error"`
}

func (dbx *apiImpl) ListFolderMembersContinue(arg *sharing.ListFolderMembersContinueArg) (res *sharing.SharedFolderMembers, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/list_folder_members/continue", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ListFolderMembersContinueWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ShareFolderWrapper struct {
	apierror.ApiError
	EndpointError *sharing.ShareFolderError `json:"error"`
}

func (dbx *apiImpl) ShareFolder(arg *sharing.ShareFolderArg) (res *sharing.ShareFolderLaunch, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/share_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ShareFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type CheckShareJobStatusWrapper struct {
	apierror.ApiError
	EndpointError *async.PollError `json:"error"`
}

func (dbx *apiImpl) CheckShareJobStatus(arg *async.PollArg) (res *sharing.ShareFolderJobStatus, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/check_share_job_status", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap CheckShareJobStatusWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type CheckJobStatusWrapper struct {
	apierror.ApiError
	EndpointError *async.PollError `json:"error"`
}

func (dbx *apiImpl) CheckJobStatus(arg *async.PollArg) (res *sharing.JobStatus, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/check_job_status", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap CheckJobStatusWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UnshareFolderWrapper struct {
	apierror.ApiError
	EndpointError *sharing.UnshareFolderError `json:"error"`
}

func (dbx *apiImpl) UnshareFolder(arg *sharing.UnshareFolderArg) (res *async.LaunchEmptyResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/unshare_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UnshareFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type TransferFolderWrapper struct {
	apierror.ApiError
	EndpointError *sharing.TransferFolderError `json:"error"`
}

func (dbx *apiImpl) TransferFolder(arg *sharing.TransferFolderArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/transfer_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap TransferFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UpdateFolderPolicyWrapper struct {
	apierror.ApiError
	EndpointError *sharing.UpdateFolderPolicyError `json:"error"`
}

func (dbx *apiImpl) UpdateFolderPolicy(arg *sharing.UpdateFolderPolicyArg) (res *sharing.SharedFolderMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/update_folder_policy", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UpdateFolderPolicyWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type AddFolderMemberWrapper struct {
	apierror.ApiError
	EndpointError *sharing.AddFolderMemberError `json:"error"`
}

func (dbx *apiImpl) AddFolderMember(arg *sharing.AddFolderMemberArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/add_folder_member", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap AddFolderMemberWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type RemoveFolderMemberWrapper struct {
	apierror.ApiError
	EndpointError *sharing.RemoveFolderMemberError `json:"error"`
}

func (dbx *apiImpl) RemoveFolderMember(arg *sharing.RemoveFolderMemberArg) (res *async.LaunchEmptyResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/remove_folder_member", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap RemoveFolderMemberWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UpdateFolderMemberWrapper struct {
	apierror.ApiError
	EndpointError *sharing.UpdateFolderMemberError `json:"error"`
}

func (dbx *apiImpl) UpdateFolderMember(arg *sharing.UpdateFolderMemberArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/update_folder_member", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UpdateFolderMemberWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MountFolderWrapper struct {
	apierror.ApiError
	EndpointError *sharing.MountFolderError `json:"error"`
}

func (dbx *apiImpl) MountFolder(arg *sharing.MountFolderArg) (res *sharing.SharedFolderMetadata, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/mount_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MountFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type UnmountFolderWrapper struct {
	apierror.ApiError
	EndpointError *sharing.UnmountFolderError `json:"error"`
}

func (dbx *apiImpl) UnmountFolder(arg *sharing.UnmountFolderArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/unmount_folder", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap UnmountFolderWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type RelinquishFolderMembershipWrapper struct {
	apierror.ApiError
	EndpointError *sharing.RelinquishFolderMembershipError `json:"error"`
}

func (dbx *apiImpl) RelinquishFolderMembership(arg *sharing.RelinquishFolderMembershipArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/sharing/relinquish_folder_membership", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap RelinquishFolderMembershipWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetInfoWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) GetInfo() (res *team.TeamGetInfoResult, err error) {
	cli := dbx.client

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/get_info", nil)
	if err != nil {
		return
	}

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetInfoWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type DevicesListMemberDevicesWrapper struct {
	apierror.ApiError
	EndpointError *team.ListMemberDevicesError `json:"error"`
}

func (dbx *apiImpl) DevicesListMemberDevices(arg *team.ListMemberDevicesArg) (res *team.ListMemberDevicesResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/devices/list_member_devices", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap DevicesListMemberDevicesWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type DevicesListTeamDevicesWrapper struct {
	apierror.ApiError
	EndpointError *team.ListTeamDevicesError `json:"error"`
}

func (dbx *apiImpl) DevicesListTeamDevices(arg *team.ListTeamDevicesArg) (res *team.ListTeamDevicesResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/devices/list_team_devices", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap DevicesListTeamDevicesWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type DevicesRevokeDeviceSessionWrapper struct {
	apierror.ApiError
	EndpointError *team.RevokeDeviceSessionError `json:"error"`
}

func (dbx *apiImpl) DevicesRevokeDeviceSession(arg *team.RevokeDeviceSessionArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/devices/revoke_device_session", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap DevicesRevokeDeviceSessionWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type DevicesRevokeDeviceSessionBatchWrapper struct {
	apierror.ApiError
	EndpointError *team.RevokeDeviceSessionBatchError `json:"error"`
}

func (dbx *apiImpl) DevicesRevokeDeviceSessionBatch(arg *team.RevokeDeviceSessionBatchArg) (res *team.RevokeDeviceSessionBatchResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/devices/revoke_device_session_batch", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap DevicesRevokeDeviceSessionBatchWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsListWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) GroupsList(arg *team.GroupsListArg) (res *team.GroupsListResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/list", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsListWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsListContinueWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupsListContinueError `json:"error"`
}

func (dbx *apiImpl) GroupsListContinue(arg *team.GroupsListContinueArg) (res *team.GroupsListResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/list/continue", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsListContinueWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsGetInfoWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupsGetInfoError `json:"error"`
}

func (dbx *apiImpl) GroupsGetInfo(arg *team.GroupsSelector) (res []*team.GroupsGetInfoItem, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/get_info", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsGetInfoWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsCreateWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupCreateError `json:"error"`
}

func (dbx *apiImpl) GroupsCreate(arg *team.GroupCreateArg) (res *team.GroupFullInfo, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/create", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsCreateWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsDeleteWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupDeleteError `json:"error"`
}

func (dbx *apiImpl) GroupsDelete(arg *team.GroupSelector) (res *async.LaunchEmptyResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/delete", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsDeleteWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsUpdateWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupUpdateError `json:"error"`
}

func (dbx *apiImpl) GroupsUpdate(arg *team.GroupUpdateArgs) (res *team.GroupFullInfo, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/update", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsUpdateWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsMembersAddWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupMembersAddError `json:"error"`
}

func (dbx *apiImpl) GroupsMembersAdd(arg *team.GroupMembersAddArg) (res *team.GroupMembersChangeResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/members/add", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsMembersAddWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsMembersRemoveWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupMembersRemoveError `json:"error"`
}

func (dbx *apiImpl) GroupsMembersRemove(arg *team.GroupMembersRemoveArg) (res *team.GroupMembersChangeResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/members/remove", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsMembersRemoveWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsMembersSetAccessTypeWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupMemberSelectorError `json:"error"`
}

func (dbx *apiImpl) GroupsMembersSetAccessType(arg *team.GroupMembersSetAccessTypeArg) (res []*team.GroupsGetInfoItem, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/members/set_access_type", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsMembersSetAccessTypeWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GroupsJobStatusGetWrapper struct {
	apierror.ApiError
	EndpointError *team.GroupsPollError `json:"error"`
}

func (dbx *apiImpl) GroupsJobStatusGet(arg *async.PollArg) (res *async.PollEmptyResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/groups/job_status/get", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GroupsJobStatusGetWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type LinkedAppsListMemberLinkedAppsWrapper struct {
	apierror.ApiError
	EndpointError *team.ListMemberAppsError `json:"error"`
}

func (dbx *apiImpl) LinkedAppsListMemberLinkedApps(arg *team.ListMemberAppsArg) (res *team.ListMemberAppsResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/linked_apps/list_member_linked_apps", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap LinkedAppsListMemberLinkedAppsWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type LinkedAppsListTeamLinkedAppsWrapper struct {
	apierror.ApiError
	EndpointError *team.ListTeamAppsError `json:"error"`
}

func (dbx *apiImpl) LinkedAppsListTeamLinkedApps(arg *team.ListTeamAppsArg) (res *team.ListTeamAppsResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/linked_apps/list_team_linked_apps", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap LinkedAppsListTeamLinkedAppsWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type LinkedAppsRevokeLinkedAppWrapper struct {
	apierror.ApiError
	EndpointError *team.RevokeLinkedAppError `json:"error"`
}

func (dbx *apiImpl) LinkedAppsRevokeLinkedApp(arg *team.RevokeLinkedApiAppArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/linked_apps/revoke_linked_app", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap LinkedAppsRevokeLinkedAppWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type LinkedAppsRevokeLinkedAppBatchWrapper struct {
	apierror.ApiError
	EndpointError *team.RevokeLinkedAppBatchError `json:"error"`
}

func (dbx *apiImpl) LinkedAppsRevokeLinkedAppBatch(arg *team.RevokeLinkedApiAppBatchArg) (res *team.RevokeLinkedAppBatchResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/linked_apps/revoke_linked_app_batch", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap LinkedAppsRevokeLinkedAppBatchWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersListWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersListError `json:"error"`
}

func (dbx *apiImpl) MembersList(arg *team.MembersListArg) (res *team.MembersListResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/list", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersListWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersListContinueWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersListContinueError `json:"error"`
}

func (dbx *apiImpl) MembersListContinue(arg *team.MembersListContinueArg) (res *team.MembersListResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/list/continue", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersListContinueWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersGetInfoWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersGetInfoError `json:"error"`
}

func (dbx *apiImpl) MembersGetInfo(arg *team.MembersGetInfoArgs) (res []*team.MembersGetInfoItem, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/get_info", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersGetInfoWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersAddWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) MembersAdd(arg *team.MembersAddArg) (res *team.MembersAddLaunch, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/add", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersAddWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersAddJobStatusGetWrapper struct {
	apierror.ApiError
	EndpointError *async.PollError `json:"error"`
}

func (dbx *apiImpl) MembersAddJobStatusGet(arg *async.PollArg) (res *team.MembersAddJobStatus, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/add/job_status/get", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersAddJobStatusGetWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersSetProfileWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersSetProfileError `json:"error"`
}

func (dbx *apiImpl) MembersSetProfile(arg *team.MembersSetProfileArg) (res *team.TeamMemberInfo, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/set_profile", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersSetProfileWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersSetAdminPermissionsWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersSetPermissionsError `json:"error"`
}

func (dbx *apiImpl) MembersSetAdminPermissions(arg *team.MembersSetPermissionsArg) (res *team.MembersSetPermissionsResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/set_admin_permissions", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersSetAdminPermissionsWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersSendWelcomeEmailWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersSendWelcomeError `json:"error"`
}

func (dbx *apiImpl) MembersSendWelcomeEmail(arg *team.UserSelectorArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/send_welcome_email", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersSendWelcomeEmailWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersRemoveWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersRemoveError `json:"error"`
}

func (dbx *apiImpl) MembersRemove(arg *team.MembersRemoveArg) (res *async.LaunchEmptyResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/remove", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersRemoveWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersRemoveJobStatusGetWrapper struct {
	apierror.ApiError
	EndpointError *async.PollError `json:"error"`
}

func (dbx *apiImpl) MembersRemoveJobStatusGet(arg *async.PollArg) (res *async.PollEmptyResult, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/remove/job_status/get", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersRemoveJobStatusGetWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersSuspendWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersSuspendError `json:"error"`
}

func (dbx *apiImpl) MembersSuspend(arg *team.MembersDeactivateArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/suspend", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersSuspendWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type MembersUnsuspendWrapper struct {
	apierror.ApiError
	EndpointError *team.MembersUnsuspendError `json:"error"`
}

func (dbx *apiImpl) MembersUnsuspend(arg *team.MembersUnsuspendArg) (res struct{}, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/members/unsuspend", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap MembersUnsuspendWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ReportsGetStorageWrapper struct {
	apierror.ApiError
	EndpointError *team.DateRangeError `json:"error"`
}

func (dbx *apiImpl) ReportsGetStorage(arg *team.DateRange) (res *team.GetStorageReport, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/reports/get_storage", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ReportsGetStorageWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ReportsGetActivityWrapper struct {
	apierror.ApiError
	EndpointError *team.DateRangeError `json:"error"`
}

func (dbx *apiImpl) ReportsGetActivity(arg *team.DateRange) (res *team.GetActivityReport, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/reports/get_activity", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ReportsGetActivityWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ReportsGetMembershipWrapper struct {
	apierror.ApiError
	EndpointError *team.DateRangeError `json:"error"`
}

func (dbx *apiImpl) ReportsGetMembership(arg *team.DateRange) (res *team.GetMembershipReport, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/reports/get_membership", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ReportsGetMembershipWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type ReportsGetDevicesWrapper struct {
	apierror.ApiError
	EndpointError *team.DateRangeError `json:"error"`
}

func (dbx *apiImpl) ReportsGetDevices(arg *team.DateRange) (res *team.GetDevicesReport, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/team/reports/get_devices", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap ReportsGetDevicesWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetAccountWrapper struct {
	apierror.ApiError
	EndpointError *users.GetAccountError `json:"error"`
}

func (dbx *apiImpl) GetAccount(arg *users.GetAccountArg) (res *users.BasicAccount, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/users/get_account", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetAccountWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetCurrentAccountWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) GetCurrentAccount() (res *users.FullAccount, err error) {
	cli := dbx.client

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/users/get_current_account", nil)
	if err != nil {
		return
	}

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetCurrentAccountWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetSpaceUsageWrapper struct {
	apierror.ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) GetSpaceUsage() (res *users.SpaceUsage, err error) {
	cli := dbx.client

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/users/get_space_usage", nil)
	if err != nil {
		return
	}

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetSpaceUsageWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}

type GetAccountBatchWrapper struct {
	apierror.ApiError
	EndpointError *users.GetAccountBatchError `json:"error"`
}

func (dbx *apiImpl) GetAccountBatch(arg *users.GetAccountBatchArg) (res []*users.BasicAccount, err error) {
	cli := dbx.client

	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/users/get_account_batch", bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			var errWrap GetAccountBatchWrapper
			err = json.Unmarshal(body, &errWrap)
			if err != nil {
				return
			}
			err = errWrap
			return
		}
		var apiError apierror.ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	return
}
