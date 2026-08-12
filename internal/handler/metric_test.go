package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cymiam/metircs-store/internal/handler"
	models "github.com/cymiam/metircs-store/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestData(t *testing.T, server *httptest.Server,
	path string) {

	req := resty.New().R()
	req.Method = "POST"
	req.URL = server.URL + path

	_, err := req.Send()
	assert.NoError(t, err, "error making HTTP request")
}

func TestMetricHandler(t *testing.T) {
	server := httptest.NewServer(handler.MetricRouter())
	defer server.Close()

	createTestData(t, server, "/update/counter/testCounter/1")
	createTestData(t, server, "/update/counter/testCounter2/2")
	createTestData(t, server, "/update/counter/testCounter2/2")
	createTestData(t, server, "/update/counter/testCounter2/2")
	createTestData(t, server, "/update/counter/testCounter2/2")
	createTestData(t, server, "/update/counter/testCounter2/2")
	createTestData(t, server, "/update/gauge/testGauge/3.14")

	type want struct {
		statusCode  int
		contentType string
		body        io.ReadCloser
		value       string
	}

	tests := []struct {
		name   string
		url    string
		body   models.Metrics
		method string
		want   want
	}{
		{
			name:   "GET Gauge",
			url:    "/value/gauge/testGauge",
			method: "GET",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				value:       "3.14",
			},
		},
		{
			name:   "GET Single Counter",
			url:    "/value/counter/testCounter",
			method: "GET",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				value:       "1",
			},
		},
		{
			name:   "GET Multiple Counters",
			url:    "/value/counter/testCounter2",
			method: "GET",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				value:       "10",
			},
		},
		{
			name:   "GET Unknown Counter",
			url:    "/value/counter/unknownCounter",
			method: "GET",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				value:       "",
			},
		},
		{
			name:   "GET Unknow gauge",
			url:    "/value/gauge/unknownGauge",
			method: "GET",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				value:       "",
			},
		},
		{
			name:   "POST Gauge",
			url:    "/update/gauge/test/3.0",
			method: "POST",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "POST Counter",
			url:    "/update/counter/test2/5",
			method: "POST",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "POST Unknown metric type",
			url:    "/update/timeseries/requests/5",
			method: "POST",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				value:       "Неизвестный тип метрики\n",
			},
		},
		{
			name:   "POST No metric name",
			url:    "/update/counter/5",
			method: "POST",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				value:       "404 page not found\n",
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = server.URL + tc.url

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			assert.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))

			if resp.Header().Get("Content-Type") == "text/plain; charset=utf-8" {
				require.Equal(t, tc.want.value, string(resp.Body()))
			}
		})

	}
}
