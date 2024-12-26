package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/emersion/go-smtp"
)

func main() {
	// サーバーの起動
	server := NewSmtp2DiscordServer()
	Cfg.Logger.Debugf("listen and serve on %s", Cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		Cfg.Logger.Errorf("Failed to start server: %v", err)
	}
}

func NewSmtp2DiscordServer() *smtp.Server {
	backend := NewBackend()
	backend.SendEmailFuncs = append(backend.SendEmailFuncs, func(e EmailData) error {
		Cfg.Logger.Debugf("will send email data from(%s) to(%v). %s", e.From, e.To, e.Text)
		return nil
	})
	backend.SendEmailFuncs = append(backend.SendEmailFuncs, sendEmailDataToDiscord)
	server := NewServer(backend)
	return server
}

func sendEmailDataToDiscord(e EmailData) error {
	templateData := map[string]interface{}{
		"From": e.From,
		"To":   e.To,
		"Text": e.Text,
	}
	message, err := RenderDiscordMessageTemplate(templateData)
	if err != nil {
		return err
	}
	if err := SendToDiscord(message); err != nil {
		return err
	}
	return nil
}

func RenderDiscordMessageTemplate(data interface{}) (string, error) {
	templateContent := Cfg.Discord.MessageTemplate
	tmpl, err := template.New("discord message").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("error creating template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}
	return buf.String(), nil
}

func SendToDiscord(message string) error {
	payload := map[string]string{
		"content": message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}
	resp, err := http.Post(Cfg.Discord.WebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending request to Discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code from Discord: %d", resp.StatusCode)
	}
	return nil
}
