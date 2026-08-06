package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// UploadSource provides sequential access to stable upload content.
//
// The source content must remain unchanged until Upload returns.
type UploadSource interface {
	// Open opens the source from its beginning.
	//
	// Upload closes the returned reader before returning.
	Open(ctx context.Context) (io.ReadCloser, error)
}

// SizedUploadSource provides upload content with a known size.
type SizedUploadSource interface {
	UploadSource

	// Size returns the total size of the source in bytes.
	Size() int64
}

// RangedUploadSource provides repeatable random access to upload content.
type RangedUploadSource interface {
	SizedUploadSource

	// OpenRange opens an independent reader for the half-open range
	// [offset, offset+length).
	//
	// The same range may be opened repeatedly or concurrently. Each call for
	// the same range must return the same bytes until Upload returns.
	OpenRange(
		ctx context.Context,
		offset int64,
		length int64,
	) (io.ReadCloser, error)
}

// -----------------------------------------------------------------------------
// Built-in upload sources
// -----------------------------------------------------------------------------

// FileSource reads upload content from a file.
type FileSource struct {
	path string
	size int64
}

// FileUpload opens path as a ranged upload source.
//
// The file must remain unchanged until Upload returns.
func FileUpload(path string) (*FileSource, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("upload source is not a regular file: %s", path)
	}

	return &FileSource{
		path: path,
		size: info.Size(),
	}, nil
}

// Size returns the size of the file in bytes.
func (s *FileSource) Size() int64 {
	return s.size
}

// Open opens the file from its beginning.
func (s *FileSource) Open(ctx context.Context) (io.ReadCloser, error) {
	return s.OpenRange(ctx, 0, s.size)
}

// OpenRange opens an independent reader for the requested file range.
func (s *FileSource) OpenRange(
	ctx context.Context,
	offset int64,
	length int64,
) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateRange(s.size, offset, length); err != nil {
		return nil, err
	}

	file, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}

	return &sectionReadCloser{
		Reader: io.NewSectionReader(file, offset, length),
		close:  file.Close,
	}, nil
}

// BytesSource reads upload content from memory.
type BytesSource struct {
	data []byte
}

// BytesUpload returns a ranged upload source backed by data.
//
// The caller must not modify data until Upload returns.
func BytesUpload(data []byte) *BytesSource {
	return &BytesSource{
		data: data,
	}
}

// Size returns the number of bytes in the source.
func (s *BytesSource) Size() int64 {
	return int64(len(s.data))
}

// Open opens the source from its beginning.
func (s *BytesSource) Open(ctx context.Context) (io.ReadCloser, error) {
	return s.OpenRange(ctx, 0, int64(len(s.data)))
}

// OpenRange opens an independent reader for the requested byte range.
func (s *BytesSource) OpenRange(
	ctx context.Context,
	offset int64,
	length int64,
) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateRange(int64(len(s.data)), offset, length); err != nil {
		return nil, err
	}

	start := int(offset)
	end := int(offset + length)

	return io.NopCloser(bytes.NewReader(s.data[start:end])), nil
}

// ReaderSource reads upload content from a one-shot stream of unknown size.
type ReaderSource struct {
	mu     sync.Mutex
	reader io.Reader
	opened bool
}

// ReaderUpload returns a sequential upload source with an unknown size.
//
// The source may be opened only once.
func ReaderUpload(reader io.Reader) (*ReaderSource, error) {
	if reader == nil {
		return nil, errors.New("upload reader is required")
	}

	return &ReaderSource{
		reader: reader,
	}, nil
}

// Open opens the source.
func (s *ReaderSource) Open(ctx context.Context) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.opened {
		return nil, errors.New("upload source has already been opened")
	}

	s.opened = true

	if closer, ok := s.reader.(io.ReadCloser); ok {
		return closer, nil
	}

	return io.NopCloser(s.reader), nil
}

// SizedReaderSource reads upload content from a one-shot stream of known size.
type SizedReaderSource struct {
	*ReaderSource
	size int64
}

// SizedReaderUpload returns a sequential upload source with a known size.
//
// The source may be opened only once.
func SizedReaderUpload(
	reader io.Reader,
	size int64,
) (*SizedReaderSource, error) {
	if reader == nil {
		return nil, errors.New("upload reader is required")
	}
	if size < 0 {
		return nil, errors.New("upload size must not be negative")
	}

	source, err := ReaderUpload(reader)
	if err != nil {
		return nil, err
	}

	return &SizedReaderSource{
		ReaderSource: source,
		size:         size,
	}, nil
}

// Size returns the source size.
func (s *SizedReaderSource) Size() int64 {
	return s.size
}

type sectionReadCloser struct {
	io.Reader
	close func() error
}

func (r *sectionReadCloser) Close() error {
	return r.close()
}

func validateRange(size int64, offset int64, length int64) error {
	if size < 0 {
		return errors.New("source size must not be negative")
	}
	if offset < 0 {
		return errors.New("range offset must not be negative")
	}
	if length < 0 {
		return errors.New("range length must not be negative")
	}
	if offset > size || length > size-offset {
		return fmt.Errorf(
			"range [%d,%d) exceeds source size %d",
			offset,
			offset+length,
			size,
		)
	}

	return nil
}
