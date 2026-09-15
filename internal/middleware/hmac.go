package middleware

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"io"
	"net/http"

	h "github.com/cymiam/metrics-store/pkg/hmac"
)

func HMACMiddleware(key string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// проверяем, что клиент отправил серверу хеш сумму
			hashSum := r.Header.Get("HashSHA256")
			if hashSum != "" {
				body, err := io.ReadAll(r.Body)
				if err != nil {

					http.Error(w, "Bad Request", http.StatusBadRequest)
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))

				got, err := hex.DecodeString(hashSum)
				if err != nil {
					http.Error(w, "Bad Request", http.StatusBadRequest)
				}
				hash, err := hex.DecodeString(h.CalculateSha256Sum(body, key))
				if err != nil {

					http.Error(w, "Bad Request", http.StatusBadRequest)
				}
				if !hmac.Equal(got, hash) {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
			}

			// передаём управление хендлеру
			next.ServeHTTP(w, r)
		})
	}
}
