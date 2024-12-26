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
	text := e.Text
	//extractTextFromEmailText(e.Text)

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

// ExtractTextFromEmailText extracts text from an HTML input or returns the raw text if the input is not valid HTML.
// TODO: implement this well
func ExtractTextFromEmailText(input string) string {
	// Try parsing the input as HTML
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		// If parsing fails, assume it's plain text
		return input
	}

	// Extract text from the HTML
	var buf bytes.Buffer
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			extractText(n.FirstChild)
		}
		if n.NextSibling != nil {
			extractText(n.NextSibling)
		}
	}
	extractText(doc)
	return strings.TrimSpace(buf.String())
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
