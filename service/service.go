package service

import (
	"context"

	pb "github.com/Mubinabd/payment_service/genproto"
	"github.com/Mubinabd/payment_service/storage"
	"github.com/google/uuid"
)

type PaymentService struct {
	stg storage.StorageI
	pb.UnimplementedPaymentServiceServer
}

func NewPaymentService(stg storage.StorageI) *PaymentService {
	return &PaymentService{stg: stg}
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment *pb.PaymentCreate) (*pb.Payment, error) {
	id := uuid.NewString()
	payment.Id = id
	return s.stg.Payment().CreatePayment(payment)
}

func (s *PaymentService) GetPayment(ctx context.Context, id *pb.ById) (*pb.Payment, error) {
    return s.stg.Payment().GetPayment(id)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, payment *pb.PaymentCreate) (*pb.Payment, error) {
    return s.stg.Payment().UpdatePayment(payment)
}

func (s *PaymentService) GetPayments(ctx context.Context, filter *pb.PaymentFilter) (*pb.Payments, error) {
    return s.stg.Payment().GetPayments(filter)
}
