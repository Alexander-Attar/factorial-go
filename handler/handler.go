package handler

import (
	"io"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/reservemedia/factorial-go"
)

// Factorial Handler generates factorials for n requested via the `n=` query
// parameter.
func Factorial(w http.ResponseWriter, r *http.Request) {
	n, ok := new(big.Int).SetString(r.URL.Query().Get("n"), 10)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Invalid number")
		return
	}

	f := factorial.GenerateFactorial(n)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, f.String())
}

// Flake middleware returns 502s depending on the provided threshold value.
func Flake(next http.Handler, threshold float32) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if rand.Float32() < threshold {
			w.WriteHeader(http.StatusBadGateway)
			io.WriteString(w, "Bad gateway")
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
