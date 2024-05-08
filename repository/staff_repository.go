package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
)

type StaffRepository interface {
	Register(ctx context.Context, pool *pgxpool.Pool, staff *entity.StaffRegisterRequest) (string, error)
	Login(ctx context.Context, pool *pgxpool.Pool, phoneNumber string) (*entity.Staff, error)
}

type staffRepositoryImp struct {
}

func NewStaffRepository() StaffRepository {
	return &staffRepositoryImp{}
}

// Login implements staffRepository.
func (i *staffRepositoryImp) Login(ctx context.Context, pool *pgxpool.Pool, phoneNumber string) (*entity.Staff, error) {
	query := "SELECT id, phone_number, name, password FROM staffs WHERE phone_number = $1"

	rows, err := pool.Query(ctx, query, phoneNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() == false {
		return nil, errors.New("Phone number not found")
	}

	staff := entity.Staff{}

	rows.Scan(&staff.Id, &staff.PhoneNumber, &staff.Name, &staff.Password)

	return &staff, nil
}

//
// Register returning staff id.
func (i *staffRepositoryImp) Register(ctx context.Context, pool *pgxpool.Pool, staff *entity.StaffRegisterRequest) (string, error) {
	var id string
	query := "INSERT INTO staffs (phone_number, name, password) VALUES ($1, $2, $3) RETURNING id"

	row := pool.QueryRow(ctx, query, staff.PhoneNumber, staff.Name, staff.Password)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}
