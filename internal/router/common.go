package router

import "net/http"

type ErrorResponse func(w http.ResponseWriter, r *http.Request) error

func HandleErrors(h ErrorResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HomeRoute(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte("Hello world"))
	return err
}
