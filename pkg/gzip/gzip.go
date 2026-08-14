package pkg

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipDecompressMidlleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = gz
			defer gz.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}
