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

//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// TestUserCreateAndDeleteFolder exercises a round-trip against the live user
// API: create a uniquely named folder, then delete it. It authenticates with a
// refresh token, so it also proves the SDK can mint an access token on demand.
func TestUserCreateAndDeleteFolder(t *testing.T) {
	client := files.New(userConfig(t))

	path := fmt.Sprintf("/sdk-integration-%d", time.Now().UnixNano())

	created, err := client.CreateFolderV2(files.NewCreateFolderArg(path))
	if err != nil {
		t.Fatalf("CreateFolderV2(%q): %v", path, err)
	}
	if created.Metadata.PathDisplay != path {
		t.Fatalf("created folder path = %q, want %q", created.Metadata.PathDisplay, path)
	}

	// Best-effort cleanup so an early failure below does not leak the folder.
	// Once the folder is deleted successfully this becomes a no-op.
	deleted := false
	defer func() {
		if deleted {
			return
		}
		if _, err := client.DeleteV2(files.NewDeleteArg(path)); err != nil {
			t.Errorf("cleanup DeleteV2(%q): %v", path, err)
		}
	}()

	result, err := client.DeleteV2(files.NewDeleteArg(path))
	if err != nil {
		t.Fatalf("DeleteV2(%q): %v", path, err)
	}
	deleted = true
	// DeleteResult.Metadata is the IsMetadata interface; a deleted folder comes
	// back as *FolderMetadata.
	folder, ok := result.Metadata.(*files.FolderMetadata)
	if !ok {
		t.Fatalf("deleted metadata type = %T, want *files.FolderMetadata", result.Metadata)
	}
	if folder.PathDisplay != path {
		t.Fatalf("deleted folder path = %q, want %q", folder.PathDisplay, path)
	}
}
