package middleware

import (
	"context"
	"net/http"

	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/pkg"
)

const AuthStaffID = "auth.staff.id"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authorization := r.Header.Get("Authorization")

		Unauthorize := exception.NewUnauthorized("Invalid token")
		if len(Authorization) < 8 {
			Unauthorize.Send(w)
			return
		}

		if Authorization[:7] != "Bearer " {
			Unauthorize.Send(w)
			return
		}

		token := Authorization[7:]
		jwt, err := pkg.ClaimToken(token)
		if err != nil {
			Unauthorize.Send(w)
			return
		}

		ctx := context.WithValue(r.Context(), AuthStaffID, jwt.StaffId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
