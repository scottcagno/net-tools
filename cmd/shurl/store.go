package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Store interface {
	Add(string, []string)
	Set(string, []string)
	Get(string) []string
	Del(string)
}

type SimpleStore struct {
	mu   sync.Mutex
	path string
	data map[string][]string
	fd   *os.File
}

func NewSimpleStore(path string) *SimpleStore {
	err := CreateDirIfNotExist(path)
	if err != nil {
		log.Fatal(err)
	}
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	ss := &SimpleStore{
		path: path,
		data: make(map[string][]string),
	}
	ss.readFileIntoMap()
	return ss
}

func (s *SimpleStore) readFileIntoMap() {
	files, err := os.ReadDir(s.path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b, err := os.ReadFile(s.path + file.Name())
		if err != nil {
			log.Fatalf("error reading file: %v\n", err)
		}
		var ss []string
		err = json.Unmarshal(b, &ss)
		if err != nil {
			log.Fatalf("error unmarshaling: %v\n", err)
		}
		s.data[file.Name()] = ss
	}
}

func (s *SimpleStore) deleteEntry(k string) error {
	// look for contents on disk
	fi, err := os.Stat(s.path + k + ".json")
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("error, can't locate file (%s): %v\n",
			s.path+k+".json", err)
	}
	// contents are there, so lets remove them
	err = os.Remove(fi.Name())
	if err != nil {
		return fmt.Errorf("error removing file: %v\n", err)
	}
	// check for entry in map and remove
	if _, ok := s.data[k]; ok {
		delete(s.data, k)
	}
	return nil
}

func (s *SimpleStore) writeEntry(overwrite bool, k string, v []string) error {
	// check if key exists
	if _, ok := s.data[k]; ok {
		if overwrite {
			err := s.deleteEntry(k)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// if no key found, then marshal data...
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshaling: %v\n", err)
	}
	// ...and write contents to file
	err = os.WriteFile(s.path+k+".json", b, 0666)
	if err != nil {
		return fmt.Errorf("error writing file: %v\n", err)
	}
	// update map if overwrite is true
	if overwrite {
		s.data[k] = v
	}
	return nil
}

func (s *SimpleStore) readEntryFile(k string) ([]string, error) {
	// look for contents in map
	if v, ok := s.data[k]; ok {
		return v, nil
	}
	// look for contents on disk
	fi, err := os.Stat(s.path + k + ".json")
	if err != nil && os.IsNotExist(err) {
		return nil,
			fmt.Errorf("error, can't locate file (%s): %v\n",
				s.path+k+".json", err)
	}
	// read file contents from disk (i really dont think this should happen)
	b, err := os.ReadFile(fi.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v\n", err)
	}
	// unmarshal contents and add back into map (this also really should happen)
	var ss []string
	err = json.Unmarshal(b, &ss)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling: %v\n", err)
	}
	s.data[k] = ss
	return s.data[k], nil
}

func (s *SimpleStore) Add(k string, v []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.writeEntry(false, k, v)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SimpleStore) Set(k string, v []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.writeEntry(true, k, v)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SimpleStore) Get(k string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.readEntryFile(k)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func (s *SimpleStore) Del(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.deleteEntry(k)
	if err != nil {
		log.Fatal(err)
	}
}

func StoreTest() {

	fmt.Println("Creating new simple store")
	s := NewSimpleStore(".")
	time.Sleep(time.Second * 3)

	fmt.Printf("adding %q\n", "foo")
	s.Add("foo", []string{"foo1", "foo2", "foo3"})
	time.Sleep(time.Second * 3)

	fmt.Printf("adding %q\n", "bar")
	s.Add("bar", []string{"baz"})
	time.Sleep(time.Second * 3)

	fmt.Printf("getting %q->%v\n", "foo", s.Get("foo"))
	time.Sleep(time.Second * 3)

	fmt.Printf("setting %q\n", "foo")
	s.Set("foo", []string{"bar"})
	time.Sleep(time.Second * 3)

	fmt.Printf("getting %q->%v\n", "foo", s.Get("foo"))
	time.Sleep(time.Second * 3)

	fmt.Printf("setting %q\n", "bar")
	s.Set("bar", []string{"dang"})
	time.Sleep(time.Second * 3)

	fmt.Printf("getting %q->%v\n", "bar", s.Get("bar"))
	time.Sleep(time.Second * 3)

	fmt.Printf("removing %q\n", "bar")
	s.Del("bar")
	time.Sleep(time.Second * 3)
}

func listDir(path string) []os.FileInfo {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}
