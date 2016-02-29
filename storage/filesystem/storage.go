package filesystem

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v3/core"
)

type ObjectStorage struct {
	root string
	temp string
}

// NewObjectStorage returns an ObjectStorage that will store its data in files
// within the provided root directory.
func NewObjectStorage(root string) (*ObjectStorage, error) {
	return &ObjectStorage{
		root: filepath.Dir(root),
		temp: os.TempDir(),
	}, nil
}

// New returns a new Object.
func (o *ObjectStorage) New() (core.Object, error) {
	return &Object{}, nil
}

func (o *ObjectStorage) Set(obj core.Object) (core.Hash, error) {
	var h core.Hash

	if !obj.Type().Valid() {
		return h, core.ErrInvalidType
	}

	fso, ok := obj.(*Object)
	if !ok {
		return h, errors.New("filesystem: object not created by filesystem storage")
	}

	fsop := fso.Path()
	if fsop == "" {
		return h, errors.New("filesystem: object not ready")
	}

	path := o.Path(h)
	if fsop == "" {
		return h, errors.New("filesystem: unable to determine correct file path for object")
	}

	h = obj.Hash()
	if h == core.ZeroHash {
		return h, errors.New("filesystem: object not ready")
	}

	if path == fsop {
		// Already exists in correct location. No-op.
		return h, nil
	}

	if err := o.mkdir(o.Dir(h)); err != nil {
		// FIXME: What do we do with the object if it's a temporary file? Delete it?
		return h, err
	}

	if err := os.Rename(fsop, path); err != nil {
		// FIXME: What do we do with the object if it's a temporary file? Delete it?
		return h, err
	}

	// TODO: Set permissions on path? If it was a temp file it was created with 0600.
	// TODO: Consider making file and directory permissions configurable.

	return h, nil
}

// Get returns the object associated with the given hash. If an object with
// the requested hash is not present ErrNotFound is returned.
func (o *ObjectStorage) Get(h core.Hash) (obj core.Object, err error) {
	if h == core.ZeroHash {
		return nil, ErrZeroHash
	}

	path := o.Path(h)

	f, err := os.Open(path)
	if err != nil {
		// TODO: If err == ErrNotFound, try packfile next
		return nil, err
	}

	r, err := NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	defer checkClose(r, &err)

	return &Object{
		t:    r.Type(),
		h:    h,
		size: r.Size(),
		path: path,
	}, nil
}

// Iter returns a core.ObjectIter for the given core.ObjectType.
func (o *ObjectStorage) Iter(t core.ObjectType) core.ObjectIter {
	// TODO: Filter by type

	series := make([]core.Hash, 0, 32)
	prefix := ""
	filepath.Walk(o.root, func(path string, info os.FileInfo, err error) error {
		name := info.Name()
		if info.IsDir() {
			if name == "pack" || name == "info" {
				return filepath.SkipDir
			}
			if len(name) == 2 {
				if _, err := hex.DecodeString(name); err == nil {
					prefix = name
				}
			}
			return nil
		}
		sha1 := prefix + name
		if len(sha1) == 40 {
			if b, err := hex.DecodeString(sha1); err == nil && len(b) == 20 {
				var h core.Hash
				copy(h[:], b)
				series = append(series, h)
			}
		}
		return nil
	})
	return core.NewObjectLookupIter(o, series)
}

// Init ensures that all of the standard git directories have been created.
func (o *ObjectStorage) Init() {
	// TODO: Write this
	return
}

// Root returns the directory that serves as the root of the repository.
func (o *ObjectStorage) Root() string {
	return o.root
}

// Dir returns the path of the directory for an object with the given hash.
func (o *ObjectStorage) Dir(h core.Hash) string {
	hx := h.String()
	return filepath.Join(o.root, hx[:2])
}

// Path returns the path of the file for an object with the given hash.
func (o *ObjectStorage) Path(h core.Hash) string {
	hx := h.String()
	return filepath.Join(o.root, hx[:2], hx[2:])
}

// mkdir will create the directory specified by path if it doesn't exist
// already.
func (o *ObjectStorage) mkdir(path string) error {
	return os.Mkdir(path, 0660)
}
