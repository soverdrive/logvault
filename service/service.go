package service

import (
	pb "github.com/albert-widi/logvault/pb"
	"github.com/albert-widi/logvault/storage"
	"golang.org/x/net/context"
)

// RPC Service
type RPC struct {
	Logger storage.Storage
}

// IngestLog for ingest log to logvault service
func (rpc *RPC) IngestLog(ctx context.Context, req *pb.IngestRequest) (*pb.IngestResponse, error) {
	resp := &pb.IngestResponse{Status: "OK"}
	err := rpc.Logger.WriteLog(req)
	if err != nil {
		resp.Status = "FAILED"
	}
	return resp, err
}
