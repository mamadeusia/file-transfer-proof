package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	FilePath   string
	OutputFile *os.File
}

func NewFile() *File {
	return &File{}
}

func (f *File) SetFile(index int, collectionHash string, serverPath string) error {
	err := os.MkdirAll(filepath.Join(serverPath, collectionHash), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	f.FilePath = filepath.Join(serverPath, collectionHash, fmt.Sprintf("F_%d", index))
	if _, err := os.Stat(f.FilePath); !os.IsNotExist(err) {
		return errors.New("file already exist")
	}
	file, err := os.Create(f.FilePath)
	if err != nil {
		return err
	}
	f.OutputFile = file
	return nil
}

func (f *File) Write(chunk []byte) error {
	if f.OutputFile == nil {
		return nil
	}
	_, err := f.OutputFile.Write(chunk)
	return err
}

func (f *File) WriteFromReader(reader io.Reader) error {
	if f.OutputFile == nil {
		return nil
	}
	_, err := io.Copy(f.OutputFile, reader)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (f *File) Close() error {
	return f.OutputFile.Close()
}

func (f *File) GetFile(index int, collectionHash string, serverPath string) error {
	f.FilePath = filepath.Join(serverPath, collectionHash, fmt.Sprintf("F_%d", index))
	if _, err := os.Stat(f.FilePath); os.IsNotExist(err) {
		return errors.New("file not exist")
	}
	file, err := os.Open(f.FilePath)
	if err != nil {
		return err
	}
	f.OutputFile = file
	return nil
}

func (f *File) GetBytes() ([]byte, error) {
	return io.ReadAll(f.OutputFile)

}
