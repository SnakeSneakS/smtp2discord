package main_test

import (
	"fmt"
	"net/smtp"
	"strings"
	"testing"
	"time"

	. "github.com/snakesneaks/smtp2discord"
)

func TestMain(t *testing.T) {
	backend := NewBackend()

	called := false
	backend.SendEmailFuncs = append(
		backend.SendEmailFuncs,
		func(e EmailData) error {
			called = true
			return nil
		},
	)

	server := NewServer(backend)
	defer server.Close()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !IsClosedNetworkError(err) {
				t.Errorf("Failed to start server: %v", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	const sender = "sender@example.com"
	const recipient = "recipient@example.com"

	const message = `Subject: SUBJECT

test mail`

	serverURL := fmt.Sprintf(
		"%s:%s",
		Cfg.Server.Domain,
		strings.Split(Cfg.Server.Addr, ":")[1],
	)

	auth := smtp.PlainAuth(
		"",
		Cfg.Auth.Username,
		Cfg.Auth.Password,
		"localhost",
	)

	err := smtp.SendMail(
		serverURL,
		auth,
		sender,
		[]string{recipient},
		[]byte(message),
	)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	if !called {
		t.Fatal("SendEmailFunc was not called")
	}
}