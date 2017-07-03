package storage

// Storage interface
type Storage interface {
	WriteLog(prefix string, hostname, content string) error
	Close() error
}
