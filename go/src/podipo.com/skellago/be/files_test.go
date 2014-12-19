package be

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/chai2010/assert"
)

func TestLocalFileStorage(t *testing.T) {
	_, err := NewLocalFileStorage("/bogus/mcboog")
	AssertNotNil(t, err, "Should return an error if handed a non-existing directory")

	// Directory for the LocalFileStorage
	fsDir, err := ioutil.TempDir(os.TempDir(), "skellago-test-fs")
	AssertNil(t, err, "Could not create fsDir: "+fsDir)
	defer func() {
		err = os.RemoveAll(fsDir)
		AssertNil(t, err, "Could not clean up fsDir: "+fsDir)
	}()

	// Directory for temporary source files
	tempDir, err := ioutil.TempDir(os.TempDir(), "skellago-temp")
	AssertNil(t, err, "Could not create tempDir: "+tempDir)
	defer func() {
		err = os.RemoveAll(tempDir)
		AssertNil(t, err, "Could not clean up tempDir: "+tempDir)
	}()

	testFS, err := NewLocalFileStorage(fsDir)
	AssertNil(t, err)

	exists, err := testFS.Exists("bogus-key")
	AssertNotNil(t, err)
	AssertFalse(t, exists)

	f1, err := tempFile(tempDir, 10)
	AssertNil(t, err)
	key, err := testFS.Put("foo.bin", f1)
	AssertNil(t, err)
	lf1, err := testFS.Get(key)
	AssertNil(t, err)
	name, err := lf1.Name()
	AssertNil(t, err)
	AssertEqual(t, name, "foo.bin")
	exists, err = lf1.Exists()
	AssertNil(t, err)
	AssertTrue(t, exists)
	size, err := lf1.Size()
	AssertNil(t, err)
	AssertEqual(t, size, int64(10*1024))
	f1.Seek(0, 0)
	lf1Reader, err := lf1.Reader()
	AssertNil(t, err)
	AssertTrue(t, compareReaderData(f1, lf1Reader))
	f1.Seek(0, 0)
	lf1Reader, err = lf1.Reader()
	AssertNil(t, err, "Could not read a second time")
	AssertTrue(t, compareReaderData(f1, lf1Reader), "Second comparison failed")

	f2, err := tempFile(tempDir, 5)
	AssertNil(t, err)
	key2, err := testFS.Put("foo.bin", f2)
	AssertNil(t, err)
	f2.Seek(0, 0)
	lf2, err := testFS.Get(key2)
	AssertNil(t, err)
	lf2Reader, err := lf2.Reader()
	AssertNil(t, err)
	AssertTrue(t, compareReaderData(f2, lf2Reader))

	err = testFS.Delete("bogus-key-2")
	AssertNotNil(t, err)
	err = testFS.Delete(lf1.Key())
	AssertNil(t, err)
	err = testFS.Delete(lf2.Key())
	AssertNil(t, err)
	lf2Reader, err = lf2.Reader()
	AssertNotNil(t, err)
	AssertNotNil(t, testFS.Delete(lf1.Key()), "That key should no longer exist")
	AssertNotNil(t, testFS.Delete(lf2.Key()), "That key should no longer exist")
}

func compareReaderData(file1 io.Reader, file2 io.Reader) bool {
	buf1 := make([]byte, 1024)
	n1 := 0
	buf2 := make([]byte, 1024)
	n2 := 0
	for {
		n1, _ = file1.Read(buf1)
		n2, _ = file2.Read(buf2)
		if n1 != n2 {
			logger.Print("Unbalanced read: ", n1, " ", n2)
			return false
		}
		if bytes.Compare(buf1[0:n1], buf2[0:n2]) != 0 {
			logger.Print("Different buffers: ", buf1[0:n1], " ", buf2[0:n2])
			return false
		}
		if n1 == 0 {
			return true
		}
	}
}

func tempFile(dir string, kilobytes int) (*os.File, error) {
	f, err := ioutil.TempFile(dir, "skella-test-file")
	if err != nil {
		return nil, err
	}
	if kilobytes > 0 {
		data := make([]byte, 1024)
		n := 0
		for i := 0; i < kilobytes; i++ {
			n, err = f.Write(data)
			if err != nil || n != len(data) {
				f.Close()
				return nil, err
			}
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}
