package main_test

import (
	"net/smtp"
	"testing"
	"time"

	. "github.com/snakesneaks/smtp2discord"
)

func TestSMTPServer(t *testing.T) {
	backend := NewBackend()
	backend.SendEmailFuncs = append(backend.SendEmailFuncs, func(e EmailData) error {
		t.Logf("will send email: %+v", e)
		return nil
	})
	server := NewServer(backend)
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
	err := smtp.SendMail(server.Addr, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	/* # manual
	// SMTPクライアントで接続
	client, err := smtp.Dial(server.Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()
	// EHLOを送信
	if err := client.Hello("localhost"); err != nil {
		t.Fatalf("EHLO failed: %v", err)
	}
	// MAILコマンド
	if err := client.Mail("sender@example.com", nil); err != nil {
		t.Fatalf("MAIL command failed: %v", err)
	}
	// RCPTコマンド
	if err := client.Rcpt("recipient@example.com", nil); err != nil {
		t.Fatalf("RCPT command failed: %v", err)
	}
	// DATAコマンド
	wc, err := client.Data()
	if err != nil {
		t.Fatalf("DATA command failed: %v", err)
	}
	_, err = wc.Write([]byte("Subject: Test\n\nThis is a test email."))
	if err != nil {
		t.Fatalf("Failed to write email data: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("Failed to close DATA: %v", err)
	}
	// QUITコマンド
	if err := client.Quit(); err != nil {
		t.Fatalf("QUIT command failed: %v", err)
	}
	*/
}
