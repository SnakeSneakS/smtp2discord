package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func NewHTTPServer(backend *Backend) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: NewHTTPHandler(backend),
	}
}

func NewHTTPHandler(backend *Backend) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/send", BasicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleSend(backend, w, r)
	})))

	return mux
}

func Authenticate(username, password string) error {
	if username != Cfg.Auth.Username || password != Cfg.Auth.Password {
		return errors.New("invalid username or password")
	}
	return nil
}

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="smtp2discord"`)
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}

		if err := Authenticate(username, password); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleSend(backend *Backend, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailData EmailData
	if err := json.NewDecoder(r.Body).Decode(&emailData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, fn := range backend.SendEmailFuncs {
		if err := fn(emailData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}