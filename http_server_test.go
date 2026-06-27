package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/snakesneaks/smtp2discord"
)

func TestHTTPServer(t *testing.T) {
	called := false

	backend := NewBackend()
	backend.SendEmailFuncs = append(backend.SendEmailFuncs, func(e EmailData) error {
		called = true

		if e.From != "sender@example.com" {
			t.Fatalf("unexpected from: %s", e.From)
		}

		if len(e.To) != 1 || e.To[0] != "recipient@example.com" {
			t.Fatalf("unexpected recipient: %+v", e.To)
		}

		if e.Text != "hello from http" {
			t.Fatalf("unexpected text: %s", e.Text)
		}

		return nil
	})

	server := httptest.NewServer(NewHTTPHandler(backend))
	defer server.Close()

	body := EmailData{
		From: "sender@example.com",
		To:   []string{"recipient@example.com"},
		Text: "hello from http",
	}

	payload, _ := json.Marshal(body)

	req, err := http.NewRequest(
		http.MethodPost,
		server.URL+"/send",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(Cfg.Auth.Username, Cfg.Auth.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	if !called {
		t.Fatal("SendEmailFunc was not called")
	}
}