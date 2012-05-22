// Package resourcemanager provides a simple way to load files relative to the installed package
// directory. It computes the root directory of the package from the package name and the GOPATH
// environment variable.
package resourcemanager

import (
    "os"
    "strings"
)

// Type ResourceManager represents a resource manager.
type ResourceManager struct {
    rootDir string
}

// Function NewResourceManager creates and returns a new ResourceManager. `packageName` is a string
// containing the package's path (e.g. "github.com/kierdavis/go/resourcemanager").
func NewResourceManager(packageName string) (rm *ResourceManager) {
    gopaths := os.Getenv("GOPATH")
    if gopaths == "" {
        panic("$GOPATH environment variable not found")
    }

    gopath := strings.TrimRight(strings.Split(gopaths, ":")[0], "/")
    rootDir := gopath + "/src/" + packageName + "/"

    rm = new(ResourceManager)
    rm.rootDir = rootDir
    return rm
}

// Function ResourceManager.GetFilename computes the absolute filename of the path `path`, which is
// relative to the root package directory. For example, GetFilename("res/image.png") in package
// "github.com/me/test", where the first element of GOPATH is "/home/me/gocode", would return
// "/home/kier/gocode/src/github/me/test/res/image.png".
func (rm *ResourceManager) GetFilename(path string) (filename string) {
    return rm.rootDir + path
}

// Function ResourceManager.GetFile opens and returns the file specified by the absolute filename of
// the path `path`, which is relative to the root package directory.
func (rm *ResourceManager) GetFile(path string) (file *os.File, err error) {
    return os.Open(rm.rootDir + path)
}
