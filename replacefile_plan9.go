package replacefile

import (
	"io"
	"os"
)

// ReplaceFile replaces file with new data.
// If the file does not exist, ReplaceFile creates it with permissions perm;
// otherwise ReplaceFile creates a new file and renames it over the old one.
//
// This won't work on Windows.
//
//   Here is a summary on the operation:
//   1) readlink to get the real destination
//   2) create a temporary file
//   3) chown+chmod it like the configuration file
//   4) write and fsync
//   5) rename
func ReplaceFile(filename string, data []byte, perm uint32) (err error) {
	// Open the file
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	// Write all the data
	n, err := f.Write(data)
	if err != nil {
		return
	}
	if n < len(data) {
		return io.ErrShortWrite
	}

	f.Close()
	return
}
