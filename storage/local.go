package storage

import (
	"bytes"
	"errors"
	writablefs "github.com/thewizardplusplus/go-writable-fs"
	fsutils "github.com/thewizardplusplus/go-writable-fs/fs-utils"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type File struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Mime         string    `json:"mime"`
	Ext          string    `json:"ext"`
	Size         int       `json:"size"`
	IsDir        bool      `json:"isDir"`
	LastModified time.Time `json:"lastModified"`
	Editable     bool      `json:"editable"`
}

type LocalStorage struct {
	fs writablefs.DirFS
}

func NewLocalStorage(root string) *LocalStorage {
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(root, os.ModeDir)
	}

	return &LocalStorage{fs: writablefs.NewDirFS(root)}
}

func (ls LocalStorage) Create(path string, force bool) (*File, error) {
	path = strings.TrimLeft(path, "/")

	if force {
		dir := strings.Trim(strings.ReplaceAll(path, filepath.Base(path), ""), "/")
		_ = fsutils.MkdirAll(ls.fs, dir, fs.ModeDir)
	}

	f, err := ls.fs.Create(path)
	if err != nil {
		return nil, unwrap(err)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, unwrap(err)
	}
	_ = f.Close()

	return file(stat, path), nil
}

func (ls LocalStorage) CreateDir(path string) (*File, error) {
	path = strings.TrimLeft(path, "/")

	err := fsutils.MkdirAll(ls.fs, path, fs.ModeDir)
	if err != nil {
		return nil, unwrap(err)
	}

	fl, err := ls.fs.Open(path)
	if err != nil {
		return nil, unwrap(err)
	}

	stat, err := fl.Stat()
	if err != nil {
		return nil, unwrap(err)
	}

	_ = fl.Close()

	return file(stat, path), nil
}

func (ls LocalStorage) Move(old string, new string) (*File, error) {
	old = strings.TrimLeft(old, "/")
	new = strings.TrimLeft(new, "/")

	err := ls.fs.Rename(old, new)
	if err != nil {
		return nil, unwrap(err)
	}

	fl, err := ls.fs.Open(new)
	if err != nil {
		return nil, unwrap(err)
	}

	stat, err := fl.Stat()
	if err != nil {
		return nil, unwrap(err)
	}

	_ = fl.Close()

	return file(stat, new), nil
}

func (ls LocalStorage) Remove(path string, force bool) error {
	path = strings.TrimLeft(path, "/")

	if force == true {
		err := fsutils.RemoveAll(ls.fs, path)
		if err != nil {
			return unwrap(err)
		}
	} else {
		err := ls.fs.Remove(path)
		if err != nil {
			return unwrap(err)
		}
	}

	return nil
}

func (ls LocalStorage) Update(path string, body []byte) (*File, error) {
	path = strings.TrimLeft(path, "/")

	fl, err := ls.fs.Create(path)
	if err != nil {
		return nil, unwrap(err)
	}

	stat, err := fl.Stat()
	if err != nil {
		return nil, unwrap(err)
	}

	if !isEditableStat(stat) {
		_ = fl.Close()
		return nil, errors.New("file not editable")
	}

	_, err = fl.Write(body)
	if err != nil {
		return nil, unwrap(err)
	}

	_ = fl.Close()

	return file(stat, path), nil
}

func (ls LocalStorage) Cat(path string) (*string, error) {
	path = strings.TrimLeft(path, "/")

	fl, err := ls.fs.Open(path)
	if err != nil {
		return nil, unwrap(err)
	}

	stat, err := fl.Stat()
	if err != nil {
		return nil, unwrap(err)
	}

	if !isEditableStat(stat) {
		_ = fl.Close()
		return nil, errors.New("file not editable")
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(fl)
	if err != nil {
		return nil, unwrap(err)
	}

	body := buf.String()
	return &body, nil
}

func (ls LocalStorage) Tree() ([]*File, error) {
	//path = strings.TrimLeft(path, "/")
	files := make([]*File, 0)

	err := fs.WalkDir(ls.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "." {
			return nil
		}

		stat, _ := fs.Stat(ls.fs, path)
		files = append(files, file(stat, path))

		return nil
	})
	if err != nil {
		return nil, unwrap(err)
	}

	return files, nil
}

func (ls LocalStorage) Upload(path string, body []byte) (*File, error) {
	path = strings.TrimLeft(path, "/")

	fl, err := ls.fs.Create(path)
	if err != nil {
		return nil, unwrap(err)
	}

	stat, err := fl.Stat()
	if err != nil {
		return nil, unwrap(err)
	}

	_, err = fl.Write(body)
	if err != nil {
		return nil, unwrap(err)
	}

	_ = fl.Close()

	return file(stat, path), nil
}

func (ls LocalStorage) Exists(path string) bool {
	path = strings.TrimLeft(path, "/")

	_, err := ls.fs.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func isEditable(mimeType string) bool {
	mimes := []string{"text/css", "text/csv", "text/html", "application/json", "text/javascript", "text/javascript", "application/x-httpd-php", "text/plain", "application/xml"}
	return slices.Contains(mimes, mimeType)
}

func isEditableStat(st fs.FileInfo) bool {
	ext := strings.ToLower(strings.Trim(filepath.Ext(st.Name()), "."))
	mimeType := strings.Split(mime.TypeByExtension("."+ext), ";")[0]
	return !st.IsDir() && isEditable(mimeType)
}

func file(st fs.FileInfo, path string) *File {
	size := int(st.Size())
	ext := strings.ToLower(strings.Trim(filepath.Ext(st.Name()), "."))
	mimeType := strings.Split(mime.TypeByExtension("."+ext), ";")[0]

	return &File{
		Name:         st.Name(),
		Path:         path,
		Mime:         mimeType,
		Ext:          ext,
		Size:         size,
		IsDir:        st.IsDir(),
		LastModified: st.ModTime(),
		Editable:     !st.IsDir() && isEditable(mimeType),
	}
}

func unwrap(err error) error {
	value, ok := err.(*fs.PathError)
	if ok {
		return value.Err
	}
	return err
}
