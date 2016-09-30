package handler_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reservemedia/factorial-go/handler"
)

func TestFactorialHandler(t *testing.T) {
	testCases := map[string]struct {
		Code int
		Body string
	}{
		// Valid requests
		"n=1": {http.StatusOK, "1"},
		"n=5": {http.StatusOK, "120"},

		// Invalid requests
		"n=": {http.StatusBadRequest, "Invalid number"},
		"n=some-invalid-number": {http.StatusBadRequest, "Invalid number"},
	}

	for name, expected := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			qs := fmt.Sprintf("/?%s", name)
			req, err := http.NewRequest("GET", qs, nil)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(handler.Factorial)
			handler.ServeHTTP(rec, req)

			if gotCode := rec.Code; gotCode != expected.Code {
				t.Errorf("handler returned invalid code: got %v, expected %v",
					gotCode,
					expected.Code,
				)
			}

			if gotBody := rec.Body; gotBody.String() != expected.Body {
				t.Errorf("handler returned invalid body: got %v, expected %v",
					gotBody,
					expected.Body,
				)
			}
		})
	}
}

func TestFlakeMiddleware(t *testing.T) {
	rand.Seed(1234567890)

	testCases := []struct {
		Pct       float32
		Threshold float32
		Code      int
	}{
		// rand32() calls based off the seed above
		{0.082695216, 0.1, http.StatusBadGateway},
		{0.7527473, 0.1, http.StatusOK},
		{0.045034427, 0.1, http.StatusBadGateway},
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("rand_%v", c.Pct), func(t *testing.T) {
			threshold := c.Threshold
			expectedCode := c.Code

			req, err := http.NewRequest("GET", "?n=1", nil)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Factorial)
			flake := handler.Flake(h, threshold)
			flake.ServeHTTP(rec, req)

			if gotCode := rec.Code; gotCode != expectedCode {
				t.Errorf("handler returned invalid code: got %v, expected %v",
					gotCode,
					expectedCode,
				)
			}
		})
	}

}
