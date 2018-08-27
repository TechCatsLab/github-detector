package yaml

import (
	"errors"
	"io"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/fengyfei/github-detector/pkg/filetool"
)

// File -
type File struct {
	path string
	objs *sync.Map
}

// Open -
func Open(path string) (*File, error) {
	abspath, err := filetool.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := filetool.Open(abspath, filetool.RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[interface{}]interface{})
	err = yaml.NewDecoder(f).Decode(&m)
	if err != nil && err != io.EOF {
		return nil, err
	}

	objs := &sync.Map{}
	for k, v := range m {
		objs.Store(k, v)
	}

	return &File{
		path: abspath,
		objs: objs,
	}, nil
}

// SaveAndClose -
func (f *File) SaveAndClose() (err error) {
	file, err := filetool.Open(f.path, filetool.TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		f.objs = nil
	}()
	defer file.Close()

	m := make(map[interface{}]interface{})
	f.objs.Range(func(key interface{}, value interface{}) bool {
		m[key] = value
		return true
	})

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(&m)
}

// Delete -
func (f *File) Delete(key interface{}) {
	f.objs.Delete(key)
}

// Upsert -
func (f *File) Upsert(key, value interface{}) {
	f.objs.Store(key, value)
}

// Find -
func (f *File) Find(key interface{}, value interface{}) (err error) {
	v, ok := f.objs.Load(key)
	if !ok {
		return errors.New("not found")
	}

	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, value)
}
