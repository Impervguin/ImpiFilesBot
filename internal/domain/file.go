package domain

import "io"

type File struct {
	Name string
	Path string
}

type FileData struct {
	File
	Content io.Reader
}

type Directory struct {
	Name  string
	Path  string
	Files []File
	Dirs  []Directory
}

func NewFile(name string, path string) (*File, error) {
	return &File{Name: name, Path: path}, nil
}

func NewDirectory(name string, path string) (*Directory, error) {
	return &Directory{Name: name, Path: path}, nil
}

func (d *Directory) AddFile(file *File) {
	d.Files = append(d.Files, *file)
}

func (d *Directory) AddDir(dir *Directory) {
	d.Dirs = append(d.Dirs, *dir)
}

func NewFileData(file *File, content io.Reader) (*FileData, error) {
	return &FileData{File: *file, Content: content}, nil
}
