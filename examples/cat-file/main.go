package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v3"
	"gopkg.in/src-d/go-git.v3/core"
	"gopkg.in/src-d/go-git.v3/storage/filesystem"
)

func main() {
	hash := core.NewHash(os.Args[2])

	// TODO: Add helper function for finding the local git repository instead
	//       of forcing the caller to provide it.
	storage, err := filesystem.NewObjectStorage(os.Args[3])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reading %v from %v...\n", hash, storage.Root())

	obj, err := storage.Get(hash)
	if err != nil {
		panic(err)
	}

	fmt.Println(obj.Type().String())

	switch obj.Type() {
	case core.CommitObject:
		c := new(git.Commit)
		err := c.Decode(obj)
		if err != nil {
			panic(err)
		}
	case core.BlobObject:
	}
}
