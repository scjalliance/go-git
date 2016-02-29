package filesystem

import (
	"os"

	"gopkg.in/src-d/go-git.v3/core"
	"gopkg.in/src-d/go-git.v3/formats/objfile"
)

// Writer writes and encodes data in compressed objfile format to the provided
// os.File.
//
// Writer implements io.WriteCloser. Close should be called when finished with
// the Writer.
//
// When Writer is closed it will also close the underlying os.File.
type Writer struct {
	f *os.File
	objfile.Writer
}

// NewWriter returns a new Writer writing to f.
//
// Calling NewWriter causes it to immediately write header data to f
// containing size and type information. Any errors encountered in that
// process will be returned in err.
//
// The returned Writer implements io.WriteCloser. Close should be called when
// finished with the Writer. Close will also close the underlying os.File.
//
// If NewWriter returns an error it will close the underlying os.File.
func NewWriter(f *os.File, t core.ObjectType, size int64) (w *Writer, err error) {
	var ofw *objfile.Writer
	ofw, err = objfile.NewWriter(f, t, size)
	if err != nil {
		defer f.Close()
		return
	}

	w = &Writer{
		f:      f,
		Writer: *ofw,
	}

	return
}

// Close releases any resources consumed by the writer.
//
// Calling Close will close the underlying os.File originally passed in to
// NewWriter.
func (w *Writer) Close() (err error) {
	err = w.Writer.Close()
	if err != nil {
		w.f.Close()
		return
	}

	return w.f.Close()

	// 3. Rename file if it was temporary
	// Delete or rename file if the hash doesn't match the expected value?
	/*
		if newpath := w.s.Path(w.h); path != newpath {
			if err := os.Rename(path, s.Path(w.h)); err != nil {
				defer os.Remove(path)
				return err
			}
		}
	*/

	// TODO: Watch for errors on Write and delete the temp file if the write
	//       wasn't successful?
}
