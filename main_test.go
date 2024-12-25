package main_test

import (
	"net/smtp"
	"testing"
	"time"

	. "github.com/snakesneaks/smtp2discord"
)

func TestMain(t *testing.T) {
	server := NewSmtp2DiscordServer()
	// サーバーの起動
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !IsClosedNetworkError(err) {
				t.Errorf("Failed to start server: %v", err)
			}
		}
	}()

	// サーバーが起動するまで待機
	time.Sleep(100 * time.Millisecond)

	const sender = "sender@example.com"
	const recipient = "recipient@example.com"
	const message = "Subject: Test\nThis is a test email."
	auth := smtp.PlainAuth("", Config.Auth.Username, Config.Auth.Password, "localhost")
	wrongAuth := smtp.PlainAuth("", "fake username", "fake password", "localhost")
	err := smtp.SendMail(server.Addr, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
	mustErr := smtp.SendMail(server.Addr, wrongAuth, sender, []string{recipient}, []byte(message))
	if mustErr == nil {
		t.Fatalf("Success to send email with wrong authentication: %v", err)
	}
}
