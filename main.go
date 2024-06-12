package main

import (
	"log"
	"net"

	pb "github.com/Mubinabd/payment_service/genproto"
	"github.com/Mubinabd/payment_service/service"
	"github.com/Mubinabd/payment_service/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	db,err := postgres.ConnectDB()
	if err!= nil {
        log.Fatal(err)
    }
	liss, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s,service.NewPaymentService(db))
	
	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
