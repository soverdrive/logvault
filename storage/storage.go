package storage

// Storage interface
type Storage interface {
	WriteLog(prefix string, content string) error
	Close() error
}
