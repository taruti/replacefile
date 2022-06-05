// +build linux
package replacefile

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"syscall"
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
func ReplaceFile(filename string, data []byte, perm os.FileMode) (err error) {
	// Get the real filename behind links
	dest, err := RealPath(filename)
	if err == nil {
		filename = dest
	}

	// Get the temporary filename
	tmp := make([]byte, 4)
	_, err = io.ReadFull(rand.Reader, tmp)
	if err != nil {
		return
	}
	tmpfile := fmt.Sprintf("%s.temp.%x", filename, tmp)

	// Open the temporary file
	f, err := os.OpenFile(tmpfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, perm)
	if err != nil {
		return
	}

	// Fix the permissions, we don't want umask messing things up
	// Stat the original file and copy permissions, may fail if the original does not exist
	fi, err := os.Stat(filename)
	if err == nil {
		st := fi.Sys().(*syscall.Stat_t)
		f.Chmod(os.FileMode(st.Mode & 0777))
		f.Chown(int(st.Uid), int(st.Gid))
	}

	// Write all the data
	n, err := f.Write(data)
	if err != nil {
		return
	}
	if n < len(data) {
		return io.ErrShortWrite
	}

	// Fsync
	err = syscall.Fsync(int(f.Fd()))
	er2:= f.Close()
	if err != nil {
		return err
	}
	if er2 != nil {
		return er2
	}

	// Rename
	err = os.Rename(tmpfile, filename)

	return
}
