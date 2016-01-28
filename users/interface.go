/* DO NOT EDIT */
/* This file was generated from users.babel */

package users

import "encoding/json"

// The amount of detail revealed about an account depends on the user being
// queried and the user making the query.
type Account struct {
	// The user's unique Dropbox ID.
	AccountId string `json:"account_id"`
	// Details of a user's name.
	Name *Name `json:"name"`
}

func NewAccount() *Account {
	s := new(Account)
	return s
}

// What type of account this user has.
type AccountType struct {
	Tag string `json:".tag"`
}

func (u *AccountType) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	}
	return nil
}

// Basic information about any account.
type BasicAccount struct {
	// The user's unique Dropbox ID.
	AccountId string `json:"account_id"`
	// Details of a user's name.
	Name *Name `json:"name"`
	// Whether this user is a teammate of the current user. If this account is the
	// current user's account, then this will be :val:`true`.
	IsTeammate bool `json:"is_teammate"`
}

func NewBasicAccount() *BasicAccount {
	s := new(BasicAccount)
	return s
}

// Detailed information about the current user's account.
type FullAccount struct {
	// The user's unique Dropbox ID.
	AccountId string `json:"account_id"`
	// Details of a user's name.
	Name *Name `json:"name"`
	// The user's e-mail address. Do not rely on this without checking the
	// :field:`email_verified` field. Even then, it's possible that the user has
	// since lost access to their e-mail.
	Email string `json:"email"`
	// Whether the user has verified their e-mail address.
	EmailVerified bool `json:"email_verified"`
	// The language that the user specified. Locale tags will be :link:`IETF
	// language tags http://en.wikipedia.org/wiki/IETF_language_tag`.
	Locale string `json:"locale"`
	// The user's :link:`referral link https://www.dropbox.com/referrals`.
	ReferralLink string `json:"referral_link"`
	// Whether the user has a personal and work account. If the current account is
	// personal, then :field:`team` will always be :val:`null`, but
	// :field:`is_paired` will indicate if a work account is linked.
	IsPaired bool `json:"is_paired"`
	// What type of account this user has.
	AccountType *AccountType `json:"account_type"`
	// The user's two-letter country code, if available. Country codes are based on
	// :link:`ISO 3166-1 http://en.wikipedia.org/wiki/ISO_3166-1`.
	Country string `json:"country,omitempty"`
	// If this account is a member of a team, information about that team.
	Team *Team `json:"team,omitempty"`
}

func NewFullAccount() *FullAccount {
	s := new(FullAccount)
	return s
}

type GetAccountArg struct {
	// A user's account identifier.
	AccountId string `json:"account_id"`
}

func NewGetAccountArg() *GetAccountArg {
	s := new(GetAccountArg)
	return s
}

type GetAccountBatchArg struct {
	// List of user account identifiers.  Should not contain any duplicate account
	// IDs.
	AccountIds []string `json:"account_ids"`
}

func NewGetAccountBatchArg() *GetAccountBatchArg {
	s := new(GetAccountBatchArg)
	return s
}

type GetAccountBatchError struct {
	Tag string `json:".tag"`
	// The value is an account ID specified in
	// :field:`GetAccountBatchArg.account_ids` that does not exist.
	NoAccount string `json:"no_account,omitempty"`
}

func (u *GetAccountBatchError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The value is an account ID specified in
		// :field:`GetAccountBatchArg.account_ids` that does not exist.
		NoAccount json.RawMessage `json:"no_account"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "no_account":
		{
			if len(w.NoAccount) == 0 {
				break
			}
			if err := json.Unmarshal(w.NoAccount, &u.NoAccount); err != nil {
				return err
			}
		}
	}
	return nil
}

type GetAccountError struct {
	Tag string `json:".tag"`
}

func (u *GetAccountError) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	}
	return nil
}

type IndividualSpaceAllocation struct {
	// The total space allocated to the user's account (bytes).
	Allocated uint64 `json:"allocated"`
}

func NewIndividualSpaceAllocation() *IndividualSpaceAllocation {
	s := new(IndividualSpaceAllocation)
	return s
}

// Representations for a person's name to assist with internationalization.
type Name struct {
	// Also known as a first name.
	GivenName string `json:"given_name"`
	// Also known as a last name or family name.
	Surname string `json:"surname"`
	// Locale-dependent name. In the US, a person's familiar name is their
	// :field:`given_name`, but elsewhere, it could be any combination of a
	// person's :field:`given_name` and :field:`surname`.
	FamiliarName string `json:"familiar_name"`
	// A name that can be used directly to represent the name of a user's Dropbox
	// account.
	DisplayName string `json:"display_name"`
}

func NewName() *Name {
	s := new(Name)
	return s
}

// Space is allocated differently based on the type of account.
type SpaceAllocation struct {
	Tag string `json:".tag"`
	// The user's space allocation applies only to their individual account.
	Individual *IndividualSpaceAllocation `json:"individual,omitempty"`
	// The user shares space with other members of their team.
	Team *TeamSpaceAllocation `json:"team,omitempty"`
}

func (u *SpaceAllocation) UnmarshalJSON(body []byte) error {
	type wrap struct {
		Tag string `json:".tag"`
		// The user's space allocation applies only to their individual account.
		Individual json.RawMessage `json:"individual"`
		// The user shares space with other members of their team.
		Team json.RawMessage `json:"team"`
	}
	var w wrap
	if err := json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch w.Tag {
	case "individual":
		{
			if err := json.Unmarshal(body, &u.Individual); err != nil {
				return err
			}
		}
	case "team":
		{
			if err := json.Unmarshal(body, &u.Team); err != nil {
				return err
			}
		}
	}
	return nil
}

// Information about a user's space usage and quota.
type SpaceUsage struct {
	// The user's total space usage (bytes).
	Used uint64 `json:"used"`
	// The user's space allocation.
	Allocation *SpaceAllocation `json:"allocation"`
}

func NewSpaceUsage() *SpaceUsage {
	s := new(SpaceUsage)
	return s
}

// Information about a team.
type Team struct {
	// The team's unique ID.
	Id string `json:"id"`
	// The name of the team.
	Name string `json:"name"`
}

func NewTeam() *Team {
	s := new(Team)
	return s
}

type TeamSpaceAllocation struct {
	// The total space currently used by the user's team (bytes).
	Used uint64 `json:"used"`
	// The total space allocated to the user's team (bytes).
	Allocated uint64 `json:"allocated"`
}

func NewTeamSpaceAllocation() *TeamSpaceAllocation {
	s := new(TeamSpaceAllocation)
	return s
}

type Users interface {
	// Get information about a user's account.
	GetAccount(arg *GetAccountArg) (res *BasicAccount, err error)
	// Get information about multiple user accounts.  At most 300 accounts may be
	// queried per request.
	GetAccountBatch(arg *GetAccountBatchArg) (res []*BasicAccount, err error)
	// Get information about the current user's account.
	GetCurrentAccount() (res *FullAccount, err error)
	// Get the space usage information for the current user's account.
	GetSpaceUsage() (res *SpaceUsage, err error)
}
