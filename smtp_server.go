package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type EmailData struct {
	//email.Email
	From string
	To   []string
	Text string
}

type SendEmailFunc = func(email EmailData) error

type Backend struct {
	SendEmailFuncs []SendEmailFunc
}

func NewBackend() *Backend {
	return &Backend{}
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return NewSession(c, b.SendEmailFuncs), nil
}

type Session struct {
	Conn           *smtp.Conn
	Email          *EmailData
	SendEmailFuncs []SendEmailFunc
}

func NewSession(
	conn *smtp.Conn,
	sendEmailFuncs []SendEmailFunc,
) *Session {
	return &Session{
		Conn:           conn,
		Email:          &EmailData{},
		SendEmailFuncs: sendEmailFuncs,
	}
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != Config.Auth.Username || password != Config.Auth.Password {
			return errors.New("invalid username or password")
		}
		return nil
	}), nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	//Logger.Println("Mail from:", from)
	s.Email.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	//Logger.Println("Rcpt to:", to)
	s.Email.To = append(s.Email.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		return err
	} else {
		//Logger.Println("Data:", string(b))
		//s.Email.Text = append(s.Email.Text, []byte(b)...)
		s.Email.Text = string(b)
		for _, sendEmailFunc := range s.SendEmailFuncs {
			if err := sendEmailFunc(*s.Email); err != nil {

			}
		}
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	//Logger.Println("logout")
	//Logger.Println(s)
	return nil
}

func NewServer(backend *Backend) *smtp.Server {
	server := smtp.NewServer(backend)
	addr := fmt.Sprintf("localhost:%d", Config.Server.Port)
	server.Addr = addr
	server.Domain = Config.Server.Domain
	server.WriteTimeout = Config.Server.WriteTimeout
	server.ReadTimeout = Config.Server.ReadTimeout
	server.MaxMessageBytes = int64(Config.Server.MaxMessageBytes)
	server.AllowInsecureAuth = Config.Server.AllowInsecureAuth
	//Logger.Println("Starting server at", server.Addr)
	//if err := server.ListenAndServe(); err != nil {
	//	Logger.Fatal(err)
	//}
	//defer server.Close()
	return server
}
