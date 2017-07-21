package storage

import (
	pb "github.com/albert-widi/logvault/pb"
)

// Storage interface
type Storage interface {
	WriteLog(*pb.IngestRequest) error
	Close() error
}
