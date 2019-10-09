package sharing_test

import (
	"fmt"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"

	"encoding/json"
	"testing"
)

func TestSharedLinkAlreadyExistsMetadata(t *testing.T) {
	slaem := sharing.SharedLinkAlreadyExistsMetadata{}
	slaem.Tag = sharing.SharedLinkAlreadyExistsMetadataMetadata
	slm := sharing.NewFileLinkMetadata("http://test.com", "test", nil, time.Now(), time.Now(), "rev", 10)
	slaem.Metadata = slm
	// see if serializing works: However, the SDK doesn't correctly serialize strunions
	data, err := json.Marshal(slaem)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("testt", string(data))
	data = []byte(`{".tag":"metadata","metadata":{".tag":"file","url":"http://test.com","name":"test","expires":"0001-01-01T00:00:00Z","link_permissions":null,"client_modified":"2019-10-08T20:52:10.418397-07:00","server_modified":"2019-10-08T20:52:10.418397-07:00","rev":"rev","size":10}}`)
	//data = []byte(`{".tag":"metadata","metadata":{".tag":"file","url":"http://test.com","name":"test","expires":"0001-01-01T00:00:00Z","link_permissions":null,"client_modified":"2019-10-08T22:41:30.937397-07:00","server_modified":"2019-10-08T22:41:30.937397-07:00","rev":"rev","size":10}}`)
	var slaem2 sharing.SharedLinkAlreadyExistsMetadata
	err = json.Unmarshal(data, &slaem2)
	if err != nil {
		t.Fatal(err)
	}

	if slaem2.Metadata == nil {
		t.Fatal("expected non-nil metadata")
	}

	slm2 := slaem2.Metadata.(*sharing.FileLinkMetadata)
	if slm.Name != slm2.Name {
		t.Fatalf("%v != %v", slm.Name, slm2.Name)
	}
}

func TestInviteeInfo(t *testing.T) {
	ii := sharing.InviteeInfo{Email: "hi"}
	ii.Tag = sharing.InviteeInfoEmail
	data, err := json.Marshal(ii)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("data", string(data))
	var ii2 sharing.InviteeInfo
	err = json.Unmarshal(data, &ii2)
	if err != nil {
		t.Fatal(err)
	}
	if ii.Email != ii2.Email {
		t.Fatalf("%v != %v", ii.Email, ii2.Email)
	}
}
