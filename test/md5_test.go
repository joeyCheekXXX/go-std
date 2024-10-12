package test

import (
	"github.com/joeyCheek888/go-std/encrypt/md5"
	"testing"
)

func TestMd5(t *testing.T) {

	t.Log(md5.MD5("hello"))

}
