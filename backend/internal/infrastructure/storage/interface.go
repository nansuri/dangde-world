package storage

type Storage interface {
	Save(path string, content []byte) error
	Read(path string) ([]byte, error)
	Delete(path string) error
}
