package main

import (
	"log"
	"net"

	pb "github.com/albert-widi/logvault/pb"
	"github.com/albert-widi/logvault/service"
	"github.com/albert-widi/logvault/storage/logger"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":9300")
	if err != nil {
		log.Fatal("Failed to listen to default address ", err.Error())
	}
	l := logger.NewFileLogger()
	defer l.Close()

	grpcServer := grpc.NewServer()
	rpcService := &service.RPC{Logger: l}
	pb.RegisterLogvaultServer(grpcServer, rpcService)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Cannot start grpc server ", err.Error())
	}
}
