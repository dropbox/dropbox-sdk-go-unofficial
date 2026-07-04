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

package contenthash

import (
	"crypto/sha256"
	"errors"
	"hash"
	"strings"
	"testing"
)

func TestNewDoesNotPreallocateBlock(t *testing.T) {
	h := New()
	if h.overall != nil {
		t.Fatal("New initialized overall hash before it was needed")
	}
	if cap(h.block) != 0 {
		t.Fatalf("New block capacity = %d, want 0", cap(h.block))
	}
}

func TestCloneSHA256PanicsOnMarshalError(t *testing.T) {
	assertPanicContains(t, "failed to marshal SHA-256 state", func() {
		cloneSHA256(failingMarshalHash{Hash: sha256.New()})
	})
}

func TestCloneSHA256PanicsOnUnmarshalError(t *testing.T) {
	assertPanicContains(t, "failed to unmarshal SHA-256 state", func() {
		cloneSHA256(invalidStateHash{Hash: sha256.New()})
	})
}

type failingMarshalHash struct {
	hash.Hash
}

func (h failingMarshalHash) MarshalBinary() ([]byte, error) {
	return nil, errors.New("marshal failed")
}

type invalidStateHash struct {
	hash.Hash
}

func (h invalidStateHash) MarshalBinary() ([]byte, error) {
	return []byte("invalid SHA-256 state"), nil
}

func assertPanicContains(t *testing.T, want string, f func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("expected panic containing %q", want)
		}
		gotText, ok := got.(string)
		if !ok {
			t.Fatalf("panic = %v, want string containing %q", got, want)
		}
		if !strings.Contains(gotText, want) {
			t.Fatalf("panic = %q, want string containing %q", gotText, want)
		}
	}()

	f()
}
