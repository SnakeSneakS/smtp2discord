package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/emersion/go-smtp"
	"golang.org/x/net/html"
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
	text, err := extractTextFromHTML(e.Text)
	if err != nil {
		Cfg.Logger.Debugf("failed to parse email text as html: %s", e.Text)
		text = e.Text
	}
	templateData := map[string]interface{}{
		"From": e.From,
		"To":   e.To,
		"Text": text,
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

func extractTextFromHTML(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}
	var result strings.Builder
	var extractText func(*html.Node)
	extractText = func(node *html.Node) {
		if node.Type == html.TextNode {
			result.WriteString(node.Data)
			result.WriteString("\n")
		} else if node.Type == html.ElementNode {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				extractText(child)
			}
		}
	}
	extractText(doc)
	return strings.TrimSpace(result.String()), nil
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
