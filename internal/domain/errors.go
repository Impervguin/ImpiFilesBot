package domain

import "fmt"

var ErrFileSizeLimit = fmt.Errorf("file size limit(%v) exceeded", FileSizeLimit)
var ErrNoAccess = fmt.Errorf("no access")

var ErrDirNotFound = fmt.Errorf("directory not found")

var ErrNotDirectory = fmt.Errorf("not a directory")
var ErrNotFile = fmt.Errorf("not a file")

var ErrCwdNotFound = fmt.Errorf("cwd not found")

var ErrUserNotFound = fmt.Errorf("user not found")
