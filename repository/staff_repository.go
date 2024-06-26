package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
)

type StaffRepository interface {
	Register(ctx context.Context, pool *pgxpool.Pool, staff *entity.StaffRegisterRequest) (string, error)
	Login(ctx context.Context, pool *pgxpool.Pool, phoneNumber string) (*entity.Staff, error)
	PhoneIsExist(ctx context.Context, pool *pgxpool.Pool, phoneNumber string) bool
}

type staffRepositoryImp struct {
}

func NewStaffRepository() StaffRepository {
	return &staffRepositoryImp{}
}

// Login implements staffRepository.
func (i *staffRepositoryImp) Login(ctx context.Context, pool *pgxpool.Pool, phoneNumber string) (*entity.Staff, error) {
	query := "SELECT id, phone_number, name, password FROM staffs WHERE phone_number = $1 LIMIT 1"
	staff := &entity.Staff{}

	err := pool.QueryRow(ctx, query, phoneNumber).Scan(&staff.Id, &staff.PhoneNumber, &staff.Name, &staff.Password)
	if err != nil {
		return nil, errors.New("Phone number not found")
	}

	return staff, nil
}

//
// Register returning staff id.
func (i *staffRepositoryImp) Register(ctx context.Context, pool *pgxpool.Pool, staff *entity.StaffRegisterRequest) (string, error) {
	var id string
	query := "INSERT INTO staffs (phone_number, name, password) VALUES ($1, $2, $3) RETURNING id"

	row := pool.QueryRow(ctx, query, staff.PhoneNumber, staff.Name, staff.Password)
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}

	return id, nil
}

func (i *staffRepositoryImp) PhoneIsExist(ctx context.Context, pool *pgxpool.Pool, phoneNumber string) bool {
	var n int
	query := "SELECT 1 FROM staffs WHERE phone_number = $1 LIMIT 1"

	row := pool.QueryRow(ctx, query, phoneNumber)
	err := row.Scan(&n)
	if err != nil {
		log.Print(err)
		return false
	}

	log.Println(phoneNumber, "exists")

	return true
}
