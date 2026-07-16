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
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/team"
)

// TestTeamMembersList exercises a team-scoped endpoint with team credentials.
func TestTeamMembersList(t *testing.T) {
	client := team.New(teamConfig(t))

	if _, err := client.MembersList(team.NewMembersListArg()); err != nil {
		t.Fatalf("MembersList: %v", err)
	}
}

// TestTeamMemberCreateAndDeleteFolder acts as a team member (the team analogue
// of the user create/delete round-trip) to create a uniquely named folder in
// that member's Dropbox and then delete it. It requires team credentials and a
// team with at least one member.
func TestTeamMemberCreateAndDeleteFolder(t *testing.T) {
	cfg := teamConfig(t)

	members, err := team.New(cfg).MembersList(team.NewMembersListArg())
	if err != nil {
		t.Fatalf("MembersList: %v", err)
	}
	if len(members.Members) == 0 {
		t.Skip("team has no members; skipping member file operations")
	}
	memberID := members.Members[0].Profile.TeamMemberId

	// Perform file operations as the selected member (equivalent to the Python
	// SDK's DropboxTeam.as_user).
	cfg.AsMemberID = memberID
	client := files.New(cfg)

	path := fmt.Sprintf("/sdk-integration-team-%d", time.Now().UnixNano())

	created, err := client.CreateFolderV2(files.NewCreateFolderArg(path))
	if err != nil {
		t.Fatalf("CreateFolderV2(%q) as member %q: %v", path, memberID, err)
	}
	if created.Metadata.PathDisplay != path {
		t.Fatalf("created folder path = %q, want %q", created.Metadata.PathDisplay, path)
	}

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
		t.Fatalf("DeleteV2(%q) as member %q: %v", path, memberID, err)
	}
	deleted = true
	folder, ok := result.Metadata.(*files.FolderMetadata)
	if !ok {
		t.Fatalf("deleted metadata type = %T, want *files.FolderMetadata", result.Metadata)
	}
	if folder.PathDisplay != path {
		t.Fatalf("deleted folder path = %q, want %q", folder.PathDisplay, path)
	}
}
