package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cymiam/metircs-store/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestMetricHandler_HandleUpdate(t *testing.T) {

	mux := http.NewServeMux()
	metricHandler := handler.NewMetricHandler()
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", metricHandler.HandleUpdate)

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "Test Positive Gauge",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/gauge/test/3.0",
		},
		{
			name: "Test Positive Counter",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/counter/test2/5",
		},
		{
			name: "Test Unknown metric type",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/timeseries/requests/5",
		},
		{
			name: "Test no metric name",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
			request: "/update/counter/5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			fmt.Println(tt.request)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

		})
	}
}
