package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func HMACMiddleware(key string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// проверяем, что клиент отправил серверу хеш сумму
			hashSum := r.Header.Get("HashSHA256")
			if hashSum != "" {
				h := hmac.New(sha256.New, []byte(key))
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Bad Request", http.StatusBadRequest)
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				h.Write(body)
				hash := h.Sum(nil)

				got, err := hex.DecodeString(hashSum)
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
