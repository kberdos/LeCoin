package vfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// a dummy interface for an FS with some basic operations

type VFS struct {
	root string
	cwd  string
}

func NewVFS(root string) *VFS {
	return &VFS{
		root: root,
		cwd:  "/",
	}
}

func (vfs *VFS) WriteFile(path string, data []byte) error {
	truepath, _ := vfs.resolvePath(path)
	return os.WriteFile(truepath, data, fs.ModePerm)
}

func (vfs *VFS) ReadFile(path string) ([]byte, error) {
	truepath, _ := vfs.resolvePath(path)
	return os.ReadFile(truepath)
}

func (vfs *VFS) Mkdir(path string) error {
	truepath, _ := vfs.resolvePath(path)
	return os.Mkdir(truepath, fs.ModePerm)
}

func (vfs *VFS) Cd(path string) error {
	_, abspath := vfs.resolvePath(path)
	if !vfs.Exists(abspath) {
		return errors.New("cd: path dne")
	}

	vfs.cwd = abspath

	return nil
}

func (vfs *VFS) GetCwd() string {
	return vfs.cwd
}

func (vfs *VFS) Exists(path string) bool {
	truepath, _ := vfs.resolvePath(path)
	_, err := os.Stat(truepath)
	if err == nil {
		return true
	}

	return false
}

func (vfs *VFS) Listdir(path string) ([]os.DirEntry, error) {
	truepath, _ := vfs.resolvePath(path)

	entries, err := os.ReadDir(truepath)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// convert path to true os path, and resolved rel path
func (vfs *VFS) resolvePath(path string) (truepath string, abspath string) {
	if !strings.HasPrefix(path, "/") {
		path = filepath.Join("/", vfs.cwd, path)
	}

	abspath = filepath.Clean(path)
	if strings.HasPrefix(abspath, "..") {
		abspath = "/"
	}
	truepath = filepath.Join(vfs.root, abspath)

	return truepath, abspath
}
