package replacefile

import (
	"testing"
)

func TestReplaceFile(t *testing.T) {
	ReplaceFile("/tmp/go-TestReplaceFile", []byte("foobar\n"), 0600)
	RealPath("/home/.///../home//.///././taruti/.emacs")
	RealPath("/../../../../..")
	RealPath("/tmp/..")
	RealPath("/////")
	RealPath("")
	RealPath(".")
}
