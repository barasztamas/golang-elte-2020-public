// Binary cksum calculates checksums for files.
package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// TODO: parallelize the checksum calculation
	for _, path := range Files() {
		hash, err := Hash(path)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		fmt.Printf("%x\t%s\n", hash, path)
	}
	// END OMIT
}

// Hash calculates a checksum of a file.
// It returns an error, if the file was not readable.
func Hash(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Files returns the list of file paths that are expanded from walking the tree
// of every command line arguments.
func Files() []string {
	var files []string
	flag.Parse()
	for _, path := range flag.Args() {
		// Walk will return no error, because all WalkFunc always returns nil.
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("ERROR: unable to access %q\n", path)
				return nil
			}
			if info.Mode()&os.ModeType != 0 {
				return nil // Not a regular file.
			}
			files = append(files, path)
			return nil
		})
	}
	return files
}
