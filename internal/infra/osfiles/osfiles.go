package osfiles

import (
	domain "ImpiFilesBot/internal/domain"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// interface implementation
var _ domain.FileRepository = &OsFileSystem{}

type OsFileSystem struct {
	root string
}

func NewOsFileSystem(root string) *OsFileSystem {
	return &OsFileSystem{root: root}
}

func (fs *OsFileSystem) GetByPath(path string) (*domain.File, error) {
	fullPath := filepath.Join(fs.root, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return domain.NewFile(info.Name(), path)
	}

	return nil, fmt.Errorf("directory %s is not file", fullPath)
}

func (fs *OsFileSystem) Save(file *domain.FileData) error {
	fullPath := filepath.Join(fs.root, file.Path)
	finfo, err := os.Stat(fullPath)
	if err == nil && !finfo.IsDir() {
		return fmt.Errorf("file %s already exists", fullPath)
	}
	osFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer osFile.Close()
	_, err = io.Copy(osFile, file.Content)

	return err
}

func (fs *OsFileSystem) Download(path string) (*domain.FileData, error) {
	fullPath := filepath.Join(fs.root, path)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("directory %s is not file", fullPath)
	}
	if fileInfo.Size() > domain.FileSizeLimit {
		return nil, domain.ErrFileSizeLimit
	}
	f, err := domain.NewFile(fileInfo.Name(), path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return domain.NewFileData(f, bytes.NewReader(content))
}

func (fs *OsFileSystem) Delete(path string) error {
	fullPath := filepath.Join(fs.root, path)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (fs *OsFileSystem) List(path string) (*domain.Directory, error) {
	fullPath := filepath.Join(fs.root, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, domain.ErrDirNotFound
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("file %s is not directory", fullPath)
	}
	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	dir, err := domain.NewDirectory(info.Name(), fullPath)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			d, err := domain.NewDirectory(fileInfo.Name(), filepath.Join(fullPath, fileInfo.Name()))
			if err != nil {
				return nil, err
			}
			dir.AddDir(d)
		} else {
			f, err := domain.NewFile(fileInfo.Name(), filepath.Join(fullPath, fileInfo.Name()))
			if err != nil {
				return nil, err
			}
			dir.AddFile(f)
		}
	}
	return dir, nil
}
