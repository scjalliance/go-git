package filesystem

import (
	"os"

	"gopkg.in/src-d/go-git.v3/formats/objfile"
)

// Reader reads and decodes data in compressed objfile format from the provided
// os.File.
//
// Reader implements io.ReadCloser. Close should be called when finished with
// the reader.
//
// Reader will close the underlying os.File when it is closed.
type Reader struct {
	f *os.File
	objfile.Reader
}

// NewReader returns a new Reader reading from f.
//
// Calling NewReader causes the it to immediately read in the header data
// from f containing size and type information. Any errors encountered in that
// process will be returned in err.
//
// The returned Reader implements io.ReadCloser. Close should be called when
// finished with the Reader. Close will also close the underlying os.File.
func NewReader(f *os.File) (r *Reader, err error) {
	ofr, err := objfile.NewReader(f)
	if err != nil {
		f.Close() // Close file on init failure
		return
	}

	r = &Reader{
		f:      f,
		Reader: *ofr,
	}

	return
}

// Close releases any resources consumed by the reader.
//
// Calling Close will close the underlying os.File originally passed in to
// NewReader.
func (r *Reader) Close() (err error) {
	err = r.Reader.Close()
	if err != nil {
		r.f.Close() // Close file even if reader close fails
		return
	}
	// TODO: Somehow assert that the hash matches what we expected?
	return r.f.Close()
}
