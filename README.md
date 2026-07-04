# Dropbox SDK for Go [UNOFFICIAL] [![GoDoc](https://pkg.go.dev/badge/github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox)](https://pkg.go.dev/github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox) [![Actions Status](https://github.com/dropbox/dropbox-sdk-go-unofficial/workflows/Test/badge.svg)](https://github.com/dropbox/dropbox-sdk-go-unofficial/actions) [![Actions Status](https://github.com/dropbox/dropbox-sdk-go-unofficial/workflows/Lint/badge.svg)](https://github.com/dropbox/dropbox-sdk-go-unofficial/actions)

An **UNOFFICIAL** Go SDK for integrating with the Dropbox API v2. Requires Go 1.23+

:warning: WARNING: This SDK is **NOT yet official**. What does this mean?

  * There is no formal Dropbox [support](https://www.dropbox.com/developers/support) for this SDK at this point
  * Bugs may or may not get fixed
  * Not all SDK features may be implemented and implemented features may be buggy or incorrect


### Uh OK, so why are you releasing this?

  * the SDK, while unofficial, _is_ usable. See [dbxcli](https://github.com/dropbox/dbxcli) for an example application built using the SDK
  * we would like to get feedback from the community and evaluate the level of interest/enthusiasm before investing into official supporting one more SDK

## Installation

```sh
$ go get github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/...
```

For most applications, you should just import the relevant namespace(s) only. The SDK exports the following sub-packages:

* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/team`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users`

Additionally, the base `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox` package exports some configuration and helper methods.

## Usage

First, you need to [register a new "app"](https://dropbox.com/developers/apps) to start making API requests. Once you have created an app, you can either use the SDK via an access token (useful for testing) or via the regular OAuth2 flow (recommended for production).

### Using OAuth token

Once you've created an app, you can get an access token from the app's console. Note that this token will only work for the Dropbox account the token is associated with.

```go
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"

func main() {
  config := dropbox.Config{
      Token: token,
      LogLevel: dropbox.LogInfo, // if needed, set the desired logging level. Default is off
  }
  dbx := users.New(config)
  // start making API calls
}
```

### Using OAuth2 flow

For this, you will need your `APP_KEY` and `APP_SECRET` from the developers console. Your app will then have to take users though the oauth flow, as part of which users will explicitly grant permissions to your app. At the end of this process, users will get a token that the app can then use for subsequent authentication. See [this](https://pkg.go.dev/golang.org/x/oauth2#example-Config) for an example of oauth2 flow in Go.

Once you have the token, usage is same as above.

### Making API calls

Each Dropbox API takes in a request type and returns a response type. For instance, [/users/get_account](https://www.dropbox.com/developers/documentation/http/documentation#users-get_account) takes as input a `GetAccountArg` and returns a `BasicAccount`. The typical pattern for making API calls is:

  * Instantiate the argument via the `New*` convenience functions in the SDK
  * Invoke the API
  * Process the response (or handle error, as below)

Here's an example:

```go
  arg := users.NewGetAccountArg(accountId)
  if resp, err := dbx.GetAccount(arg); err != nil {
    return err
  } else {
    fmt.Printf("Name: %v", resp.Name)
  }
```

### Automatic upload integrity (`content_hash`)

For file upload calls whose argument includes `ContentHash`, the SDK automatically
computes a Dropbox `content_hash` when the content reader implements `io.ReadSeeker`.
This covers `Upload`, `UploadSessionStart`, `UploadSessionAppendV2`,
`UploadSessionAppendBatch`, and `UploadSessionFinish`. The server rejects the
request if the hash does not match the received bytes, providing end-to-end
integrity protection with no code changes required.

Because the hash must appear in the HTTP request header, the SDK reads the content once
to compute the hash, seeks back, then reads again for the upload body. For local files
the second read is typically served from the OS page cache. To disable auto-hashing:

```go
_, err = client.Upload(arg, files.WithoutAutoContentHash(f))
```

The SDK treats a reader as seekable when it implements `io.ReadSeeker`.
Non-seekable readers (e.g. `io.Pipe`, `bytes.Buffer`) are never hashed
automatically because the SDK cannot rewind them for the upload body.

To set `content_hash` manually instead, assign it on the arg. The SDK skips
auto-hashing when a value is already present, even if the reader is wrapped with
`WithoutAutoContentHash`. This also works for non-seekable readers if you already
know the Dropbox content hash:

```go
arg := files.NewUploadArg("/remote-file.txt")
arg.ContentHash = myPrecomputedHash
_, err = client.Upload(arg, f)
```

### Computing content hashes manually

Use the `contenthash` package directly when you need the hash value (e.g. for
local-vs-remote comparison):

```go
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"

f, err := os.Open("local-file.txt")
if err != nil {
    return err
}
defer f.Close()

hash, err := contenthash.Compute(f)
if err != nil {
    return err
}
fmt.Println(hash) // hex-encoded Dropbox content_hash
```

### Working with polymorphic responses

Some API methods return interface types (e.g. `IsSharedLinkMetadata`, `IsMetadata`). Use a type switch to access the concrete type:

```go
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"

// CreateSharedLinkWithSettings returns IsSharedLinkMetadata
res, err := client.CreateSharedLinkWithSettings(arg)
if err != nil {
    return err
}
switch link := res.(type) {
case *sharing.FileLinkMetadata:
    fmt.Println("File link:", link.Url)
case *sharing.FolderLinkMetadata:
    fmt.Println("Folder link:", link.Url)
case *sharing.SharedLinkMetadata:
    fmt.Println("Link:", link.Url)
}
```

Similarly, when listing folder contents, entries are returned as `IsMetadata`:

```go
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"

res, err := client.ListFolder(files.NewListFolderArg("/path"))
if err != nil {
    return err
}
for _, entry := range res.Entries {
    switch e := entry.(type) {
    case *files.FileMetadata:
        fmt.Printf("File: %s (%d bytes)\n", e.Name, e.Size)
    case *files.FolderMetadata:
        fmt.Printf("Folder: %s\n", e.Name)
    case *files.DeletedMetadata:
        fmt.Printf("Deleted: %s\n", e.Name)
    }
}
```

### Timestamps

All timestamp fields use `dropbox.DBXTime` which serializes to the format the Dropbox API expects (`2006-01-02T15:04:05Z`). Convert to/from `time.Time`:

```go
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"

// time.Time → DBXTime
modified := dropbox.DBXTime(time.Now())

// DBXTime → time.Time
t := time.Time(metadata.ClientModified)
```

### Error Handling

As described in the [API docs](https://www.dropbox.com/developers/documentation/http/documentation#error-handling), all HTTP errors _except_ 409 are returned as-is to the client (with a helpful text message where possible). In case of a 409, the SDK will return an endpoint-specific error as described in the API. This will be made available as `EndpointError` member in the error.

## Note on using the Teams API

To use the Team API, you will need to create a Dropbox Business App. The OAuth token from this app will _only_ work for the Team API.

Please read the [API docs](https://www.dropbox.com/developers/documentation/http/teams) carefully to appropriate secure your apps and tokens when using the Team API.

## Code Generation

This SDK is automatically generated using the public [Dropbox API spec](https://github.com/dropbox/dropbox-api-spec) and [Stone](https://github.com/dropbox/stone). See this [README](https://github.com/dropbox/dropbox-sdk-go-unofficial/blob/master/generator/README.md)
for more details on how code is generated. 

## Caveats

  * To re-iterate, this is an **UNOFFICIAL** SDK and thus has no official support from Dropbox
  * Only supports the v2 API. Parts of the v2 API are still in beta, and thus subject to change
  * This SDK itself is in beta, and so interfaces may change at any point
