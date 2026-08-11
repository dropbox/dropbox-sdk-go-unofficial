# File transfer examples

These examples demonstrate reliable Dropbox uploads and downloads using the
`dropbox/filetransfer` package.

Each subdirectory is an independent `main` package and can be run from this
directory, for example:

```sh
cd examples/go/filetransfer

DROPBOX_ACCESS_TOKEN=YOUR_TOKEN \
DROPBOX_DOWNLOAD_PATH=/large-file.bin \
LOCAL_DOWNLOAD_PATH=large-file.bin \
go run ./download-file
```

## Examples

- `download-file`: download to a local file.
- `download-bytes`: download into memory.
- `download-progress`: report sequential download progress.
- `download-parallel`: use parallel ranged downloads.
- `upload-file`: upload a local file.
- `upload-bytes`: upload an in-memory byte slice.
- `upload-progress`: report upload progress.
- `upload-http`: stream an HTTP response into Dropbox.

## Environment

All examples require `DROPBOX_ACCESS_TOKEN`.

Download examples use `DROPBOX_DOWNLOAD_PATH`. File-based examples also use
`LOCAL_DOWNLOAD_PATH`.

Upload examples use `DROPBOX_UPLOAD_PATH`. File-based examples also use
`LOCAL_UPLOAD_PATH`. The HTTP example uses `SOURCE_URL`.
