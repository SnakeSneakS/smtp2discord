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
		if username != Cfg.Auth.Username || password != Cfg.Auth.Password {
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
		var err error
		s.Email.Text = string(b)
		for _, sendEmailFunc := range s.SendEmailFuncs {
			if _err := sendEmailFunc(*s.Email); _err != nil {
				Cfg.Logger.Errorf("failed to send email: %v", _err)
				err = fmt.Errorf("%w, %w", err, _err)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) Reset() {
	s.Email = &EmailData{}
}

func (s *Session) Logout() error {
	//Logger.Println("logout")
	//Logger.Println(s)
	return nil
}

func NewServer(backend *Backend) *smtp.Server {
	server := smtp.NewServer(backend)
	server.Addr = Cfg.Server.Addr
	server.Domain = Cfg.Server.Domain
	server.WriteTimeout = Cfg.Server.WriteTimeout
	server.ReadTimeout = Cfg.Server.ReadTimeout
	server.MaxMessageBytes = int64(Cfg.Server.MaxMessageBytes)
	server.AllowInsecureAuth = Cfg.Server.AllowInsecureAuth
	//Logger.Println("Starting server at", server.Addr)
	//if err := server.ListenAndServe(); err != nil {
	//	Logger.Fatal(err)
	//}
	//defer server.Close()
	return server
}
