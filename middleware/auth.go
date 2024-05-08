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

		if Authorization[:7] != "Bearer " {
			e := exception.NewUnauthorized("Invalid token")
			e.Send(w)
			return
		}

		token := Authorization[7:]
		jwt, err := pkg.ClaimToken(token)
		if err != nil {
			e := exception.NewUnauthorized("Invalid token")
			e.Send(w)
			return
		}

		ctx := context.WithValue(r.Context(), AuthStaffID, jwt.StaffId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
