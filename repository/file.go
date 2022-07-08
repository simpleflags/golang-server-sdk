package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	path string
	mux  sync.RWMutex
}

func NewFileStorage(path string) (FileStorage, error) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return FileStorage{}, err
	}
	return FileStorage{
		path: path,
	}, nil
}

func (f *FileStorage) Get(key string, output interface{}) error {
	f.mux.RLock()
	defer f.mux.RUnlock()

	filename := key + ".json"
	fullpath := filepath.Join(f.path, filename)
	bytes, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileStorage) Set(key string, value interface{}) error {
	f.mux.RLock()
	f.mux.RUnlock()

	filename := key + ".json"
	fullpath := filepath.Join(f.path, filename)
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullpath, bytes, 0777)
}

func (f *FileStorage) Remove(s string) error {
	return nil
}

func (f *FileStorage) List() []interface{} {
	return nil
}
