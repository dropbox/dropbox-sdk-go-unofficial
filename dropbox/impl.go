/* DO NOT EDIT */

package dropbox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Api interface {
	Files
	Sharing
	Users
}

type GetMetadataWrapper struct {
	ApiError
	EndpointError *GetMetadataError `json:"error"`
}

func (dbx *apiImpl) GetMetadata(arg *GetMetadataArg) (res *Metadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListFolderLongpollError `json:"error"`
}

func (dbx *apiImpl) ListFolderLongpoll(arg *ListFolderLongpollArg) (res *ListFolderLongpollResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListFolderError `json:"error"`
}

func (dbx *apiImpl) ListFolder(arg *ListFolderArg) (res *ListFolderResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListFolderContinueError `json:"error"`
}

func (dbx *apiImpl) ListFolderContinue(arg *ListFolderContinueArg) (res *ListFolderResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListFolderError `json:"error"`
}

func (dbx *apiImpl) ListFolderGetLatestCursor(arg *ListFolderArg) (res *ListFolderGetLatestCursorResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *DownloadError `json:"error"`
}

func (dbx *apiImpl) Download(arg *DownloadArg) (res *FileMetadata, content io.ReadCloser, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) UploadSessionStart(content io.Reader) (res *UploadSessionStartResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UploadSessionLookupError `json:"error"`
}

func (dbx *apiImpl) UploadSessionAppend(arg *UploadSessionCursor, content io.Reader) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UploadSessionFinishError `json:"error"`
}

func (dbx *apiImpl) UploadSessionFinish(arg *UploadSessionFinishArg, content io.Reader) (res *FileMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UploadError `json:"error"`
}

func (dbx *apiImpl) Upload(arg *CommitInfo, content io.Reader) (res *FileMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *SearchError `json:"error"`
}

func (dbx *apiImpl) Search(arg *SearchArg) (res *SearchResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *CreateFolderError `json:"error"`
}

func (dbx *apiImpl) CreateFolder(arg *CreateFolderArg) (res *FolderMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *DeleteError `json:"error"`
}

func (dbx *apiImpl) Delete(arg *DeleteArg) (res *Metadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *DeleteError `json:"error"`
}

func (dbx *apiImpl) PermanentlyDelete(arg *DeleteArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *RelocationError `json:"error"`
}

func (dbx *apiImpl) Copy(arg *RelocationArg) (res *Metadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *RelocationError `json:"error"`
}

func (dbx *apiImpl) Move(arg *RelocationArg) (res *Metadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ThumbnailError `json:"error"`
}

func (dbx *apiImpl) GetThumbnail(arg *ThumbnailArg) (res *FileMetadata, content io.ReadCloser, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *PreviewError `json:"error"`
}

func (dbx *apiImpl) GetPreview(arg *PreviewArg) (res *FileMetadata, content io.ReadCloser, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListRevisionsError `json:"error"`
}

func (dbx *apiImpl) ListRevisions(arg *ListRevisionsArg) (res *ListRevisionsResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *RestoreError `json:"error"`
}

func (dbx *apiImpl) Restore(arg *RestoreArg) (res *FileMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *SharedLinkError `json:"error"`
}

func (dbx *apiImpl) GetSharedLinkMetadata(arg *GetSharedLinkMetadataArg) (res *SharedLinkMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListSharedLinksError `json:"error"`
}

func (dbx *apiImpl) ListSharedLinks(arg *ListSharedLinksArg) (res *ListSharedLinksResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ModifySharedLinkSettingsError `json:"error"`
}

func (dbx *apiImpl) ModifySharedLinkSettings(arg *ModifySharedLinkSettingsArgs) (res *SharedLinkMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *CreateSharedLinkWithSettingsError `json:"error"`
}

func (dbx *apiImpl) CreateSharedLinkWithSettings(arg *CreateSharedLinkWithSettingsArg) (res *SharedLinkMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *RevokeSharedLinkError `json:"error"`
}

func (dbx *apiImpl) RevokeSharedLink(arg *RevokeSharedLinkArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *GetSharedLinkFileError `json:"error"`
}

func (dbx *apiImpl) GetSharedLinkFile(arg *GetSharedLinkMetadataArg) (res *SharedLinkMetadata, content io.ReadCloser, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *GetSharedLinksError `json:"error"`
}

func (dbx *apiImpl) GetSharedLinks(arg *GetSharedLinksArg) (res *GetSharedLinksResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *CreateSharedLinkError `json:"error"`
}

func (dbx *apiImpl) CreateSharedLink(arg *CreateSharedLinkArg) (res *PathLinkMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) ListFolders() (res *ListFoldersResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListFoldersContinueError `json:"error"`
}

func (dbx *apiImpl) ListFoldersContinue(arg *ListFoldersContinueArg) (res *ListFoldersResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *SharedFolderAccessError `json:"error"`
}

func (dbx *apiImpl) GetFolderMetadata(arg *GetMetadataArgs) (res *SharedFolderMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *SharedFolderAccessError `json:"error"`
}

func (dbx *apiImpl) ListFolderMembers(arg *ListFolderMembersArgs) (res *SharedFolderMembers, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ListFolderMembersContinueError `json:"error"`
}

func (dbx *apiImpl) ListFolderMembersContinue(arg *ListFolderMembersContinueArg) (res *SharedFolderMembers, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *ShareFolderError `json:"error"`
}

func (dbx *apiImpl) ShareFolder(arg *ShareFolderArg) (res *ShareFolderLaunch, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *PollError `json:"error"`
}

func (dbx *apiImpl) CheckShareJobStatus(arg *PollArg) (res *ShareFolderJobStatus, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *PollError `json:"error"`
}

func (dbx *apiImpl) CheckJobStatus(arg *PollArg) (res *JobStatus, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UnshareFolderError `json:"error"`
}

func (dbx *apiImpl) UnshareFolder(arg *UnshareFolderArg) (res *LaunchEmptyResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *TransferFolderError `json:"error"`
}

func (dbx *apiImpl) TransferFolder(arg *TransferFolderArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UpdateFolderPolicyError `json:"error"`
}

func (dbx *apiImpl) UpdateFolderPolicy(arg *UpdateFolderPolicyArg) (res *SharedFolderMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *AddFolderMemberError `json:"error"`
}

func (dbx *apiImpl) AddFolderMember(arg *AddFolderMemberArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *RemoveFolderMemberError `json:"error"`
}

func (dbx *apiImpl) RemoveFolderMember(arg *RemoveFolderMemberArg) (res *LaunchEmptyResult, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UpdateFolderMemberError `json:"error"`
}

func (dbx *apiImpl) UpdateFolderMember(arg *UpdateFolderMemberArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *MountFolderError `json:"error"`
}

func (dbx *apiImpl) MountFolder(arg *MountFolderArg) (res *SharedFolderMetadata, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *UnmountFolderError `json:"error"`
}

func (dbx *apiImpl) UnmountFolder(arg *UnmountFolderArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *RelinquishFolderMembershipError `json:"error"`
}

func (dbx *apiImpl) RelinquishFolderMembership(arg *RelinquishFolderMembershipArg) (res struct{}, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *GetAccountError `json:"error"`
}

func (dbx *apiImpl) GetAccount(arg *GetAccountArg) (res *BasicAccount, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) GetCurrentAccount() (res *FullAccount, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError struct{} `json:"error"`
}

func (dbx *apiImpl) GetSpaceUsage() (res *SpaceUsage, err error) {
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
		var apiError ApiError
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
	ApiError
	EndpointError *GetAccountBatchError `json:"error"`
}

func (dbx *apiImpl) GetAccountBatch(arg *GetAccountBatchArg) (res []*BasicAccount, err error) {
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
		var apiError ApiError
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
