package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Mubinabd/payment_service/client"
	pb "github.com/Mubinabd/payment_service/genproto"
)

type PaymentStorage struct {
	db           *sql.DB
	ReservationC pb.ReservationServiceClient
}

func NewPaymentStorage(db *sql.DB) *PaymentStorage {
	ResClient := client.NewClients()
	return &PaymentStorage{
		db:           db,
		ReservationC: ResClient.ReservationClient,
	}
}

func (s *PaymentStorage) CreatePayment(payment *pb.PaymentCreate) (*pb.Payment, error) {
	res_id := pb.ById{Id: payment.ReservationId}
	total_sum, err := s.ReservationC.GetTotalSum(context.Background(), &res_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get total sum from reservation service: %w", err)
	}

	if payment.Amount > total_sum.Total {
		error_message := fmt.Sprintf("your payment is over by %v, you should pay %v", payment.Amount-total_sum.Total, total_sum.Total)
		return nil, errors.New(error_message)
	} else if payment.Amount < total_sum.Total {
		error_message := fmt.Sprintf("your payment is under by %v, you should pay %v", total_sum.Total-payment.Amount, total_sum.Total)
		return nil, errors.New(error_message)
	}

	query := `
		INSERT INTO payment (
			id, 
			reservation_id, 
			amount,
			payment_method,
			payment_status
		) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id,
			reservation_id, 
			amount,
			payment_method,
			payment_status
	`
	status := "completed"
	row := s.db.QueryRow(query,
		payment.Id,
		payment.ReservationId,
		payment.Amount,
		payment.PaymentMethod,
		status,
	)

	var res pb.Payment
	err = row.Scan(&res.Id,
		&res.ReservationId,
		&res.Amount,
		&res.PaymentMethod,
		&res.PaymentStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan inserted payment: %w", err)
	}

	return &res, nil
}
func (s *PaymentStorage) GetPayment(id *pb.ById) (*pb.Payment, error) {
	query := `
        SELECT 
            id,
            reservation_id, 
            amount,
            payment_method,
            payment_status
        FROM payment
        WHERE id = $1
    `
	row := s.db.QueryRow(query, id.Id)

	var res pb.Payment
	err := row.Scan(&res.Id,
		&res.ReservationId,
		&res.Amount,
		&res.PaymentMethod,
		&res.PaymentStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan payment: %w", err)
	}

	return &res, nil
}
func (s *PaymentStorage) UpdatePayment(payment *pb.PaymentCreate) (*pb.Payment, error) {
	res_id := pb.ById{Id: payment.ReservationId}
	total_sum, err := s.ReservationC.GetTotalSum(context.Background(), &res_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get total sum from reservation service: %w", err)
	}

	if payment.Amount > total_sum.Total {
		error_message := fmt.Sprintf("your payment is over by %v, you should pay %v", payment.Amount-total_sum.Total, total_sum.Total)
		return nil, errors.New(error_message)
	} else if payment.Amount < total_sum.Total {
		error_message := fmt.Sprintf("your payment is under by %v, you should pay %v", total_sum.Total-payment.Amount, total_sum.Total)
		return nil, errors.New(error_message)
	}
	query := `
		update payment set 
				reservation_id = $1, 
				amount = $2, 
				payment_method = $3, 
				payment_status = $4, 
			where id = $5
			returning id,reservation_id,amount,payment_method,payment_status 	
	`
	status := "completed"
	row := s.db.QueryRow(query,
		payment.ReservationId,
		payment.Amount,
		payment.PaymentMethod,
		status,
		payment.Id,
	)
	var res pb.Payment
	err = row.Scan(&res.Id,
		&res.ReservationId,
		&res.Amount,
		&res.PaymentMethod,
		&res.PaymentStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan updated payment: %w", err)
	}
	return &res, nil
}
func (s *PaymentStorage) GetPayments(filter *pb.PaymentFilter) (*pb.Payments, error) {
	query := `
        SELECT 
            id,
            reservation_id, 
            amount,
            payment_method,
            payment_status
        FROM payment
    `
	var conditons []string
	var args []interface{}
	if filter.PaymentMethod != "" {
		conditons = append(conditons, fmt.Sprintf("payment_method = $%d", len(args)+1))
		args = append(args, filter.PaymentMethod)
	}
	if filter.PaymentStatus != "" {
		conditons = append(conditons, fmt.Sprintf("payment_status = $%d", len(args)+1))
		args = append(args, filter.PaymentStatus)
	}
	if filter.AmountFrom != 0 {
		conditons = append(conditons, fmt.Sprintf("amount >= $%d", len(args)+1))
		args = append(args, filter.AmountFrom)
	}
	if filter.AmountTo != 0 {
		conditons = append(conditons, fmt.Sprintf("amount <= $%d", len(args)+1))
		args = append(args, filter.AmountTo)
	}
	if len(conditons) > 0 {
		query += " WHERE " + strings.Join(conditons, " AND ")
	}
	rows, err := s.db.Query(query, args...)
	if err!= nil {
        return nil, fmt.Errorf("failed to query payments: %w", err)
    }
	defer rows.Close()
	var payments pb.Payments
	for rows.Next() {
		var payment pb.Payment
        err := rows.Scan(
            &payment.Id,
            &payment.ReservationId,
            &payment.Amount,
            &payment.PaymentMethod,
            &payment.PaymentStatus,
        )
        if err!= nil {
            return nil, fmt.Errorf("failed to scan payment: %w", err)
        }
        payments.Payments = append(payments.Payments, &payment)
	}
	return &payments, nil
}
