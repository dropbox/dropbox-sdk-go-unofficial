# Dropbox SDK for Go [![GoDoc](https://pkg.go.dev/badge/github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox)](https://pkg.go.dev/github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox) [![Actions Status](https://github.com/dropbox/dropbox-sdk-go-unofficial/workflows/Test/badge.svg)](https://github.com/dropbox/dropbox-sdk-go-unofficial/actions) [![Actions Status](https://github.com/dropbox/dropbox-sdk-go-unofficial/workflows/Lint/badge.svg)](https://github.com/dropbox/dropbox-sdk-go-unofficial/actions)

A Go SDK for integrating with the Dropbox API v2. Requires Go 1.25+

## Installation

```sh
$ go get github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/...
```

For most applications, you should just import the relevant namespace(s) only. The SDK exports the following sub-packages:

* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/oauth`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/openid`
* `github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/retry`
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

For PKCE, you will need your `APP_KEY` from the developers console. Your app
will then have to take users through the OAuth flow, as part of which users
will explicitly grant permissions to your app. At the end of this process, your
app exchanges the authorization code for a token.

Choose the OAuth helper that matches your app type:

| App type | Helper |
| --- | --- |
| CLI or native public app | `NewPKCEFlow` |
| Web public app | `NewWebPKCEFlow` |
| Confidential app without redirect | `NewOAuth2FlowNoRedirect` |
| Confidential web app | `NewOAuth2Flow` |
| Refreshable SDK client | `TokenSource` |

```go
import "context"
import "fmt"

import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
import dropboxoauth "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/oauth"
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"

func main() {
  ctx := context.Background()
  flow, err := dropboxoauth.NewPKCEFlow(appKey)
  if err != nil {
    panic(err)
  }

  fmt.Printf("Go to %s\n", flow.AuthCodeURL())
  // Read the authorization code from the user.
  token, err := flow.Exchange(ctx, code)
  if err != nil {
    panic(err)
  }

  config := dropbox.Config{
    TokenSource: dropboxoauth.TokenSource(ctx, appKey, token),
  }
  dbx := users.New(config)
  // start making API calls
}
```

The PKCE helpers request offline access by default so refresh tokens are
returned. Use `dropboxoauth.TokenSource` in `dropbox.Config` for automatic
refresh, or call `dropboxoauth.Refresh` directly when you manage token storage
yourself. If your app only calls Dropbox while the user is actively present and
does not need background access, pass
`dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessTypeOnline)`.

For redirect-based web apps, use `NewWebPKCEFlow`. Each login attempt needs a
single PKCE verifier for both `Start` and `Finish`; generate it with
`dropboxoauth.GenerateVerifier`, store it with the returned CSRF token in the
user's session, then pass both values back when finishing the flow:

```go
verifier, err := dropboxoauth.GenerateVerifier()
if err != nil {
  return err
}

flow, err := dropboxoauth.NewWebPKCEFlow(
  appKey,
  redirectURL,
  dropboxoauth.WithVerifier(verifier),
)
if err != nil {
  return err
}

authURL, csrfToken, err := flow.Start("optional-url-state")
if err != nil {
  return err
}
// Redirect the user to authURL and store csrfToken and verifier in your web session.

// In the callback handler, load csrfToken and verifier from your web session.
flow, err = dropboxoauth.NewWebPKCEFlow(
  appKey,
  redirectURL,
  dropboxoauth.WithVerifier(verifier),
)
if err != nil {
  return err
}

result, err := flow.Finish(r.Context(), r.URL.Query(), csrfToken)
if err != nil {
  return err
}

config := dropbox.Config{
  TokenSource: dropboxoauth.TokenSource(r.Context(), appKey, result.Token),
}
```

Confidential apps that can safely keep an app secret can use
`NewOAuth2FlowNoRedirect` or `NewOAuth2Flow` instead. These helpers support the
same token access types as Dropbox's Python SDK: omit the token access type to
use the app default, use `TokenAccessTypeOffline` for refresh tokens, or use the
deprecated `TokenAccessTypeLegacy` only for legacy compatibility.

```go
flow, err := dropboxoauth.NewOAuth2FlowNoRedirect(
  appKey,
  dropboxoauth.WithAppSecret(appSecret),
  dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessTypeOffline),
)
if err != nil {
  return err
}

fmt.Printf("Go to %s\n", flow.Start())
// Read the authorization code from the user.
result, err := flow.Finish(ctx, code)
if err != nil {
  return err
}

config := dropbox.Config{
  TokenSource: dropboxoauth.TokenSource(
    ctx,
    appKey,
    result.Token,
    dropboxoauth.WithAppSecret(appSecret),
  ),
}
```

The SDK still accepts a static access token through `dropbox.Config.Token`; for
refreshable tokens, pass a token source through `dropbox.Config.TokenSource`.

If your app uses Dropbox as a login identity provider, request OpenID Connect
scopes during OAuth, then call `UserinfoContext` on an `openid` client with the
resulting token. Normal Dropbox API access does not require OpenID scopes.

```go
flow, err := dropboxoauth.NewWebPKCEFlow(
  appKey,
  redirectURL,
  dropboxoauth.WithVerifier(verifier),
  dropboxoauth.WithScopes(
    dropboxoauth.ScopeOpenID,
    dropboxoauth.ScopeEmail,
    dropboxoauth.ScopeProfile,
  ),
)
if err != nil {
  return err
}

// Finish the OAuth flow as shown above.
identityClient := openid.NewContext(dropbox.Config{
  TokenSource: dropboxoauth.TokenSource(ctx, appKey, result.Token),
})
info, err := identityClient.UserinfoContext(ctx, openid.NewUserInfoArgs())
if err != nil {
  return err
}
_ = info
```

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

### Retrying API calls

Retries are disabled by default. To opt in for all calls made by a client, set
`Config.RetryPolicy`:

```go
import "time"

import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/retry"
import "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"

policy := retry.Policy{
  MaxRetries:     3,
  InitialBackoff: 500 * time.Millisecond,
  MaxBackoff:     30 * time.Second,
  MaxRetryAfter:  30 * time.Second,
}

dbx := users.New(dropbox.Config{
  Token:       accessToken,
  RetryPolicy: &policy,
})

arg := users.NewGetAccountArg(accountId)
resp, err := dbx.GetAccount(arg)
```

`MaxRetries` is the number of attempts after the initial request. If left unset,
`InitialBackoff`, `MaxBackoff`, and `MaxRetryAfter` default to 500ms, 30s, and
30s respectively.

The SDK retries HTTP `429 Too Many Requests` and `503 Service Unavailable`
responses. `Retry-After` headers and Dropbox 429 `retry_after` response bodies
are honored up to `MaxRetryAfter`; longer delays are not retried. Network errors
where no HTTP response is received are not retried.

Retry policies are client-wide; there is no per-call retry override. To set a
per-call time budget, use a context deadline, which covers retry sleeps and all
attempts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := dbx.GetAccountContext(ctx, arg)
```

Some Dropbox API `409 Conflict` responses are also safe to retry. Opt in by
listing the top-level Stone error tags that should be retried:

```go
policy := retry.Policy{
  MaxRetries:       3,
  InitialBackoff:   time.Second,
  MaxBackoff:       30 * time.Second,
  MaxRetryAfter:    30 * time.Second,
  Retryable409Tags: []string{"too_many_write_operations"},
}
```

Upload retries require a seekable request body (for example `*os.File`,
`*bytes.Reader`, or `*strings.Reader`) so the SDK can rewind the body before
each attempt. Non-seekable upload bodies are sent once and are not retried.
Leave `Config.RetryPolicy` nil to disable retries.

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

## Support posture

This SDK is maintained in the Dropbox GitHub organization by Dropbox engineers,
but it is not a formally supported Dropbox product. Use GitHub issues and pull
requests for bugs and contributions; Dropbox [Support](https://www.dropbox.com/developers/support)
does not provide support for this SDK. The SDK implements a practical subset of
Dropbox API features, not the full API surface.
