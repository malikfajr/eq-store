package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)

				w.WriteHeader(500)
				w.Write([]byte("internal server error"))
				return
			}
		}()

		next.ServeHTTP(w, r)

		log.Println(r.Method, r.URL, time.Since(start))
	})
}
