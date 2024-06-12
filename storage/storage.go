package storage

import pb "github.com/Mubinabd/payment_service/genproto"

type StorageI interface {
	Payment() PaymentI
}

type PaymentI interface {
	CreatePayment(payment *pb.PaymentCreate) (*pb.Payment, error)
	GetPayment(id *pb.ById) (*pb.Payment, error)
	UpdatePayment(payment *pb.PaymentCreate) (*pb.Payment, error) 
	GetPayments(filter *pb.PaymentFilter) (*pb.Payments, error)
}