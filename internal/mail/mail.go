// Package mail sends the few emails the server sends: for now, password reset
// links.
//
// With no SMTP host configured, a message is written to the log instead of
// being sent, so a developer can follow a reset link without a mail server.
package mail

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"xermess/internal/config"
)

// Message is one plain-text email.
type Message struct {
	To      string
	Subject string
	Body    string
}

// Sender sends a message. The server holds one; tests hand in their own to
// read what would have been sent.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// New returns the sender the configuration asks for.
func New(cfg config.Mail, log *slog.Logger) Sender {
	if cfg.Host == "" {
		return Log{log: log}
	}

	return &SMTP{cfg: cfg}
}

// Log writes messages to the log instead of sending them.
type Log struct {
	log *slog.Logger
}

// Send logs the message, body included: it is only used where no mail is
// configured, which is to say while developing.
func (l Log) Send(_ context.Context, msg Message) error {
	l.log.Info("email not sent: XERMESS_SMTP_HOST is not set", "to", msg.To, "subject", msg.Subject, "body", msg.Body)
	return nil
}

// SMTP sends through a mail server, with STARTTLS when the server offers it
// and plain authentication when a username is set.
type SMTP struct {
	cfg config.Mail
}

// Send sends one message.
func (s *SMTP) Send(ctx context.Context, msg Message) error {
	from, err := mail.ParseAddress(s.cfg.From)
	if err != nil {
		return fmt.Errorf("mail: XERMESS_SMTP_FROM: %w", err)
	}

	to, err := mail.ParseAddress(msg.To)
	if err != nil {
		return fmt.Errorf("mail: recipient: %w", err)
	}

	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	// net/smtp has no context; a deadline on the whole send stands in for it.
	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(addr, auth, from.Address, []string{to.Address}, compose(from, to, msg))
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// compose writes the message with the headers a mail client expects. The
// subject is encoded, since it may carry an application's name in any script.
func compose(from, to *mail.Address, msg Message) []byte {
	var b strings.Builder

	headers := [][2]string{
		{"From", from.String()},
		{"To", to.String()},
		{"Subject", mime.QEncoding.Encode("utf-8", msg.Subject)},
		{"Date", time.Now().Format(time.RFC1123Z)},
		{"MIME-Version", "1.0"},
		{"Content-Type", `text/plain; charset="utf-8"`},
		{"Content-Transfer-Encoding", "8bit"},
	}
	for _, header := range headers {
		fmt.Fprintf(&b, "%s: %s\r\n", header[0], header[1])
	}

	b.WriteString("\r\n")
	b.WriteString(strings.ReplaceAll(msg.Body, "\n", "\r\n"))

	return []byte(b.String())
}
