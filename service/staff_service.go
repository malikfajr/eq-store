package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/pkg"
	"github.com/malikfajr/eq-store/repository"
)

type StaffService interface {
	Login(ctx context.Context, req *entity.StaffLoginRequest) (*StaffResponse, error)
	Register(ctx context.Context, req *entity.StaffRegisterRequest) (*StaffResponse, error)
}

type staffService struct {
	staffRepository repository.StaffRepository
	pool            *pgxpool.Pool
}

func NewStaffService(staff repository.StaffRepository, pool *pgxpool.Pool) StaffService {
	return &staffService{
		staffRepository: staff,
		pool:            pool,
	}
}

type StaffResponse struct {
	UserId      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

// Login implements iStaffService.
func (i *staffService) Login(ctx context.Context, req *entity.StaffLoginRequest) (*StaffResponse, error) {
	staff, err := i.staffRepository.Login(context.Background(), i.pool, req.PhoneNumber)
	if err != nil {
		return nil, exception.NewNotFound("user is not found")
	}

	if ok := pkg.ValidPassword(staff.Password, req.Password); ok == false {
		return nil, exception.NewBadRequest("password is wrong")
	}

	token := pkg.CreateToken(staff.Id, staff.Name)

	data := &StaffResponse{
		UserId:      staff.Id,
		PhoneNumber: req.PhoneNumber,
		Name:        staff.Name,
		AccessToken: token,
	}

	return data, nil
}

// Register implements iStaffService.
func (s *staffService) Register(ctx context.Context, req *entity.StaffRegisterRequest) (*StaffResponse, error) {
	hashPassword := pkg.HashPassword(req.Password)
	req.Password = hashPassword

	staffId, err := s.staffRepository.Register(ctx, s.pool, req)
	if err != nil {
		return nil, err
	}

	token := pkg.CreateToken(staffId, req.Name)

	data := &StaffResponse{
		UserId:      staffId,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		AccessToken: token,
	}

	return data, nil
}
