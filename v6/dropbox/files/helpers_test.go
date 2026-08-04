package files

import (
	"math"
	"testing"
)

func TestSetRangeLength(t *testing.T) {
	tests := []struct {
		name       string
		arg        *DownloadArg
		offset     int64
		length     int64
		wantHeader string
		wantErr    bool
	}{
		{
			name:       "sets range header",
			arg:        NewDownloadArg("/file.bin"),
			offset:     100,
			length:     50,
			wantHeader: "bytes=100-149",
		},
		{
			name:       "preserves existing headers",
			arg:        &DownloadArg{ExtraHeaders: map[string]string{"If-None-Match": "etag"}},
			offset:     0,
			length:     1,
			wantHeader: "bytes=0-0",
		},
		{
			name:    "nil argument",
			arg:     nil,
			offset:  0,
			length:  1,
			wantErr: true,
		},
		{
			name:    "negative offset",
			arg:     NewDownloadArg("/file.bin"),
			offset:  -1,
			length:  1,
			wantErr: true,
		},
		{
			name:    "zero length",
			arg:     NewDownloadArg("/file.bin"),
			offset:  0,
			length:  0,
			wantErr: true,
		},
		{
			name:    "negative length",
			arg:     NewDownloadArg("/file.bin"),
			offset:  0,
			length:  -1,
			wantErr: true,
		},
		{
			name:    "range overflow",
			arg:     NewDownloadArg("/file.bin"),
			offset:  math.MaxInt64,
			length:  2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetRangeLength(tt.arg, tt.offset, tt.length)
			if tt.wantErr {
				if err == nil {
					t.Fatal("SetRangeLength() expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("SetRangeLength() error = %v", err)
			}
			if got := tt.arg.ExtraHeaders["Range"]; got != tt.wantHeader {
				t.Errorf("Range header = %q, want %q", got, tt.wantHeader)
			}
		})

	}
}

func TestSetRange(t *testing.T) {
	tests := []struct {
		name       string
		arg        *DownloadArg
		offset     int64
		wantHeader string
		wantErr    bool
	}{
		{
			name:       "sets range header",
			arg:        NewDownloadArg("/file.bin"),
			offset:     100,
			wantHeader: "bytes=100-",
		},
		{
			name:       "sets range from zero",
			arg:        NewDownloadArg("/file.bin"),
			offset:     0,
			wantHeader: "bytes=0-",
		},
		{
			name: "preserves existing headers",
			arg: &DownloadArg{
				ExtraHeaders: map[string]string{
					"If-None-Match": "etag",
				},
			},
			offset:     100,
			wantHeader: "bytes=100-",
		},
		{
			name:    "nil argument",
			arg:     nil,
			offset:  0,
			wantErr: true,
		},
		{
			name:    "negative offset",
			arg:     NewDownloadArg("/file.bin"),
			offset:  -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetRange(tt.arg, tt.offset)
			if tt.wantErr {
				if err == nil {
					t.Fatal("SetRange() expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("SetRange() error = %v", err)
			}
			if got := tt.arg.ExtraHeaders["Range"]; got != tt.wantHeader {
				t.Errorf("Range header = %q, want %q", got, tt.wantHeader)
			}
			if got := tt.arg.ExtraHeaders["If-None-Match"]; tt.name == "preserves existing headers" && got != "etag" {
				t.Errorf("existing header changed: got %q", got)
			}
		})
	}
}
