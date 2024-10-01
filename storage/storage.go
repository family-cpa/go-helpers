package storage

type Storage interface {
	UploadFile(key string, upload Upload) error
	URL(path string) string
}
