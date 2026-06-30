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

package dropbox

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDBXTimeMarshalUTC(t *testing.T) {
	// A UTC time with no nanoseconds should serialize cleanly
	ts := time.Date(2022, 2, 24, 8, 50, 46, 0, time.UTC)
	dbxTime := DBXTime(ts)

	b, err := json.Marshal(dbxTime)
	if err != nil {
		t.Fatal(err)
	}

	want := `"2022-02-24T08:50:46Z"`
	if string(b) != want {
		t.Errorf("got %s, want %s", string(b), want)
	}
}

func TestDBXTimeMarshalNonUTCWithNanos(t *testing.T) {
	// This is the exact bug from issue #109:
	// time with nanoseconds and non-UTC timezone must be converted to UTC and truncated
	loc := time.FixedZone("CST", 8*3600)
	ts := time.Date(2022, 2, 24, 16, 50, 46, 921728694, loc)
	dbxTime := DBXTime(ts)

	b, err := json.Marshal(dbxTime)
	if err != nil {
		t.Fatal(err)
	}

	// Must produce UTC, no nanos, Z suffix
	want := `"2022-02-24T08:50:46Z"`
	if string(b) != want {
		t.Errorf("got %s, want %s", string(b), want)
	}
}

func TestDBXTimeUnmarshal(t *testing.T) {
	input := `"2022-02-24T08:50:46Z"`
	var dbxTime DBXTime
	if err := json.Unmarshal([]byte(input), &dbxTime); err != nil {
		t.Fatal(err)
	}

	got := time.Time(dbxTime)
	want := time.Date(2022, 2, 24, 8, 50, 46, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDBXTimeRoundTrip(t *testing.T) {
	loc := time.FixedZone("PST", -8*3600)
	original := time.Date(2023, 6, 15, 10, 30, 0, 123456789, loc)
	dbxTime := DBXTime(original)

	b, err := json.Marshal(dbxTime)
	if err != nil {
		t.Fatal(err)
	}

	var decoded DBXTime
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}

	// After roundtrip, should be UTC-truncated
	got := time.Time(decoded)
	want := time.Date(2023, 6, 15, 18, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
