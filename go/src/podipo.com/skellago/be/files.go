package be

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	keySeparator = "___"
)

/*
FileStorage is the interface used by be.API to persist and retrieve files.
LocalFileStorage implements FileStorage and uses the local file system, but it should also be possible to persist files to services like S3.
*/
type FileStorage interface {
	Put(name string, reader io.Reader) (key string, err error)
	Get(key string) (File, error)
	Exists(key string) (bool, error)
	Delete(key string) error
}

/*
File is the type stored by FileStorage.  LocalFile is an example of File which uses the local file system.
*/
type File interface {
	Key() string
	Name() (string, error)

	Exists() (bool, error)
	Size() (int64, error)

	Reader() (io.Reader, error)
}

/*
LocalFileStorage is a FileStorage persisted by the local file system
*/
type LocalFileStorage struct {
	RootDir string
}

/*
NewLocalFileStorage requires the rootDir exist and be a directory
*/
func NewLocalFileStorage(rootDir string) (*LocalFileStorage, error) {
	stat, err := os.Stat(rootDir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, os.ErrNotExist
	}
	fs := LocalFileStorage{
		RootDir: rootDir,
	}
	return &fs, nil
}

/*
Put stores the data from reader and returns its key. It does this by reading data into a temp file and then moving the file into place
*/
func (fs LocalFileStorage) Put(name string, reader io.Reader) (key string, err error) {
	key = fs.generateKey(name)

	// Create the temp file with the data from reader
	tempDir, err := ioutil.TempDir(os.TempDir(), "lfs")
	if err != nil {
		return "", err
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()
	tempFile, err := ioutil.TempFile(tempDir, "lfs")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(tempFile, reader)
	tempFile.Close()
	if err != nil {
		return "", err
	}

	// the temp file contains the data, now move it into place
	err = os.Rename(tempFile.Name(), path.Join(fs.RootDir, key))
	if err != nil {
		return "", err
	}
	return key, nil
}

func (fs LocalFileStorage) Get(key string) (File, error) {
	lf := LocalFile{
		key:     fs.clean(key),
		rootDir: fs.RootDir,
	}
	exists, err := lf.Exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("No such key")
	}
	return lf, nil
}

func (fs LocalFileStorage) Exists(key string) (bool, error) {
	key = fs.clean(key)
	if key == "" {
		return false, errors.New("Empty file key")
	}
	lf := LocalFile{
		key:     key,
		rootDir: fs.RootDir,
	}
	return lf.Exists()
}

/*
Delete returns an error only if the file exists and there is a problem deleting it.
If it doesn't exist or there is no problem deleting the file then Delete returns nil.
*/
func (fs LocalFileStorage) Delete(key string) error {
	lf := LocalFile{
		key:     key,
		rootDir: fs.RootDir,
	}
	exists, err := lf.Exists()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return lf.delete()
}

// clean removes any characters which may cause trouble in FS names
func (fs LocalFileStorage) clean(token string) string {
	token = strings.Replace(token, "/", "_", -1)
	token = strings.Replace(token, "..", "_", -1)
	token = strings.Replace(token, "...", "_", -1)
	token = strings.Replace(token, "@", "_", -1)
	return token
}

// generateKey returns a FS friendly name which is highly likely to be unique
func (fs LocalFileStorage) generateKey(name string) string {
	slashIndex := strings.Index(name, "/")
	if slashIndex != -1 {
		name = strings.Split(name, "/")[0]
	}
	return UUID() + keySeparator + fs.clean(name)
}

/*
LocalFile is a be.File backed by a LocalFileStorage
*/
type LocalFile struct {
	key     string
	rootDir string
}

func (lf LocalFile) Key() string {
	return lf.key
}

func (lf LocalFile) Name() (string, error) {
	return lf.key[strings.Index(lf.key, keySeparator)+len(keySeparator):], nil
}

func (lf LocalFile) Exists() (bool, error) {
	_, err := os.Stat(lf.path())
	return err == nil, err
}

func (lf LocalFile) Size() (int64, error) {
	stat, err := os.Stat(lf.path())
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}

func (lf LocalFile) Reader() (io.Reader, error) {
	return os.OpenFile(lf.path(), os.O_RDONLY, os.ModePerm)
}

func (lf LocalFile) delete() error {
	return os.Remove(lf.path())
}

// path returns the full path using the rootDir
func (lf LocalFile) path() string {
	return path.Join(lf.rootDir, lf.key)
}
