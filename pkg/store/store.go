package store

type Store interface {
	Save(path string, v interface{}) error
	SaveAll(path string, vv []interface{}) error
	Load(path string, v interface{}) error
	LoadAll(path string, vv []interface{}) error
	Delete(path string) error
	DeleteAll(path string) error
}
