package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/TechCatsLab/github-detector/pkg/filetool"
)

// Map -
type Map struct {
	objs *sync.Map
}

// Open -
func Open(path string) (*Map, error) {
	if path == "" {
		return &Map{
			objs: &sync.Map{},
		}, nil
	}

	abspath, err := filetool.Abs(path)
	if err != nil {
		return nil, err
	}
	ext := filetool.Ext(abspath)

	f, err := filetool.Open(abspath, filetool.RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[interface{}]interface{})

	switch ext {
	case ".json":
		err = json.NewDecoder(f).Decode(&m)
		if err != nil && err != io.EOF {
			return nil, err
		}
	case ".yaml", ".yml":
		err = yaml.NewDecoder(f).Decode(&m)
		if err != nil && err != io.EOF {
			return nil, err
		}
	default:
		return nil, errors.New("not support")
	}

	objs := &sync.Map{}
	for k, v := range m {
		objs.Store(k, v)
	}

	return &Map{
		objs: objs,
	}, nil
}

// Save -
func (m *Map) Save(path string) error {
	if m.objs == nil {
		return errors.New("map is closed")
	}

	abspath, err := filetool.Abs(path)
	if err != nil {
		return err
	}
	ext := filetool.Ext(abspath)

	f, err := filetool.Open(abspath, filetool.TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	kvs := make(map[string]interface{})
	m.objs.Range(func(key interface{}, value interface{}) bool {
		kvs[fmt.Sprintf("%v", key)] = value
		return true
	})

	switch ext {
	case ".json":
		return filetool.NewEncoder(f).Encode(&kvs)
	case ".yaml", ".yml":
		encoder := yaml.NewEncoder(f)
		defer encoder.Close()
		return encoder.Encode(&kvs)

	default:
		return errors.New("not support")
	}
}

// Close -
func (m *Map) Close() {
	m.objs = nil
}

// Upsert -
func (m *Map) Upsert(key, value interface{}) {
	if m.objs != nil {
		m.objs.Store(key, value)
	}
}

// Delete -
func (m *Map) Delete(key interface{}) {
	if m.objs != nil {
		m.objs.Delete(key)
	}
}

// Find -
func (m *Map) Find(key interface{}, value interface{}) (err error) {
	if m.objs == nil {
		return errors.New("map is closed")
	}

	v, ok := m.objs.Load(key)
	if !ok {
		return errors.New("not found")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, value)
}
