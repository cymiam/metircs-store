package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cymiam/metircs-store/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMetricHandler_HandleGet(t *testing.T) {
	ts := httptest.NewServer(handler.MetricRouter())
	defer ts.Close()

	testRequest(t, ts, "POST", "/update/counter/testCounter/1")
	testRequest(t, ts, "POST", "/update/gauge/testGauge/3.14")

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
		get  string
	}{
		{
			name: "Test Positive Get Gauge",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/value/gauge/testGauge",
			get: "3.140000",
		},
		{
			name: "Test Positive Get Counter",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/value/counter/testCounter",
			get: "[1]",
		},
		{
			name: "Test Get unknow counter",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/value/counter/unknownCounter",
			get: "",
		},
		{
			name: "Test Get unknow gauge",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/value/gauge/unknownGauge",
			get: "",
		},
	}
	for _, v := range tests {
		resp, get := testRequest(t, ts, "GET", v.url)
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
		assert.Equal(t, v.get, get)
	}
}

func TestMetricHandler_HandleUpdate(t *testing.T) {
	ts := httptest.NewServer(handler.MetricRouter())
	defer ts.Close()

	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Test Positive Gauge",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/gauge/test/3.0",
		},
		{
			name: "Test Positive Counter",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/test2/5",
		},
		{
			name: "Test Unknown metric type",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/timeseries/requests/5",
		},
		{
			name: "Test no metric name",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/5",
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, "POST", v.url)
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}
