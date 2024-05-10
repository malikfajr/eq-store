package entity

type Customer struct {
	UserId      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

type CustomerQueryParams struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

type CustomerInsertUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=5,max=50"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16,startswith=+,valid_phone"`
}
