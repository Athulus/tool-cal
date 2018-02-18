package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func toolValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tool, exists := vars["tool"]
		if exists {
			_, ok := tools[tool]
			if !ok {
				http.Error(w, "a calendar for this tool does not exist", 400)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
