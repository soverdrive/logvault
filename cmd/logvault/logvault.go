package main

import (
	"flag"
	"log"
	"net"

	pb "github.com/albert-widi/logvault/pb"
	"github.com/albert-widi/logvault/service"
	"github.com/albert-widi/logvault/storage/logger"
	"google.golang.org/grpc"
)

var (
	fileLogDir = flag.String("file_log", "", "file log directory")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", ":9300")
	if err != nil {
		log.Fatal("Failed to listen to default address ", err.Error())
	}
	l := logger.NewFileLogger(*fileLogDir)
	defer l.Close()

	grpcServer := grpc.NewServer()
	rpcService := &service.RPC{Logger: l}
	pb.RegisterLogvaultServer(grpcServer, rpcService)

	log.Println("Logvault service start...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Cannot start grpc server ", err.Error())
	}
	log.Println("Logvault service exit...")
}
