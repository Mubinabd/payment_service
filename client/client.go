package client

import (
	"log"

	pb "github.com/Mubinabd/payment_service/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type Clients struct {
	ReservationClient pb.ReservationServiceClient
}

func NewClients() *Clients {
	conn,err := grpc.NewClient("localhost:8088",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err!=nil{
        log.Fatal(err)
    }
	reservationS := pb.NewReservationServiceClient(conn)
	
	return &Clients{
        ReservationClient: reservationS,
    }
	
}