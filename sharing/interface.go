/* DO NOT EDIT */
/* This file was generated from sharing.babel */

package sharing

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dropbox/dropbox-sdk-go/async"
	"github.com/dropbox/dropbox-sdk-go/files"
	"github.com/dropbox/dropbox-sdk-go/users"
)

// Defines the access levels for collaborators.
type AccessLevel struct {
	Tag string `json:".tag"`
}

// Policy governing who can change a shared folder's access control list (ACL).
// In other words, who can add, remove, or change the privileges of members.
type AclUpdatePolicy struct {
	Tag string `json:".tag"`
}

type AddFolderMemberArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// The intended list of members to add.  Added members will receive invites to
	// join the shared folder.
	Members []*AddMember `json:"members"`
	// Whether added members should be notified via email and device notifications
	// of their invite.
	Quiet bool `json:"quiet"`
	// Optional message to display to added members in their invitation.
	CustomMessage string `json:"custom_message,omitempty"`
}

func NewAddFolderMemberArg() *AddFolderMemberArg {
	s := new(AddFolderMemberArg)
	s.Quiet = false
	return s
}

type AddFolderMemberError struct {
	Tag string `json:".tag"`
	// Unable to access shared folder.
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
	// :field:`AddFolderMemberArg.members` contains a bad invitation recipient.
	BadMember *AddMemberSelectorError `json:"bad_member,omitempty"`
	// The value is the member limit that was reached.
	TooManyMembers uint64 `json:"too_many_members,omitempty"`
	// The value is the pending invite limit that was reached.
	TooManyPendingInvites uint64 `json:"too_many_pending_invites,omitempty"`
}

func (u *AddFolderMemberError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// Unable to access shared folder.
		AccessError json.RawMessage `json:"access_error"`
		// :field:`AddFolderMemberArg.members` contains a bad invitation recipient.
		BadMember json.RawMessage `json:"bad_member"`
		// The value is the member limit that was reached.
		TooManyMembers json.RawMessage `json:"too_many_members"`
		// The value is the pending invite limit that was reached.
		TooManyPendingInvites json.RawMessage `json:"too_many_pending_invites"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	case "bad_member":
		{
			if len(w.BadMember) == 0 {
				break
			}
			if err := json.Unmarshal(w.BadMember, &u.BadMember); err != nil {
				return err
			}
		}
	case "too_many_members":
		{
			if len(w.TooManyMembers) == 0 {
				break
			}
			if err := json.Unmarshal(w.TooManyMembers, &u.TooManyMembers); err != nil {
				return err
			}
		}
	case "too_many_pending_invites":
		{
			if len(w.TooManyPendingInvites) == 0 {
				break
			}
			if err := json.Unmarshal(w.TooManyPendingInvites, &u.TooManyPendingInvites); err != nil {
				return err
			}
		}
	}
	return nil
}

// The member and type of access the member should have when added to a shared
// folder.
type AddMember struct {
	// The member to add to the shared folder.
	Member *MemberSelector `json:"member"`
	// The access level to grant :field:`member` to the shared folder.
	// :field:`AccessLevel.owner` is disallowed.
	AccessLevel *AccessLevel `json:"access_level"`
}

func NewAddMember() *AddMember {
	s := new(AddMember)
	s.AccessLevel = &AccessLevel{Tag: "viewer"}
	return s
}

type AddMemberSelectorError struct {
	Tag string `json:".tag"`
	// The value is the ID that could not be identified.
	InvalidDropboxId string `json:"invalid_dropbox_id,omitempty"`
	// The value is the e-email address that is malformed.
	InvalidEmail string `json:"invalid_email,omitempty"`
	// The value is the ID of the Dropbox user with an unverified e-mail address.
	// Invite unverified users by e-mail address instead of by their Dropbox ID.
	UnverifiedDropboxId string `json:"unverified_dropbox_id,omitempty"`
}

func (u *AddMemberSelectorError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The value is the ID that could not be identified.
		InvalidDropboxId json.RawMessage `json:"invalid_dropbox_id"`
		// The value is the e-email address that is malformed.
		InvalidEmail json.RawMessage `json:"invalid_email"`
		// The value is the ID of the Dropbox user with an unverified e-mail address.
		// Invite unverified users by e-mail address instead of by their Dropbox ID.
		UnverifiedDropboxId json.RawMessage `json:"unverified_dropbox_id"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "invalid_dropbox_id":
		{
			if len(w.InvalidDropboxId) == 0 {
				break
			}
			if err := json.Unmarshal(w.InvalidDropboxId, &u.InvalidDropboxId); err != nil {
				return err
			}
		}
	case "invalid_email":
		{
			if len(w.InvalidEmail) == 0 {
				break
			}
			if err := json.Unmarshal(w.InvalidEmail, &u.InvalidEmail); err != nil {
				return err
			}
		}
	case "unverified_dropbox_id":
		{
			if len(w.UnverifiedDropboxId) == 0 {
				break
			}
			if err := json.Unmarshal(w.UnverifiedDropboxId, &u.UnverifiedDropboxId); err != nil {
				return err
			}
		}
	}
	return nil
}

// Metadata for a shared link. This can be either a :type:`PathLinkMetadata` or
// :type:`CollectionLinkMetadata`.
type LinkMetadata struct {
	Tag        string                  `json:".tag"`
	Path       *PathLinkMetadata       `json:"path,omitempty"`
	Collection *CollectionLinkMetadata `json:"collection,omitempty"`
}

func (u *LinkMetadata) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag        string          `json:".tag"`
		Path       json.RawMessage `json:"path"`
		Collection json.RawMessage `json:"collection"`
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
	case "collection":
		{
			if err := json.Unmarshal(body, &u.Collection); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewLinkMetadata() *LinkMetadata {
	s := new(LinkMetadata)
	return s
}

// Metadata for a collection-based shared link.
type CollectionLinkMetadata struct {
	// URL of the shared link.
	Url string `json:"url"`
	// Who can access the link.
	Visibility *Visibility `json:"visibility"`
	// Expiration time, if set. By default the link won't expire.
	Expires time.Time `json:"expires,omitempty"`
}

func NewCollectionLinkMetadata() *CollectionLinkMetadata {
	s := new(CollectionLinkMetadata)
	return s
}

type CreateSharedLinkArg struct {
	// The path to share.
	Path string `json:"path"`
	// Whether to return a shortened URL.
	ShortUrl bool `json:"short_url"`
	// If it's okay to share a path that does not yet exist, set this to either
	// :field:`PendingUploadMode.file` or :field:`PendingUploadMode.folder` to
	// indicate whether to assume it's a file or folder.
	PendingUpload *PendingUploadMode `json:"pending_upload,omitempty"`
}

func NewCreateSharedLinkArg() *CreateSharedLinkArg {
	s := new(CreateSharedLinkArg)
	s.ShortUrl = false
	return s
}

type CreateSharedLinkError struct {
	Tag  string             `json:".tag"`
	Path *files.LookupError `json:"path,omitempty"`
}

func (u *CreateSharedLinkError) UnmarshalJSON(body []byte) error {
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

type CreateSharedLinkWithSettingsArg struct {
	// The path to be shared by the shared link
	Path string `json:"path"`
	// The requested settings for the newly created shared link
	Settings *SharedLinkSettings `json:"settings,omitempty"`
}

func NewCreateSharedLinkWithSettingsArg() *CreateSharedLinkWithSettingsArg {
	s := new(CreateSharedLinkWithSettingsArg)
	return s
}

type CreateSharedLinkWithSettingsError struct {
	Tag  string             `json:".tag"`
	Path *files.LookupError `json:"path,omitempty"`
	// There is an error with the given settings
	SettingsError *SharedLinkSettingsError `json:"settings_error,omitempty"`
}

func (u *CreateSharedLinkWithSettingsError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path"`
		// There is an error with the given settings
		SettingsError json.RawMessage `json:"settings_error"`
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
	case "settings_error":
		{
			if len(w.SettingsError) == 0 {
				break
			}
			if err := json.Unmarshal(w.SettingsError, &u.SettingsError); err != nil {
				return err
			}
		}
	}
	return nil
}

// The metadata of a shared link
type SharedLinkMetadata struct {
	Tag    string              `json:".tag"`
	File   *FileLinkMetadata   `json:"file,omitempty"`
	Folder *FolderLinkMetadata `json:"folder,omitempty"`
}

func (u *SharedLinkMetadata) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag    string          `json:".tag"`
		File   json.RawMessage `json:"file"`
		Folder json.RawMessage `json:"folder"`
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
	}
	return nil
}

func NewSharedLinkMetadata() *SharedLinkMetadata {
	s := new(SharedLinkMetadata)
	return s
}

// The metadata of a file shared link
type FileLinkMetadata struct {
	// URL of the shared link.
	Url string `json:"url"`
	// The linked file name (including extension). This never contains a slash.
	Name string `json:"name"`
	// The link's access permissions.
	LinkPermissions *LinkPermissions `json:"link_permissions"`
	// The modification time set by the desktop client when the file was added to
	// Dropbox. Since this time is not verified (the Dropbox server stores whatever
	// the desktop client sends up), this should only be used for display purposes
	// (such as sorting) and not, for example, to determine if a file has changed
	// or not.
	ClientModified time.Time `json:"client_modified"`
	// The last time the file was modified on Dropbox.
	ServerModified time.Time `json:"server_modified"`
	// A unique identifier for the current revision of a file. This field is the
	// same rev as elsewhere in the API and can be used to detect changes and avoid
	// conflicts.
	Rev string `json:"rev"`
	// The file size in bytes.
	Size uint64 `json:"size"`
	// A unique identifier for the linked file.
	Id string `json:"id,omitempty"`
	// Expiration time, if set. By default the link won't expire.
	Expires time.Time `json:"expires,omitempty"`
	// The lowercased full path in the user's Dropbox. This always starts with a
	// slash. This field will only be present only if the linked file is in the
	// authenticated user's  dropbox.
	PathLower string `json:"path_lower,omitempty"`
	// The team membership information of the link's owner.  This field will only
	// be present  if the link's owner is a team member.
	TeamMemberInfo *TeamMemberInfo `json:"team_member_info,omitempty"`
	// The team information of the content's owner. This field will only be present
	// if the content's owner is a team member and the content's owner team is
	// different from the link's owner team.
	ContentOwnerTeamInfo *users.Team `json:"content_owner_team_info,omitempty"`
}

func NewFileLinkMetadata() *FileLinkMetadata {
	s := new(FileLinkMetadata)
	return s
}

// Actions that may be taken on shared folders.
type FolderAction struct {
	Tag string `json:".tag"`
}

// The metadata of a folder shared link
type FolderLinkMetadata struct {
	// URL of the shared link.
	Url string `json:"url"`
	// The linked file name (including extension). This never contains a slash.
	Name string `json:"name"`
	// The link's access permissions.
	LinkPermissions *LinkPermissions `json:"link_permissions"`
	// A unique identifier for the linked file.
	Id string `json:"id,omitempty"`
	// Expiration time, if set. By default the link won't expire.
	Expires time.Time `json:"expires,omitempty"`
	// The lowercased full path in the user's Dropbox. This always starts with a
	// slash. This field will only be present only if the linked file is in the
	// authenticated user's  dropbox.
	PathLower string `json:"path_lower,omitempty"`
	// The team membership information of the link's owner.  This field will only
	// be present  if the link's owner is a team member.
	TeamMemberInfo *TeamMemberInfo `json:"team_member_info,omitempty"`
	// The team information of the content's owner. This field will only be present
	// if the content's owner is a team member and the content's owner team is
	// different from the link's owner team.
	ContentOwnerTeamInfo *users.Team `json:"content_owner_team_info,omitempty"`
}

func NewFolderLinkMetadata() *FolderLinkMetadata {
	s := new(FolderLinkMetadata)
	return s
}

// Whether the user is allowed to take the action on the shared folder.
type FolderPermission struct {
	// The action that the user may wish to take on the folder.
	Action *FolderAction `json:"action"`
	// True if the user is allowed to take the action.
	Allow bool `json:"allow"`
	// The reason why the user is denied the permission. Not present if the action
	// is allowed
	Reason *PermissionDeniedReason `json:"reason,omitempty"`
}

func NewFolderPermission() *FolderPermission {
	s := new(FolderPermission)
	return s
}

// A set of policies governing membership and privileges for a shared folder.
type FolderPolicy struct {
	// Who can add and remove members from this shared folder.
	AclUpdatePolicy *AclUpdatePolicy `json:"acl_update_policy"`
	// Who links can be shared with.
	SharedLinkPolicy *SharedLinkPolicy `json:"shared_link_policy"`
	// Who can be a member of this shared folder. Only set if the user is a member
	// of a team.
	MemberPolicy *MemberPolicy `json:"member_policy,omitempty"`
}

func NewFolderPolicy() *FolderPolicy {
	s := new(FolderPolicy)
	return s
}

type GetMetadataArgs struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// Folder actions to query.
	Actions []*FolderAction `json:"actions,omitempty"`
}

func NewGetMetadataArgs() *GetMetadataArgs {
	s := new(GetMetadataArgs)
	return s
}

type SharedLinkError struct {
	Tag string `json:".tag"`
}

type GetSharedLinkFileError struct {
	Tag string `json:".tag"`
}

type GetSharedLinkMetadataArg struct {
	// URL of the shared link.
	Url string `json:"url"`
	// If the shared link is to a folder, this parameter can be used to retrieve
	// the metadata for a specific file or sub-folder in this folder. A relative
	// path should be used.
	Path string `json:"path,omitempty"`
	// If the shared link has a password, this parameter can be used.
	LinkPassword string `json:"link_password,omitempty"`
}

func NewGetSharedLinkMetadataArg() *GetSharedLinkMetadataArg {
	s := new(GetSharedLinkMetadataArg)
	return s
}

type GetSharedLinksArg struct {
	// See :route:`get_shared_links` description.
	Path string `json:"path,omitempty"`
}

func NewGetSharedLinksArg() *GetSharedLinksArg {
	s := new(GetSharedLinksArg)
	return s
}

type GetSharedLinksError struct {
	Tag  string `json:".tag"`
	Path string `json:"path,omitempty"`
}

func (u *GetSharedLinksError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag  string          `json:".tag"`
		Path json.RawMessage `json:"path,omitempty"`
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

type GetSharedLinksResult struct {
	// Shared links applicable to the path argument.
	Links []*LinkMetadata `json:"links"`
}

func NewGetSharedLinksResult() *GetSharedLinksResult {
	s := new(GetSharedLinksResult)
	return s
}

// The information about a group. Groups is a way to manage a list of users  who
// need same access permission to the shared folder.
type GroupInfo struct {
	GroupName string `json:"group_name"`
	GroupId   string `json:"group_id"`
	// The number of members in the group.
	MemberCount uint32 `json:"member_count"`
	// If the group is owned by the current user's team.
	SameTeam bool `json:"same_team"`
	// External ID of group. This is an arbitrary ID that an admin can attach to a
	// group.
	GroupExternalId string `json:"group_external_id,omitempty"`
}

func NewGroupInfo() *GroupInfo {
	s := new(GroupInfo)
	return s
}

// The information about a member of the shared folder.
type MembershipInfo struct {
	// The access type for this member.
	AccessType *AccessLevel `json:"access_type"`
	// The permissions that requesting user has on this member. The set of
	// permissions corresponds to the MemberActions in the request.
	Permissions []*MemberPermission `json:"permissions,omitempty"`
}

func NewMembershipInfo() *MembershipInfo {
	s := new(MembershipInfo)
	return s
}

// The information about a group member of the shared folder.
type GroupMembershipInfo struct {
	// The access type for this member.
	AccessType *AccessLevel `json:"access_type"`
	// The information about the membership group.
	Group *GroupInfo `json:"group"`
	// The permissions that requesting user has on this member. The set of
	// permissions corresponds to the MemberActions in the request.
	Permissions []*MemberPermission `json:"permissions,omitempty"`
}

func NewGroupMembershipInfo() *GroupMembershipInfo {
	s := new(GroupMembershipInfo)
	return s
}

// The information about a user invited to become a member a shared folder.
type InviteeInfo struct {
	Tag string `json:".tag"`
	// E-mail address of invited user.
	Email string `json:"email,omitempty"`
}

func (u *InviteeInfo) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// E-mail address of invited user.
		Email json.RawMessage `json:"email"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "email":
		{
			if len(w.Email) == 0 {
				break
			}
			if err := json.Unmarshal(w.Email, &u.Email); err != nil {
				return err
			}
		}
	}
	return nil
}

// The information about a user invited to become a member of a shared folder.
type InviteeMembershipInfo struct {
	// The access type for this member.
	AccessType *AccessLevel `json:"access_type"`
	// The information for the invited user.
	Invitee *InviteeInfo `json:"invitee"`
	// The permissions that requesting user has on this member. The set of
	// permissions corresponds to the MemberActions in the request.
	Permissions []*MemberPermission `json:"permissions,omitempty"`
}

func NewInviteeMembershipInfo() *InviteeMembershipInfo {
	s := new(InviteeMembershipInfo)
	return s
}

type JobError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
	MemberError *SharedFolderMemberError `json:"member_error,omitempty"`
}

func (u *JobError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
		MemberError json.RawMessage `json:"member_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	case "member_error":
		{
			if len(w.MemberError) == 0 {
				break
			}
			if err := json.Unmarshal(w.MemberError, &u.MemberError); err != nil {
				return err
			}
		}
	}
	return nil
}

type JobStatus struct {
	Tag string `json:".tag"`
	// The asynchronous job returned an error.
	Failed *JobError `json:"failed,omitempty"`
}

func (u *JobStatus) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The asynchronous job returned an error.
		Failed json.RawMessage `json:"failed"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "failed":
		{
			if len(w.Failed) == 0 {
				break
			}
			if err := json.Unmarshal(w.Failed, &u.Failed); err != nil {
				return err
			}
		}
	}
	return nil
}

type LinkPermissions struct {
	// Whether the caller can revoke the shared link
	CanRevoke bool `json:"can_revoke"`
	// The current visibility of the link after considering the shared links
	// policies of the the team (in case the link's owner is part of a team) and
	// the shared folder (in case the linked file is part of a shared folder). This
	// field is shown only if the caller has access to this info (the link's owner
	// always has access to this data).
	ResolvedVisibility *ResolvedVisibility `json:"resolved_visibility,omitempty"`
	// The shared link's requested visibility. This can be overridden by the team
	// and shared folder policies. The final visibility, after considering these
	// policies, can be found in :field:`resolved_visibility`. This is shown only
	// if the caller is the link's owner.
	RequestedVisibility *RequestedVisibility `json:"requested_visibility,omitempty"`
	// The failure reason for revoking the link. This field will only be present if
	// the :field:`can_revoke` is :val:`false`.
	RevokeFailureReason *SharedLinkAccessFailureReason `json:"revoke_failure_reason,omitempty"`
}

func NewLinkPermissions() *LinkPermissions {
	s := new(LinkPermissions)
	return s
}

type ListFolderMembersArgs struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// Member actions to query.
	Actions []*MemberAction `json:"actions,omitempty"`
}

func NewListFolderMembersArgs() *ListFolderMembersArgs {
	s := new(ListFolderMembersArgs)
	return s
}

type ListFolderMembersContinueArg struct {
	// The cursor returned by your last call to :route:`list_folder_members` or
	// :route:`list_folder_members/continue`.
	Cursor string `json:"cursor"`
}

func NewListFolderMembersContinueArg() *ListFolderMembersContinueArg {
	s := new(ListFolderMembersContinueArg)
	return s
}

type ListFolderMembersContinueError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *ListFolderMembersContinueError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

type ListFoldersContinueArg struct {
	// The cursor returned by your last call to :route:`list_folders` or
	// :route:`list_folders/continue`.
	Cursor string `json:"cursor"`
}

func NewListFoldersContinueArg() *ListFoldersContinueArg {
	s := new(ListFoldersContinueArg)
	return s
}

type ListFoldersContinueError struct {
	Tag string `json:".tag"`
}

// Result for :route:`list_folders`. Unmounted shared folders can be identified
// by the absence of :field:`SharedFolderMetadata.path_lower`.
type ListFoldersResult struct {
	// List of all shared folders the authenticated user has access to.
	Entries []*SharedFolderMetadata `json:"entries"`
	// Present if there are additional shared folders that have not been returned
	// yet. Pass the cursor into :route:`list_folders/continue` to list additional
	// folders.
	Cursor string `json:"cursor,omitempty"`
}

func NewListFoldersResult() *ListFoldersResult {
	s := new(ListFoldersResult)
	return s
}

type ListSharedLinksArg struct {
	// See :route:`list_shared_links` description.
	Path string `json:"path,omitempty"`
	// The cursor returned by your last call to :route:`list_shared_links`.
	Cursor string `json:"cursor,omitempty"`
}

func NewListSharedLinksArg() *ListSharedLinksArg {
	s := new(ListSharedLinksArg)
	return s
}

type ListSharedLinksError struct {
	Tag  string             `json:".tag"`
	Path *files.LookupError `json:"path,omitempty"`
}

func (u *ListSharedLinksError) UnmarshalJSON(body []byte) error {
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

type ListSharedLinksResult struct {
	// Shared links applicable to the path argument.
	Links []*SharedLinkMetadata `json:"links"`
	// Is true if there are additional shared links that have not been returned
	// yet. Pass the cursor into :route:`list_shared_links` to retrieve them.
	HasMore bool `json:"has_more"`
	// Pass the cursor into :route:`list_shared_links` to obtain the additional
	// links. Cursor is returned only if no path is given or the path is empty.
	Cursor string `json:"cursor,omitempty"`
}

func NewListSharedLinksResult() *ListSharedLinksResult {
	s := new(ListSharedLinksResult)
	return s
}

// Actions that may be taken on members of a shared folder.
type MemberAction struct {
	Tag string `json:".tag"`
}

// Whether the user is allowed to take the action on the associated member.
type MemberPermission struct {
	// The action that the user may wish to take on the member.
	Action *MemberAction `json:"action"`
	// True if the user is allowed to take the action.
	Allow bool `json:"allow"`
	// The reason why the user is denied the permission. Not present if the action
	// is allowed
	Reason *PermissionDeniedReason `json:"reason,omitempty"`
}

func NewMemberPermission() *MemberPermission {
	s := new(MemberPermission)
	return s
}

// Policy governing who can be a member of a shared folder. Only applicable to
// folders owned by a user on a team.
type MemberPolicy struct {
	Tag string `json:".tag"`
}

// Includes different ways to identify a member of a shared folder.
type MemberSelector struct {
	Tag string `json:".tag"`
	// Dropbox account, team member, or group ID of member.
	DropboxId string `json:"dropbox_id,omitempty"`
	// E-mail address of member.
	Email string `json:"email,omitempty"`
}

func (u *MemberSelector) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// Dropbox account, team member, or group ID of member.
		DropboxId json.RawMessage `json:"dropbox_id"`
		// E-mail address of member.
		Email json.RawMessage `json:"email"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "dropbox_id":
		{
			if len(w.DropboxId) == 0 {
				break
			}
			if err := json.Unmarshal(w.DropboxId, &u.DropboxId); err != nil {
				return err
			}
		}
	case "email":
		{
			if len(w.Email) == 0 {
				break
			}
			if err := json.Unmarshal(w.Email, &u.Email); err != nil {
				return err
			}
		}
	}
	return nil
}

type ModifySharedLinkSettingsArgs struct {
	// URL of the shared link to change its settings
	Url string `json:"url"`
	// Set of settings for the shared link.
	Settings *SharedLinkSettings `json:"settings"`
}

func NewModifySharedLinkSettingsArgs() *ModifySharedLinkSettingsArgs {
	s := new(ModifySharedLinkSettingsArgs)
	return s
}

type ModifySharedLinkSettingsError struct {
	Tag string `json:".tag"`
	// There is an error with the given settings
	SettingsError *SharedLinkSettingsError `json:"settings_error,omitempty"`
}

func (u *ModifySharedLinkSettingsError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// There is an error with the given settings
		SettingsError json.RawMessage `json:"settings_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "settings_error":
		{
			if len(w.SettingsError) == 0 {
				break
			}
			if err := json.Unmarshal(w.SettingsError, &u.SettingsError); err != nil {
				return err
			}
		}
	}
	return nil
}

type MountFolderArg struct {
	// The ID of the shared folder to mount.
	SharedFolderId string `json:"shared_folder_id"`
}

func NewMountFolderArg() *MountFolderArg {
	s := new(MountFolderArg)
	return s
}

type MountFolderError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *MountFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

// Metadata for a path-based shared link.
type PathLinkMetadata struct {
	// URL of the shared link.
	Url string `json:"url"`
	// Who can access the link.
	Visibility *Visibility `json:"visibility"`
	// Path in user's Dropbox.
	Path string `json:"path"`
	// Expiration time, if set. By default the link won't expire.
	Expires time.Time `json:"expires,omitempty"`
}

func NewPathLinkMetadata() *PathLinkMetadata {
	s := new(PathLinkMetadata)
	return s
}

// Flag to indicate pending upload default (for linking to not-yet-existing
// paths).
type PendingUploadMode struct {
	Tag string `json:".tag"`
}

// Possible reasons the user is denied a permission.
type PermissionDeniedReason struct {
	Tag string `json:".tag"`
}

type RelinquishFolderMembershipArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
}

func NewRelinquishFolderMembershipArg() *RelinquishFolderMembershipArg {
	s := new(RelinquishFolderMembershipArg)
	return s
}

type RelinquishFolderMembershipError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *RelinquishFolderMembershipError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

type RemoveFolderMemberArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// The member to remove from the folder.
	Member *MemberSelector `json:"member"`
	// If true, the removed user will keep their copy of the folder after it's
	// unshared, assuming it was mounted. Otherwise, it will be removed from their
	// Dropbox. Also, this must be set to false when kicking a group.
	LeaveACopy bool `json:"leave_a_copy"`
}

func NewRemoveFolderMemberArg() *RemoveFolderMemberArg {
	s := new(RemoveFolderMemberArg)
	return s
}

type RemoveFolderMemberError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
	MemberError *SharedFolderMemberError `json:"member_error,omitempty"`
}

func (u *RemoveFolderMemberError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
		MemberError json.RawMessage `json:"member_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	case "member_error":
		{
			if len(w.MemberError) == 0 {
				break
			}
			if err := json.Unmarshal(w.MemberError, &u.MemberError); err != nil {
				return err
			}
		}
	}
	return nil
}

// The access permission that can be requested by the caller for the shared
// link. Note that the final resolved visibility of the shared link takes into
// account other aspects, such as team and shared folder settings. Check the
// :type:`ResolvedVisibility` for more info on the possible resolved visibility
// values of shared links.
type RequestedVisibility struct {
	Tag string `json:".tag"`
}

// The actual access permissions values of shared links after taking into
// account user preferences and the team and shared folder settings. Check the
// :type:`RequestedVisibility` for more info on the possible visibility values
// that can be set by the shared link's owner.
type ResolvedVisibility struct {
	Tag string `json:".tag"`
}

type RevokeSharedLinkArg struct {
	// URL of the shared link.
	Url string `json:"url"`
}

func NewRevokeSharedLinkArg() *RevokeSharedLinkArg {
	s := new(RevokeSharedLinkArg)
	return s
}

type RevokeSharedLinkError struct {
	Tag string `json:".tag"`
}

type ShareFolderArg struct {
	// The path to the folder to share. If it does not exist, then a new one is
	// created.
	Path string `json:"path"`
	// Who can be a member of this shared folder.
	MemberPolicy *MemberPolicy `json:"member_policy"`
	// Who can add and remove members of this shared folder.
	AclUpdatePolicy *AclUpdatePolicy `json:"acl_update_policy"`
	// The policy to apply to shared links created for content inside this shared
	// folder.
	SharedLinkPolicy *SharedLinkPolicy `json:"shared_link_policy"`
	// Whether to force the share to happen asynchronously.
	ForceAsync bool `json:"force_async"`
}

func NewShareFolderArg() *ShareFolderArg {
	s := new(ShareFolderArg)
	s.MemberPolicy = &MemberPolicy{Tag: "anyone"}
	s.AclUpdatePolicy = &AclUpdatePolicy{Tag: "owner"}
	s.SharedLinkPolicy = &SharedLinkPolicy{Tag: "anyone"}
	s.ForceAsync = false
	return s
}

type ShareFolderError struct {
	Tag string `json:".tag"`
	// :field:`ShareFolderArg.path` is invalid.
	BadPath *SharePathError `json:"bad_path,omitempty"`
}

func (u *ShareFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// :field:`ShareFolderArg.path` is invalid.
		BadPath json.RawMessage `json:"bad_path"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "bad_path":
		{
			if len(w.BadPath) == 0 {
				break
			}
			if err := json.Unmarshal(w.BadPath, &u.BadPath); err != nil {
				return err
			}
		}
	}
	return nil
}

type ShareFolderJobStatus struct {
	Tag string `json:".tag"`
	// The share job has finished. The value is the metadata for the folder.
	Complete *SharedFolderMetadata `json:"complete,omitempty"`
	Failed   *ShareFolderError     `json:"failed,omitempty"`
}

func (u *ShareFolderJobStatus) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The share job has finished. The value is the metadata for the folder.
		Complete json.RawMessage `json:"complete"`
		Failed   json.RawMessage `json:"failed"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "complete":
		{
			if err := json.Unmarshal(body, &u.Complete); err != nil {
				return err
			}
		}
	case "failed":
		{
			if len(w.Failed) == 0 {
				break
			}
			if err := json.Unmarshal(w.Failed, &u.Failed); err != nil {
				return err
			}
		}
	}
	return nil
}

type ShareFolderLaunch struct {
	Tag      string                `json:".tag"`
	Complete *SharedFolderMetadata `json:"complete,omitempty"`
}

func (u *ShareFolderLaunch) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag      string          `json:".tag"`
		Complete json.RawMessage `json:"complete"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "complete":
		{
			if err := json.Unmarshal(body, &u.Complete); err != nil {
				return err
			}
		}
	}
	return nil
}

type SharePathError struct {
	Tag string `json:".tag"`
}

// There is an error accessing the shared folder.
type SharedFolderAccessError struct {
	Tag string `json:".tag"`
}

type SharedFolderMemberError struct {
	Tag string `json:".tag"`
}

// Shared folder user and group membership.
type SharedFolderMembers struct {
	// The list of user members of the shared folder.
	Users []*UserMembershipInfo `json:"users"`
	// The list of group members of the shared folder.
	Groups []*GroupMembershipInfo `json:"groups"`
	// The list of invited members of the shared folder. This list will not include
	// invitees that have already accepted or declined to join the shared folder.
	Invitees []*InviteeMembershipInfo `json:"invitees"`
	// Present if there are additional shared folder members that have not been
	// returned yet. Pass the cursor into :route:`list_folder_members/continue` to
	// list additional members.
	Cursor string `json:"cursor,omitempty"`
}

func NewSharedFolderMembers() *SharedFolderMembers {
	s := new(SharedFolderMembers)
	return s
}

// The metadata which includes basic information about the shared folder.
type SharedFolderMetadata struct {
	// The name of the this shared folder.
	Name string `json:"name"`
	// The ID of the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// The current user's access level for this shared folder.
	AccessType *AccessLevel `json:"access_type"`
	// Whether this folder is a :link:`team folder
	// https://www.dropbox.com/en/help/986`.
	IsTeamFolder bool `json:"is_team_folder"`
	// Policies governing this shared folder.
	Policy *FolderPolicy `json:"policy"`
	// The lower-cased full path of this shared folder. Absent for unmounted
	// folders.
	PathLower string `json:"path_lower,omitempty"`
	// Actions the current user may perform on the folder and its contents. The set
	// of permissions corresponds to the MemberActions in the request.
	Permissions []*FolderPermission `json:"permissions,omitempty"`
}

func NewSharedFolderMetadata() *SharedFolderMetadata {
	s := new(SharedFolderMetadata)
	return s
}

type SharedLinkAccessFailureReason struct {
	Tag string `json:".tag"`
}

// Policy governing who can view shared links.
type SharedLinkPolicy struct {
	Tag string `json:".tag"`
}

type SharedLinkSettings struct {
	// The requested access for this shared link.
	RequestedVisibility *RequestedVisibility `json:"requested_visibility,omitempty"`
	// If :field:`requested_visibility` is :field:`RequestedVisibility.password`
	// this is needed to specify the password to access the link.
	LinkPassword string `json:"link_password,omitempty"`
	// Expiration time of the shared link. By default the link won't expire.
	Expires time.Time `json:"expires,omitempty"`
}

func NewSharedLinkSettings() *SharedLinkSettings {
	s := new(SharedLinkSettings)
	return s
}

type SharedLinkSettingsError struct {
	Tag string `json:".tag"`
}

// Information about a team member.
type TeamMemberInfo struct {
	// Information about the member's team
	TeamInfo *users.Team `json:"team_info"`
	// The display name of the user.
	DisplayName string `json:"display_name"`
	// ID of user as a member of a team. This field will only be present if the
	// member is in the same team as current user.
	MemberId string `json:"member_id,omitempty"`
}

func NewTeamMemberInfo() *TeamMemberInfo {
	s := new(TeamMemberInfo)
	return s
}

type TransferFolderArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// A account or team member ID to transfer ownership to.
	ToDropboxId string `json:"to_dropbox_id"`
}

func NewTransferFolderArg() *TransferFolderArg {
	s := new(TransferFolderArg)
	return s
}

type TransferFolderError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *TransferFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

type UnmountFolderArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
}

func NewUnmountFolderArg() *UnmountFolderArg {
	s := new(UnmountFolderArg)
	return s
}

type UnmountFolderError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *UnmountFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

type UnshareFolderArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// If true, members of this shared folder will get a copy of this folder after
	// it's unshared. Otherwise, it will be removed from their Dropbox. The current
	// user, who is an owner, will always retain their copy.
	LeaveACopy bool `json:"leave_a_copy"`
}

func NewUnshareFolderArg() *UnshareFolderArg {
	s := new(UnshareFolderArg)
	return s
}

type UnshareFolderError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *UnshareFolderError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

type UpdateFolderMemberArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// The member of the shared folder to update.  Only the
	// :field:`MemberSelector.dropbox_id` may be set at this time.
	Member *MemberSelector `json:"member"`
	// The new access level for :field:`member`. :field:`AccessLevel.owner` is
	// disallowed.
	AccessLevel *AccessLevel `json:"access_level"`
}

func NewUpdateFolderMemberArg() *UpdateFolderMemberArg {
	s := new(UpdateFolderMemberArg)
	return s
}

type UpdateFolderMemberError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
	MemberError *SharedFolderMemberError `json:"member_error,omitempty"`
}

func (u *UpdateFolderMemberError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
		MemberError json.RawMessage `json:"member_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	case "member_error":
		{
			if len(w.MemberError) == 0 {
				break
			}
			if err := json.Unmarshal(w.MemberError, &u.MemberError); err != nil {
				return err
			}
		}
	}
	return nil
}

// If any of the policy's are unset, then they retain their current setting.
type UpdateFolderPolicyArg struct {
	// The ID for the shared folder.
	SharedFolderId string `json:"shared_folder_id"`
	// Who can be a member of this shared folder. Only set this if the current user
	// is on a team.
	MemberPolicy *MemberPolicy `json:"member_policy,omitempty"`
	// Who can add and remove members of this shared folder.
	AclUpdatePolicy *AclUpdatePolicy `json:"acl_update_policy,omitempty"`
	// The policy to apply to shared links created for content inside this shared
	// folder.
	SharedLinkPolicy *SharedLinkPolicy `json:"shared_link_policy,omitempty"`
}

func NewUpdateFolderPolicyArg() *UpdateFolderPolicyArg {
	s := new(UpdateFolderPolicyArg)
	return s
}

type UpdateFolderPolicyError struct {
	Tag         string                   `json:".tag"`
	AccessError *SharedFolderAccessError `json:"access_error,omitempty"`
}

func (u *UpdateFolderPolicyError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag         string          `json:".tag"`
		AccessError json.RawMessage `json:"access_error"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "access_error":
		{
			if len(w.AccessError) == 0 {
				break
			}
			if err := json.Unmarshal(w.AccessError, &u.AccessError); err != nil {
				return err
			}
		}
	}
	return nil
}

// Basic information about a user. Use :route:`users.get_account` and
// :route:`users.get_account_batch`` to obtain more detailed information.
type UserInfo struct {
	// The account ID of the user.
	AccountId string `json:"account_id"`
	// If the user is in the same team as current user.
	SameTeam bool `json:"same_team"`
	// The team member ID of the shared folder member. Only present if
	// :field:`same_team` is true.
	TeamMemberId string `json:"team_member_id,omitempty"`
}

func NewUserInfo() *UserInfo {
	s := new(UserInfo)
	return s
}

// The information about a user member of the shared folder.
type UserMembershipInfo struct {
	// The access type for this member.
	AccessType *AccessLevel `json:"access_type"`
	// The account information for the membership user.
	User *UserInfo `json:"user"`
	// The permissions that requesting user has on this member. The set of
	// permissions corresponds to the MemberActions in the request.
	Permissions []*MemberPermission `json:"permissions,omitempty"`
}

func NewUserMembershipInfo() *UserMembershipInfo {
	s := new(UserMembershipInfo)
	return s
}

// Who can access a shared link. The most open visibility is :field:`public`.
// The default depends on many aspects, such as team and user preferences and
// shared folder settings.
type Visibility struct {
	Tag string `json:".tag"`
}

type Sharing interface {
	// Allows an owner or editor (if the ACL update policy allows) of a shared
	// folder to add another member. For the new member to get access to all the
	// functionality for this folder, you will need to call :route:`mount_folder`
	// on their behalf. Apps must have full Dropbox access to use this endpoint.
	// Warning: This endpoint is in beta and is subject to minor but possibly
	// backwards-incompatible changes.
	AddFolderMember(arg *AddFolderMemberArg) (err error)
	// Returns the status of an asynchronous job. Apps must have full Dropbox
	// access to use this endpoint. Warning: This endpoint is in beta and is
	// subject to minor but possibly backwards-incompatible changes.
	CheckJobStatus(arg *async.PollArg) (res *JobStatus, err error)
	// Returns the status of an asynchronous job for sharing a folder. Apps must
	// have full Dropbox access to use this endpoint. Warning: This endpoint is in
	// beta and is subject to minor but possibly backwards-incompatible changes.
	CheckShareJobStatus(arg *async.PollArg) (res *ShareFolderJobStatus, err error)
	// Create a shared link. If a shared link already exists for the given path,
	// that link is returned. Note that in the returned :type:`PathLinkMetadata`,
	// the :field:`PathLinkMetadata.url` field is the shortened URL if
	// :field:`CreateSharedLinkArg.short_url` argument is set to :val:`true`.
	// Previously, it was technically possible to break a shared link by moving or
	// renaming the corresponding file or folder. In the future, this will no
	// longer be the case, so your app shouldn't rely on this behavior. Instead, if
	// your app needs to revoke a shared link, use :route:`revoke_shared_link`.
	CreateSharedLink(arg *CreateSharedLinkArg) (res *PathLinkMetadata, err error)
	// Create a shared link with custom settings. If no settings are given then the
	// default visibility is :field:`RequestedVisibility.public` (The resolved
	// visibility, though, may depend on other aspects such as team and shared
	// folder settings).
	CreateSharedLinkWithSettings(arg *CreateSharedLinkWithSettingsArg) (res *SharedLinkMetadata, err error)
	// Returns shared folder metadata by its folder ID. Apps must have full Dropbox
	// access to use this endpoint. Warning: This endpoint is in beta and is
	// subject to minor but possibly backwards-incompatible changes.
	GetFolderMetadata(arg *GetMetadataArgs) (res *SharedFolderMetadata, err error)
	// Download the shared link's file from a user's Dropbox.
	GetSharedLinkFile(arg *GetSharedLinkMetadataArg) (res *SharedLinkMetadata, content io.ReadCloser, err error)
	// Get the shared link's metadata.
	GetSharedLinkMetadata(arg *GetSharedLinkMetadataArg) (res *SharedLinkMetadata, err error)
	// Returns a list of :type:`LinkMetadata` objects for this user, including
	// collection links. If no path is given or the path is empty, returns a list
	// of all shared links for the current user, including collection links. If a
	// non-empty path is given, returns a list of all shared links that allow
	// access to the given path.  Collection links are never returned in this case.
	// Note that the url field in the response is never the shortened URL.
	GetSharedLinks(arg *GetSharedLinksArg) (res *GetSharedLinksResult, err error)
	// Returns shared folder membership by its folder ID. Apps must have full
	// Dropbox access to use this endpoint. Warning: This endpoint is in beta and
	// is subject to minor but possibly backwards-incompatible changes.
	ListFolderMembers(arg *ListFolderMembersArgs) (res *SharedFolderMembers, err error)
	// Once a cursor has been retrieved from :route:`list_folder_members`, use this
	// to paginate through all shared folder members. Apps must have full Dropbox
	// access to use this endpoint. Warning: This endpoint is in beta and is
	// subject to minor but possibly backwards-incompatible changes.
	ListFolderMembersContinue(arg *ListFolderMembersContinueArg) (res *SharedFolderMembers, err error)
	// Return the list of all shared folders the current user has access to. Apps
	// must have full Dropbox access to use this endpoint. Warning: This endpoint
	// is in beta and is subject to minor but possibly backwards-incompatible
	// changes.
	ListFolders() (res *ListFoldersResult, err error)
	// Once a cursor has been retrieved from :route:`list_folders`, use this to
	// paginate through all shared folders. Apps must have full Dropbox access to
	// use this endpoint. Warning: This endpoint is in beta and is subject to minor
	// but possibly backwards-incompatible changes.
	ListFoldersContinue(arg *ListFoldersContinueArg) (res *ListFoldersResult, err error)
	// List shared links of this user. If no path is given or the path is empty,
	// returns a list of all shared links for the current user. If a non-empty path
	// is given, returns a list of all shared links that allow access to the given
	// path - direct links to the given path and links to parent folders of the
	// given path.
	ListSharedLinks(arg *ListSharedLinksArg) (res *ListSharedLinksResult, err error)
	// Modify the shared link's settings. If the requested visibility conflict with
	// the shared links policy of the team or the shared folder (in case the linked
	// file is part of a shared folder) then the
	// :field:`LinkPermissions.resolved_visibility` of the returned
	// :type:`SharedLinkMetadata` will reflect the actual visibility of the shared
	// link and the :field:`LinkPermissions.requested_visibility` will reflect the
	// requested visibility.
	ModifySharedLinkSettings(arg *ModifySharedLinkSettingsArgs) (res *SharedLinkMetadata, err error)
	// The current user mounts the designated folder. Mount a shared folder for a
	// user after they have been added as a member. Once mounted, the shared folder
	// will appear in their Dropbox. Apps must have full Dropbox access to use this
	// endpoint. Warning: This endpoint is in beta and is subject to minor but
	// possibly backwards-incompatible changes.
	MountFolder(arg *MountFolderArg) (res *SharedFolderMetadata, err error)
	// The current user relinquishes their membership in the designated shared
	// folder and will no longer have access to the folder.  A folder owner cannot
	// relinquish membership in their own folder. Apps must have full Dropbox
	// access to use this endpoint. Warning: This endpoint is in beta and is
	// subject to minor but possibly backwards-incompatible changes.
	RelinquishFolderMembership(arg *RelinquishFolderMembershipArg) (err error)
	// Allows an owner or editor (if the ACL update policy allows) of a shared
	// folder to remove another member. Apps must have full Dropbox access to use
	// this endpoint. Warning: This endpoint is in beta and is subject to minor but
	// possibly backwards-incompatible changes.
	RemoveFolderMember(arg *RemoveFolderMemberArg) (res *async.LaunchEmptyResult, err error)
	// Revoke a shared link. Note that even after revoking a shared link to a file,
	// the file may be accessible if there are shared links leading to any of the
	// file parent folders. To list all shared links that enable access to a
	// specific file, you can use the :route:`list_shared_links` with the file as
	// the :field:`ListSharedLinksArg.path` argument.
	RevokeSharedLink(arg *RevokeSharedLinkArg) (err error)
	// Share a folder with collaborators. Most sharing will be completed
	// synchronously. Large folders will be completed asynchronously. To make
	// testing the async case repeatable, set `ShareFolderArg.force_async`. If a
	// :field:`ShareFolderLaunch.async_job_id` is returned, you'll need to call
	// :route:`check_share_job_status` until the action completes to get the
	// metadata for the folder. Apps must have full Dropbox access to use this
	// endpoint. Warning: This endpoint is in beta and is subject to minor but
	// possibly backwards-incompatible changes.
	ShareFolder(arg *ShareFolderArg) (res *ShareFolderLaunch, err error)
	// Transfer ownership of a shared folder to a member of the shared folder. Apps
	// must have full Dropbox access to use this endpoint. Warning: This endpoint
	// is in beta and is subject to minor but possibly backwards-incompatible
	// changes.
	TransferFolder(arg *TransferFolderArg) (err error)
	// The current user unmounts the designated folder. They can re-mount the
	// folder at a later time using :route:`mount_folder`. Apps must have full
	// Dropbox access to use this endpoint. Warning: This endpoint is in beta and
	// is subject to minor but possibly backwards-incompatible changes.
	UnmountFolder(arg *UnmountFolderArg) (err error)
	// Allows a shared folder owner to unshare the folder. You'll need to call
	// :route:`check_job_status` to determine if the action has completed
	// successfully. Apps must have full Dropbox access to use this endpoint.
	// Warning: This endpoint is in beta and is subject to minor but possibly
	// backwards-incompatible changes.
	UnshareFolder(arg *UnshareFolderArg) (res *async.LaunchEmptyResult, err error)
	// Allows an owner or editor of a shared folder to update another member's
	// permissions. Apps must have full Dropbox access to use this endpoint.
	// Warning: This endpoint is in beta and is subject to minor but possibly
	// backwards-incompatible changes.
	UpdateFolderMember(arg *UpdateFolderMemberArg) (err error)
	// Update the sharing policies for a shared folder. Apps must have full Dropbox
	// access to use this endpoint. Warning: This endpoint is in beta and is
	// subject to minor but possibly backwards-incompatible changes.
	UpdateFolderPolicy(arg *UpdateFolderPolicyArg) (res *SharedFolderMetadata, err error)
}
