package middleware

import "net/http"

type JSONContent struct{}

func NewJSONContent() *JSONContent {
	return &JSONContent{}
}

func (j *JSONContent) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
