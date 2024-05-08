package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

type JWTClaim struct {
	StaffId string           `json:"staffId"`
	Name    string           `json:"name"`
	Exp     *jwt.NumericDate `json:"exp"`
	jwt.RegisteredClaims
}

func CreateToken(staffId string, name string) string {
	claim := &JWTClaim{
		Name:    name,
		StaffId: staffId,
		Exp:     jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)

	ss, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		panic(err)
	}

	return string(ss)
}

func ClaimToken(token string) (*JWTClaim, error) {
	parsed, err := jwt.ParseWithClaims(token, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsed.Claims.(*JWTClaim); ok {
		return claims, nil
	} else {
		return nil, errors.New("Invalid token")
	}

}
