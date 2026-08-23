package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cymiam/metrics-store/internal/handler"
	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/cymiam/metrics-store/internal/service"
	"github.com/cymiam/metrics-store/pkg/compress"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func createTestData(t *testing.T, server *httptest.Server,
	path string) {

	req := resty.New().R()
	req.Method = "POST"
	req.URL = server.URL + path

	_, err := req.Send()
	assert.NoError(t, err, "error making HTTP request")
}

func createTestDataJson(t *testing.T, server *httptest.Server, path string,
	data models.Metrics) {

	req := resty.New().R()
	req.Method = "POST"
	req.URL = server.URL + path

	body, err := easyjson.Marshal(data)
	require.NoError(t, err, "error Marshall request body")
	req.Body = body
	req.Header.Set("Content-Type", "application/json")

	_, err = req.Send()
	assert.NoError(t, err, "error making HTTP request")
}

func createTestServer() chi.Router {
	repository := repository.NewStore()
	service := service.NewMetricService(service.MetricServiceParams{Store: repository, Logger: zap.NewNop()})
	metricHandler := handler.NewMetricHandler(service, zap.NewNop())

	r := handler.NewMetricRouter(metricHandler, zap.NewNop())
	return r
}

func TestMetricHandler_Update(t *testing.T) {
	server := httptest.NewServer(createTestServer())
	defer server.Close()

	type want struct {
		statusCode  int
		contentType string
		value       string
	}

	tests := []struct {
		name        string
		url         string
		contentType string
		method      string
		want        want
	}{
		{
			name:        "POST Gauge",
			url:         "/update/gauge/test/3.0",
			contentType: "text/plain; charset=utf-8",
			method:      "POST",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "POST Counter",
			url:         "/update/counter/test2/5",
			contentType: "text/plain; charset=utf-8",
			method:      "POST",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "POST Unknown metric type",
			url:         "/update/timeseries/requests/5",
			contentType: "text/plain; charset=utf-8",
			method:      "POST",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				value:       "Неизвестный тип метрики: timeseries\n",
			},
		},
		{
			name:        "POST No metric name",
			url:         "/update/counter/5",
			contentType: "text/plain; charset=utf-8",
			method:      "POST",
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
			req.Header.Set("Content-Type", tc.contentType)

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			require.Equal(t, tc.want.statusCode, resp.StatusCode())
			require.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))

			if tc.want.value != "" {
				assert.Equal(t, tc.want.value, string(resp.Body()))
			}
		})

	}
}

func TestMetricHandler_UpdateJson(t *testing.T) {
	server := httptest.NewServer(createTestServer())
	defer server.Close()

	jsonGaugeValue := 3.1415
	jsonCounterValue := int64(5)

	type want struct {
		statusCode  int
		contentType string
	}

	tests := []struct {
		name        string
		url         string
		body        models.Metrics
		contentType string
		method      string
		want        want
	}{
		{
			name:        "POST Counter JSON",
			url:         "/update",
			contentType: "application/json",
			body: models.Metrics{
				ID:    "CounterJson",
				MType: "counter",
				Delta: &jsonCounterValue,
			},
			method: "POST",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "POST Gauge JSON",
			url:         "/update",
			contentType: "application/json",
			body: models.Metrics{
				ID:    "GaugeJson",
				MType: "gauge",
				Value: &jsonGaugeValue,
			},
			method: "POST",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "POST Unknown metric type JSON",
			url:         "/update",
			contentType: "application/json",
			body: models.Metrics{
				ID:    "CounterJson",
				MType: "timeseries",
			},
			method: "POST",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "POST No metric name JSON",
			url:         "/update",
			contentType: "application/json",
			body: models.Metrics{
				MType: "counter",
				Delta: &jsonCounterValue,
			},
			method: "POST",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = server.URL + tc.url
			req.Header.Set("Content-Type", tc.contentType)

			if tc.contentType == "application/json" {
				body, err := easyjson.Marshal(tc.body)
				require.NoError(t, err, "error Marshall request body")
				req.Body = body
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			require.Equal(t, tc.want.statusCode, resp.StatusCode())
			require.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))
		})

	}
}

func TestMetricHandler_Value(t *testing.T) {
	server := httptest.NewServer(createTestServer())
	defer server.Close()

	createTestData(t, server, "/update/counter/testCounter/2")
	createTestData(t, server, "/update/counter/testCounter/2")
	createTestData(t, server, "/update/counter/testCounter/2")
	createTestData(t, server, "/update/counter/testCounter/2")
	createTestData(t, server, "/update/counter/testCounter/2")
	createTestData(t, server, "/update/gauge/testGauge/3.14")

	type want struct {
		statusCode  int
		contentType string
		value       string
	}

	tests := []struct {
		name        string
		url         string
		contentType string
		method      string
		want        want
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
			name:   "GET Counter",
			url:    "/value/counter/testCounter",
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
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = server.URL + tc.url
			req.Header.Set("Content-Type", tc.contentType)

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			require.Equal(t, tc.want.statusCode, resp.StatusCode())
			require.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))

			if tc.want.value != "" {

				assert.Equal(t, tc.want.value, string(resp.Body()))
			}
		})

	}
}

func TestMetricHandler_ValueJson(t *testing.T) {
	server := httptest.NewServer(createTestServer())
	defer server.Close()

	jsonGaugeValue := 3.1415
	jsonCounterValue := int64(5)

	createTestDataJson(t, server, "/update", models.Metrics{
		MType: "counter",
		ID:    "testCounterJson",
		Delta: &jsonCounterValue,
	})

	createTestDataJson(t, server, "/update", models.Metrics{
		MType: "gauge",
		ID:    "testGaugeJson",
		Value: &jsonGaugeValue,
	})

	type want struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name        string
		url         string
		body        models.Metrics
		contentType string
		method      string
		want        want
	}{{
		name:        "GET Counter JSON",
		url:         "/value",
		method:      "POST",
		contentType: "application/json",
		body: models.Metrics{
			ID:    "testCounterJson",
			MType: "counter",
		},
		want: want{
			statusCode:  http.StatusOK,
			contentType: "application/json",
			body: `{
				"id":    "testCounterJson",
				"type": "counter",
				"delta": 5}`,
		},
	},
		{
			name:        "GET Gauge JSON",
			url:         "/value",
			method:      "POST",
			contentType: "application/json",
			body: models.Metrics{
				ID:    "testGaugeJson",
				MType: "gauge",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				body: `{
				"id":    "testGaugeJson",
				"type": "gauge",
				"value": 3.1415}`,
			},
		},
		{
			name:        "GET Unknown Counter JSON",
			url:         "/value",
			method:      "POST",
			contentType: "application/json",
			body: models.Metrics{
				ID:    "unknownCounter",
				MType: "counter",
			},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "GET Unknown Gauge JSON",
			url:         "/value",
			method:      "POST",
			contentType: "application/json",
			body: models.Metrics{
				ID:    "unknownGauge",
				MType: "counter",
			},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = server.URL + tc.url
			req.Header.Set("Content-Type", tc.contentType)

			if tc.contentType == "application/json" {
				body, err := easyjson.Marshal(tc.body)
				require.NoError(t, err, "error Marshall request body")
				req.Body = body
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			require.Equal(t, tc.want.statusCode, resp.StatusCode())
			require.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))

			if tc.want.body != "" {

				assert.JSONEq(t, tc.want.body, string(resp.Body()))
			}
		})

	}
}

func TestMetricHandler_GzipJson(t *testing.T) {
	server := httptest.NewServer(createTestServer())
	defer server.Close()

	jsonCounterValue := int64(5)
	data := models.Metrics{
		ID:    "CounterJson",
		MType: "counter",
		Delta: &jsonCounterValue}

	req := resty.New().R()
	req.Method = "POST"
	req.URL = server.URL + "/update"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	body, err := easyjson.Marshal(data)
	require.NoError(t, err, "error Marshall request body")

	gz, err := compress.GzipCompress(body)

	require.NoError(t, err, "error Compress request body")

	req.Body = gz
	_, err = req.Send()

	data2 := models.Metrics{
		ID:    "CounterJson",
		MType: "counter"}

	body2, err := easyjson.Marshal(data2)
	require.NoError(t, err, "error Marshall request body")

	req = resty.New().R()
	req.Method = "POST"
	req.URL = server.URL + "/value"
	req.Body = body2
	req.Header.Set("Accept-Encoding", "gzip")

	res2, err := req.Send()
	require.NoError(t, err, "error Marshall request body")

	require.JSONEq(t, string(body), string(res2.Body()))
	require.Equal(t, res2.Header().Get("Content-Encoding"), "gzip")
}
