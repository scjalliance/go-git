package filesystem

import (
	"io/ioutil"
	"os"

	"gopkg.in/src-d/go-git.v3/core"
)

// Object is a filesystem-based core.Object implementation.
type Object struct {
	t    core.ObjectType
	h    core.Hash
	size int64
	path string
	w    *Writer
}

// Hash returns the object Hash.
func (o *Object) Hash() core.Hash {
	if o.w != nil {
		return o.w.Hash() // Return the hash from the writer if the object has been written to
	}
	return o.h
}

// Type returns the core.ObjectType.
func (o *Object) Type() core.ObjectType {
	//f, err := o.open()
	return o.t
}

// SetType sets the core.ObjectType.
func (o *Object) SetType(t core.ObjectType) { o.t = t }

// Size returns the size of the object.
func (o *Object) Size() int64 { return o.size }

// SetSize sets the object size.
func (o *Object) SetSize(s int64) { o.size = s }

// Reader returns a core.ObjectReader used to read the object's content.
//
// The returned reader implements io.ReadCloser. Close should be called when
// finished with the reader.
func (o *Object) Reader() (core.ObjectReader, error) {
	f, err := o.openFile()
	if err != nil {
		return nil, err
	}

	r, err := NewReader(f)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Writer returns a core.ObjectWriter used to write the object's content.
//
// The returned writer implements io.WriteCloser. Close should be called when
// finished with the writer.
func (o *Object) Writer() (core.ObjectWriter, error) {
	f, err := o.createFile() // Could return a temporary file
	if err != nil {
		return nil, err
	}

	w, err := NewWriter(f, o.t, o.size)
	if err != nil {
		return nil, err
	}

	o.w = w // retained for access to w.Hash()

	return w, nil
}

// Path returns the path of the file containing the object's data. If the object
// is new and ObjectStorage.Set has not yet been called, this will be a
// temporary file path.
func (o *Object) Path() string {
	return o.path
}

// openFile will open a os.File for writing. If o.path == "" it will return
// ErrNotPersisted.
func (o *Object) openFile() (*os.File, error) {
	if o.path == "" {
		return nil, ErrNotPersisted
	}

	return os.Open(o.path)
}

// createFile will prepare a os.File for writing. If o.path == "" it will return
// a temporary file created by os.TempFile, otherwise it will return a file
// created by os.Create.
//
// Temporary files created by this function should be renamed or removed.
func (o *Object) createFile() (*os.File, error) {
	if o.path != "" {
		return os.Create(o.path)
	}

	f, err := ioutil.TempFile("", "git-object-")
	if err == nil {
		o.path = f.Name()
	}

	return f, err
}
