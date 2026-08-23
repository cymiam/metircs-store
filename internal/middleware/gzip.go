package middleware

import (
	"net/http"
	"strings"

	"github.com/cymiam/metrics-store/pkg/compress"
)

func GzipDecompressMidlleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {

			gz, err := compress.GzipDecompress(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}
