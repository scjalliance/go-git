package filesystem

import (
	"errors"
	"io"
)

var (
	// ErrPersisted is returned when an attempt is made to obtain a Writer for
	// an object that has already been written to the file system.
	ErrPersisted = errors.New("filesystem: object cannot be written to after it has been persisted")

	// ErrNotPersisted is returned when an attempt is made to obtain a Reader for
	// an object that hasn't yet been written to the file system.
	ErrNotPersisted = errors.New("filesystem: object cannot be read from until it has been persisted")

	// ErrZeroHash is returned when a hash is epected to be non-zero.
	ErrZeroHash = errors.New("filesystem: object hash is zero")
)

// checkClose is used with defer to close the given io.Closer and check its
// returned error value. If Close returns an error and the given *error
// is not nil, *error is set to the error returned by Close.
//
// checkClose is typically used with named return values like so:
//
//   func do(obj *Object) (err error) {
//     w, err := obj.Writer()
//     if err != nil {
//       return nil
//     }
//     defer checkClose(w, &err)
//     // work with w
//   }
func checkClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
