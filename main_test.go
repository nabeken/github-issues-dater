package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestConvertRelativeDate(t *testing.T) {
	assert := assert.New(t)

	{
		now := time.Date(2015, time.January, 21, 12, 0, 0, 0, time.Local)
		assert.Equal("updated:<=2015-01-07", ConvertRelativeDate(now, "updated:within:2w"))
	}

	{
		now := time.Date(2015, time.January, 14, 12, 0, 0, 0, time.Local)
		assert.Equal("updated:<=2015-01-07", ConvertRelativeDate(now, "updated:within:1w"))
	}

	{
		now := time.Date(2015, time.January, 8, 12, 0, 0, 0, time.Local)
		assert.Equal("updated:<=2015-01-01", ConvertRelativeDate(now, "updated:within:1w"))
	}
}

func TestHandleGet(t *testing.T) {
	assert := assert.New(t)
	h := mux.NewRouter()
	Bind(h)

	request(func(rr *httptest.ResponseRecorder) {
		h.ServeHTTP(rr, mustRequest("GET", "/"))
		assert.Equal(http.StatusNotFound, rr.Code)
	})

	request(func(rr *httptest.ResponseRecorder) {
		h.ServeHTTP(rr, mustRequest("GET", "/nabeken/github-issues-dater"))
		assert.Equal(http.StatusNotFound, rr.Code)
	})

	request(func(rr *httptest.ResponseRecorder) {
		h.ServeHTTP(rr, mustRequest("GET", "/nabeken/github-issues-dater/issues"))
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	request(func(rr *httptest.ResponseRecorder) {
		h.ServeHTTP(rr, mustRequest("GET", "/nabeken/github-issues-dater/issues?q=is%3Aissue+is%3Aopen+updated%3Awithin%3A1w"))

		assert.Equal(http.StatusFound, rr.Code)
		assert.True(
			strings.HasPrefix(
				rr.Header().Get("Location"),
				githubURL+"/nabeken/github-issues-dater/issues?q=is%3Aissue+is%3Aopen+updated%3A%3C%3D",
			),
		)
	})
}

func mustRequest(method, path string) *http.Request {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func request(f func(rr *httptest.ResponseRecorder)) {
	f(httptest.NewRecorder())
}
