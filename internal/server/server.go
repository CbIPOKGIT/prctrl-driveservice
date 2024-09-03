package server

import (
	"log"
	"net"
	"os"

	"github.com/CbIPOKGIT/prctrl-driveservice/driveservice"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/usecases/middlewares"
	"google.golang.org/grpc"
)

type DriveService struct {
	driveservice.UnimplementedDriveServiceServer
}

func New() *DriveService {
	return &DriveService{}
}

func (ds *DriveService) Start() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middlewares.TokenMiddleware()),
		grpc.MaxRecvMsgSize(1e8),
		grpc.MaxSendMsgSize(1e8),
	)

	driveservice.RegisterDriveServiceServer(server, ds)

	log.Printf("Starting serverer on port %s", os.Getenv("PORT"))
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
