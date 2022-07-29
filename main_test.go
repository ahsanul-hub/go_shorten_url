package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_showLinks(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/links", nil)
	response := httptest.NewRecorder()

	showLinks(response, request)

	t.Run("test links", func(t *testing.T) {
		got := response.Body.String()
		want := "0"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
