package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/Mubinabd/payment_service/genproto"
	"github.com/stretchr/testify/assert"
)


func TestCreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &PaymentStorage{db: db}

	payment := &pb.PaymentCreate{
		Id:            "payment-id",
		ReservationId: "reservation-id",
		Amount:        100,
		PaymentMethod: "credit_card",
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
	rows := sqlmock.NewRows([]string{"id", "reservation_id", "amount", "payment_method", "payment_status"}).
		AddRow("payment-id", "reservation-id", 100, "credit_card", "completed")

	mock.ExpectQuery(query).
		WithArgs("payment-id", "reservation-id", 100, "credit_card", "completed").
		WillReturnRows(rows)

	res, err := storage.CreatePayment(payment)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "payment-id", res.Id)
	assert.Equal(t, "reservation-id", res.ReservationId)
	assert.Equal(t, 100.0, res.Amount)
	assert.Equal(t, "credit_card", res.PaymentMethod)
	assert.Equal(t, "completed", res.PaymentStatus)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePaymentOverPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &PaymentStorage{db: db}

	payment := &pb.PaymentCreate{
		Id:            "payment-id",
		ReservationId: "reservation-id",
		Amount:        150,
		PaymentMethod: "credit_card",
	}

	res, err := storage.CreatePayment(payment)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "your payment is over by 50, you should pay 100", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePaymentUnderPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &PaymentStorage{db: db}

	payment := &pb.PaymentCreate{
		Id:            "payment-id",
		ReservationId: "reservation-id",
		Amount:        50,
		PaymentMethod: "credit_card",
	}

	res, err := storage.CreatePayment(payment)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "your payment is under by 50, you should pay 100", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &PaymentStorage{db: db}

	id := &pb.ById{Id: "payment-id"}

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
	rows := sqlmock.NewRows([]string{"id", "reservation_id", "amount", "payment_method", "payment_status"}).
		AddRow("payment-id", "reservation-id", 100, "credit_card", "completed")

	mock.ExpectQuery(query).
		WithArgs(id.Id).
		WillReturnRows(rows)

	res, err := storage.GetPayment(id)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "payment-id", res.Id)
	assert.Equal(t, "reservation-id", res.ReservationId)
	assert.Equal(t, 100.0, res.Amount)
	assert.Equal(t, "credit_card", res.PaymentMethod)
	assert.Equal(t, "completed", res.PaymentStatus)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &PaymentStorage{db: db}

	payment := &pb.PaymentCreate{
		Id:            "payment-id",
		ReservationId: "reservation-id",
		Amount:        100,
		PaymentMethod: "credit_card",
	}

	query := `
		UPDATE payment SET 
			reservation_id = $1, 
			amount = $2, 
			payment_method = $3, 
			payment_status = $4 
		WHERE id = $5
		RETURNING id, reservation_id, amount, payment_method, payment_status
	`
	rows := sqlmock.NewRows([]string{"id", "reservation_id", "amount", "payment_method", "payment_status"}).
		AddRow("payment-id", "reservation-id", 100, "credit_card", "completed")

	mock.ExpectQuery(query).
		WithArgs("reservation-id", 100, "credit_card", "completed", "payment-id").
		WillReturnRows(rows)

	res, err := storage.UpdatePayment(payment)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "payment-id", res.Id)
	assert.Equal(t, "reservation-id", res.ReservationId)
	assert.Equal(t, 100.0, res.Amount)
	assert.Equal(t, "credit_card", res.PaymentMethod)
	assert.Equal(t, "completed", res.PaymentStatus)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &PaymentStorage{db: db}

	filter := &pb.PaymentFilter{
		PaymentMethod: "credit_card",
		PaymentStatus: "completed",
		AmountFrom:    50,
		AmountTo:      150,
	}

	query := `
		SELECT 
			id,
			reservation_id, 
			amount,
			payment_method,
			payment_status
		FROM payment
		WHERE payment_method = $1 AND payment_status = $2 AND amount >= $3 AND amount <= $4
	`
	rows := sqlmock.NewRows([]string{"id", "reservation_id", "amount", "payment_method", "payment_status"}).
		AddRow("payment-id-1", "reservation-id-1", 100, "credit_card", "completed").
		AddRow("payment-id-2", "reservation-id-2", 150, "credit_card", "completed")

	mock.ExpectQuery(query).
		WithArgs("credit_card", "completed", 50, 150).
		WillReturnRows(rows)

	res, err := storage.GetPayments(filter)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res.Payments))
	assert.Equal(t, "payment-id-1", res.Payments[0].Id)
	assert.Equal(t, "payment-id-2", res.Payments[1].Id)

	assert.NoError(t, mock.ExpectationsWereMet())
}
