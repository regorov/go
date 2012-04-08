package resourcemanager

import (
    "os"
    "strings"
)

type ResourceManager struct {
    rootDir string
}

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

func (rm *ResourceManager) GetFilename(path string) (filename string) {
    return rm.rootDir + path
}

func (rm *ResourceManager) GetFile(path string) (file *os.File, err error) {
    return os.Open(rm.rootDir + path)
}
