package entity

type Staff struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

type StaffLoginRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16,startswith=+"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type StaffRegisterRequest struct {
	Name        string `json:"name" validate:"required,min=5,max=50"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16,startswith=+,numeric"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}
