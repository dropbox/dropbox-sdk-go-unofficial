package riviera_test

import (
	"encoding/json"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/riviera"
)

func TestMetadataUnionIsExported(t *testing.T) {
	metadata := &riviera.MetadataUnion{
		Tagged: dropbox.Tagged{Tag: riviera.MetadataUnionExif},
		Exif:   riviera.NewApiExifMetadata(),
	}
	result := &riviera.GetMetadataResult{Metadata: metadata}
	if result.Metadata.Exif == nil {
		t.Fatal("expected EXIF metadata")
	}
}

func TestMetadataUnionUnmarshal(t *testing.T) {
	data := []byte(`{
		"metadata_type": {".tag": "metadata_type_exif"},
		"metadata": {".tag": "exif", "image_width": 640, "image_height": 480}
	}`)

	var result riviera.GetMetadataResult
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result.Metadata == nil || result.Metadata.Tag != riviera.MetadataUnionExif {
		t.Fatalf("unexpected metadata union: %#v", result.Metadata)
	}
	if result.Metadata.Exif == nil {
		t.Fatal("expected EXIF metadata")
	}
	if result.Metadata.Exif.ImageWidth != 640 || result.Metadata.Exif.ImageHeight != 480 {
		t.Fatalf("unexpected EXIF dimensions: %dx%d",
			result.Metadata.Exif.ImageWidth, result.Metadata.Exif.ImageHeight)
	}
}
