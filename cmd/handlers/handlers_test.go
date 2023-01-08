package handlers_test

import (
	"github.com/LorezV/url-shorter.git/cmd/handlers"
	"github.com/LorezV/url-shorter.git/cmd/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLHandler(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name string
		urls []storage.URL
		path string
		want want
	}{
		{
			name: "Test GET request with exiting url in repository.",
			urls: []storage.URL{
				{
					ID:       "xhxKQF",
					Original: "https://practicum.yandex.ru",
					Short:    "http://127.0.0.1:8080/xhxKQF",
				},
			},
			path: "/xhxKQF",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru",
			},
		},
		{
			name: "Test GET request with empty repository.",
			urls: []storage.URL{},
			path: "/xhxKQF",
			want: want{
				statusCode: http.StatusNotFound,
				location:   "",
			},
		},
		{
			name: "Test GET request with different urls in the request and repository.",
			urls: []storage.URL{
				{
					ID:       "ASKTTG",
					Original: "https://practicum.yandex.ru",
					Short:    "http://127.0.0.1:8080/ASKTTG",
				},
			},
			path: "/xhxKQF",
			want: want{
				statusCode: http.StatusNotFound,
				location:   "",
			},
		},
		//{
		//	name:    "Test POST path.",
		//	urls:    []storage.URL{},
		//	path: "http://127.0.0.1:8080/",
		//	method:  http.MethodPost,
		//	body:    "https://practicum.yandex.ru",
		//	want: want{
		//		statusCode: http.StatusCreated,
		//		location:   "",
		//	},
		//},
		//{
		//	name:    "Test POST path with empty body.",
		//	urls:    []storage.URL{},
		//	path: "http://127.0.0.1:8080/",
		//	method:  http.MethodPost,
		//	body:    "",
		//	want: want{
		//		statusCode: http.StatusBadRequest,
		//		location:   "",
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.Repository = storage.MakeRepository()
			for _, url := range tt.urls {
				storage.Repository.Add(url)
			}

			r := chi.NewRouter()
			r.Get("/{id}", handlers.GetURL)
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, _ := testRequest(t, ts, http.MethodGet, tt.path, nil)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
		})
	}
}

func TestCreateURL(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name string
		path string
		body string
		want want
	}{
		{
			name: "Test POST request.",
			path: "/",
			body: "https://practicum.yandex.ru",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "Test POST path with empty body.",
			path: "/",
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/", handlers.CreateURL)
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, _ := testRequest(t, ts, http.MethodPost, tt.path, strings.NewReader(tt.body))

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
