# Dropbox Go SDK Generator

This directory contains the [Stone](https://github.com/dropbox/stone) code generators
used to programmatically generate the [Dropbox Go SDK](https://github.com/dropbox/dropbox-sdk-go).

## Requirements

  * While not a hard requirement, this repo currently assumes `python3` in the path.
  * Assumes you have already installed [Stone](https://github.com/dropbox/stone)
  * Requires [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) to fix up imports in the auto-generated code

## Basic Setup

  . Clone this repo
  . Run `git submodule init` followed by `git submodule update`
  . Run `./generate-sdk.sh` to generate code under `../dropbox`

## Generated Code

### Basic Types

Here is how Stone [basic types](https://github.com/dropbox/stone/blob/master/doc/lang_ref.rst#basic-types) map to Go types:

Stone Type | Go Type
---------- | -------
Int32/Int64/UInt32/UInt64 | int32/int64/uint32/uint64
Float32/Float64 | float32/float64
Boolean | bool
String | string
Timestamp | time.Time
Void | struct{}

### Structs

Stone [structs](https://github.com/dropbox/stone/blob/master/doc/lang_ref.rst#struct) are represented as Go [structs](https://gobyexample.com/structs) in a relatively straight-forward manner. Each struct member is exported and also gets assigned the correct json tag. The latter is used for serializing requests and deserializing responses. Non-primitive types are represented as pointers to the corresponding type.

```
struct Account
    "The amount of detail revealed about an account depends on the user
    being queried and the user making the query."

    account_id AccountId
        "The user's unique Dropbox ID."
    name Name
        "Details of a user's name."
```

```go
// The amount of detail revealed about an account depends on the user being
// queried and the user making the query.
type Account struct {
	// The user's unique Dropbox ID.
	AccountId string `json:"account_id"`
	// Details of a user's name.
	Name *Name `json:"name"`
}
```

### Unions

Stone https://github.com/dropbox/stone/blob/master/doc/lang_ref.rst#union[unions] are bit more complex as Go doesn't have native support for union types (tagged or otherwise). We declare a union as a Go struct with all the possible fields as pointer types, and then use the tag value to populate the correct field during deserialization. This necessitates the use of an intermedia wrapper struct for the deserialization to work correctly, see below for a concrete example.

```
union SpaceAllocation
    "Space is allocated differently based on the type of account."

    individual IndividualSpaceAllocation
        "The user's space allocation applies only to their individual account."
    team TeamSpaceAllocation
        "The user shares space with other members of their team."
```

```go
// Space is allocated differently based on the type of account.
type SpaceAllocation struct { // <1>
  Tag string `json:".tag"` // <2>
  // The user's space allocation applies only to their individual account.
  Individual *IndividualSpaceAllocation `json:"individual,omitempty"` // <3>
  // The user shares space with other members of their team.
  Team *TeamSpaceAllocation `json:"team,omitempty"`
}

func (u *SpaceAllocation) UnmarshalJSON(body []byte) error { // <4>
  type wrap struct { // <5>
    Tag string `json:".tag"`
    // The user's space allocation applies only to their individual account.
    Individual json.RawMessage `json:"individual"` // <6>
    // The user shares space with other members of their team.
    Team json.RawMessage `json:"team"`
  }
  var w wrap
  if err := json.Unmarshal(body, &w); err != nil { // <7>
    return err
  }
  u.Tag = w.Tag
  switch w.Tag {
  case "individual":
    {
      if err := json.Unmarshal(body, &u.Individual); err != nil { // <8>
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
```
<1> A babel union is represented as Go struct
<2> The tag value is used to determine which field is actually set
<3> Possible values are represented as pointer types. Note the `omitempty` in the JSON tag; this is so values not set are automatically elided during serialization.
<4> We generate a custom `UnmarshalJSON` method for union types
<5> An intermedia wrapper struct is used to help with deserialization
<6> Note that members of the wrapper struct are of type `RawMessage` so we can capture the body without deserializing it right away
<7> When we deserialize response into the wrapper struct, it should get the tag value and the raw JSON as part of the members.
<8> We then use the tag value to deserialize the `RawMessage` into the appropriate member of the union type

### Struct with Enumerated Subtypes

Per the https://github.com/dropbox/stone/blob/master/doc/lang_ref.rst#struct-polymorphism[spec], structs with enumerated subtypes are a mechanism of inheritance:

> If a struct enumerates its subtypes, an instance of any subtype will satisfy the type constraint.

So to represent structs with enumerated subtypes in Go, we use a strategy similar to unions. The _base_ struct (the one that defines the subtypes) is represented like we represent a union above. The _subtype_ is represented as a struct that essentially duplicates all fields of the base type and includes fields specific to that subtype. Here's an example:

```
struct Metadata
    "Metadata for a file or folder."

    union
        file FileMetadata
        folder FolderMetadata
        deleted DeletedMetadata  # Used by list_folder* and search

    name String #<1>
        "The last component of the path (including extension).
        This never contains a slash."

...
struct FileMetadata extends Metadata
    id Id? #<2>
   ...
```
<1> Field common to all subtypes
<2> Field specific to `FileMetadata`


```go
// Metadata for a file or folder.
type Metadata struct { // <1>
  Tag     string           `json:".tag"`
  File    *FileMetadata    `json:"file,omitempty"`
  Folder  *FolderMetadata  `json:"folder,omitempty"`
  Deleted *DeletedMetadata `json:"deleted,omitempty"`
}
...

type FileMetadata struct {
  // The last component of the path (including extension). This never contains a
  // slash.
  Name string `json:"name"` // <2>
  ...
  Id string `json:"id,omitempty"` // <3>
}
```
<1> Subtype is represented like we represent unions as described above
<2> Common fields are duplicated in subtypes
<3> Subtype specific fields are included as usual in structs
