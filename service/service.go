package service

import (
	"log"

	pb "github.com/albert-widi/logvault/pb"
	"github.com/albert-widi/logvault/storage"
	"golang.org/x/net/context"
)

// RPC Service
type RPC struct {
	Logger storage.Storage
}

// PushLog for pushing log to logee service
func (rpc *RPC) PushLog(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	log.Printf("Writing log. Prefix: %s. Log: %s.\n", req.Prefix, req.Log)
	resp := &pb.PushResponse{Status: "OK"}
	err := rpc.Logger.WriteLog(req.Prefix, req.Log)
	if err != nil {
		resp.Status = "FAILED"
	}
	return resp, err
}
