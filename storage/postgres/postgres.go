package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Mubinabd/payment_service/config"
	"github.com/Mubinabd/payment_service/storage"
	_ "github.com/lib/pq"
	// "google.golang.org/genproto/googleapis/storage/v1"
)

type Storage struct {
	db       *sql.DB
	PaymentS storage.PaymentI
}

func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	payS := NewPaymentStorage(db)
	return &Storage{
        db:       db,
        PaymentS: payS,
    }, nil
}
func (s *Storage) Payment() storage.PaymentI {
	if s.PaymentS == nil {
		s.PaymentS = NewPaymentStorage(s.db)
	}
	return s.PaymentS
}
