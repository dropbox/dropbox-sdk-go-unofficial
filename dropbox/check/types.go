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

// Package check : has no documentation (yet)
package check

// EchoArg : EchoArg contains the arguments to be sent to the Dropbox servers.
type EchoArg struct {
	// Query : The string that you'd like to be echoed back to you.
	Query string `json:"query"`
}

// NewEchoArg returns a new EchoArg instance
func NewEchoArg() *EchoArg {
	s := new(EchoArg)
	s.Query = ""
	return s
}

// EchoResult : EchoResult contains the result returned from the Dropbox
// servers.
type EchoResult struct {
	// Result : If everything worked correctly, this would be the same as query.
	Result string `json:"result"`
}

// NewEchoResult returns a new EchoResult instance
func NewEchoResult() *EchoResult {
	s := new(EchoResult)
	s.Result = ""
	return s
}
