package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/jordan-wright/email"
)

func main() {
	backend := NewBackend()
	backend.SendEmailFuncs = append(
		backend.SendEmailFuncs,
		sendEmailDataToDiscord,
	)

	smtpServer := NewServer(backend)
	httpServer := NewHTTPServer(backend)

	go func() {
		Cfg.Logger.Infof("SMTP listen on %s", Cfg.Server.Addr)
		if err := smtpServer.ListenAndServe(); err != nil {
			Cfg.Logger.Errorf("SMTP error: %v", err)
		}
	}()

	Cfg.Logger.Infof("HTTP listen on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		Cfg.Logger.Errorf("HTTP error: %v", err)
	}
}

func sendEmailDataToDiscord(e EmailData) error {
	templateData := map[string]interface{}{
		"From":    e.From,
		"To":      e.To,
		"Cc":      []string{},
		"Bcc":     []string{},
		"Subject": "",
		"Text":    e.Text,
		"HTML":    "",
	}

	emailExtracted, err := ExtractTextFromEmailText(e.Text)
	if err == nil {
		templateData["From"] = emailExtracted.From
		templateData["To"] = emailExtracted.To
		templateData["Cc"] = emailExtracted.Cc
		templateData["Bcc"] = emailExtracted.Bcc
		templateData["Subject"] = emailExtracted.Subject
		templateData["Text"] = string(emailExtracted.Text)
		templateData["HTML"] = string(emailExtracted.HTML)
	} else {
		Cfg.Logger.Errorf("failed to parse email content: %v", err)
	}

	message, err := RenderDiscordMessageTemplate(templateData)
	if err != nil {
		return err
	}

	messages := TruncateAndSplit(message, Cfg.Discord.DiscordMsgSizeMax)
	for _, messageSplit := range messages {
		if err := SendToDiscord(messageSplit); err != nil {
			return err
		}
	}

	return nil
}

func ExtractTextFromEmailText(input string) (*email.Email, error) {
	emailParsed, err := email.NewEmailFromReader(strings.NewReader(input))
	if err != nil {
		Cfg.Logger.Warnf("failed to decode email data, so return plain text")
		return nil, err
	}
	return emailParsed, nil
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

	resp, err := http.Post(
		Cfg.Discord.WebhookURL,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return fmt.Errorf("error sending request to Discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code from Discord: %d", resp.StatusCode)
	}

	return nil
}