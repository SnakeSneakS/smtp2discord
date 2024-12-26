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
	serverURL := fmt.Sprintf("%s:%s", Config.Server.Domain, strings.Split(Config.Server.Addr, ":")[1])
	auth := smtp.PlainAuth("", Config.Auth.Username, Config.Auth.Password, "localhost")
	err := smtp.SendMail(serverURL, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}
