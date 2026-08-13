package main

import (
	"booking/internal/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// mailConfig — настройки SMTP, прочитанные из окружения (.env читается в main).
type mailConfig struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Encryption     mail.Encryption
	From           string
	TemplatesDir   string
	KeepAlive      bool
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
}

func newMailConfig() mailConfig {
	return mailConfig{
		Host:           getEnv("MAIL_HOST", "localhost"),
		Port:           getEnvInt("MAIL_PORT", 1025),
		Username:       getEnv("MAIL_USERNAME", ""),
		Password:       getEnv("MAIL_PASSWORD", ""),
		Encryption:     parseEncryption(getEnv("MAIL_ENCRYPTION", "none")),
		From:           getEnv("MAIL_FROM", ""),
		TemplatesDir:   getEnv("MAIL_TEMPLATES_DIR", "./email-templates"),
		KeepAlive:      getEnvBool("MAIL_KEEP_ALIVE", false),
		ConnectTimeout: time.Duration(getEnvInt("MAIL_CONNECT_TIMEOUT", 10)) * time.Second,
		SendTimeout:    time.Duration(getEnvInt("MAIL_SEND_TIMEOUT", 10)) * time.Second,
	}
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		errorLog.Printf("%s=%q is not a number, using %d", key, raw, fallback)
		return fallback
	}

	return value
}

func getEnvBool(key string, fallback bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		errorLog.Printf("%s=%q is not a boolean, using %t", key, raw, fallback)
		return fallback
	}

	return value
}

// parseEncryption переводит значение MAIL_ENCRYPTION в константу библиотеки.
func parseEncryption(name string) mail.Encryption {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "none":
		return mail.EncryptionNone
	case "tls", "ssl", "ssltls", "ssl/tls":
		return mail.EncryptionSSLTLS
	case "starttls":
		return mail.EncryptionSTARTTLS
	default:
		errorLog.Printf("unknown MAIL_ENCRYPTION=%q, falling back to none", name)
		return mail.EncryptionNone
	}
}

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			// паника в этой горутине убила бы весь сервер, а не только письмо
			func() {
				defer func() {
					if r := recover(); r != nil {
						errorLog.Println("recovered from panic while sending mail:", r)
					}
				}()
				sendMsg(msg)
			}()
		}
	}()
}

func sendMsg(m models.MailData) {
	cfg := newMailConfig()

	server := mail.NewSMTPClient()
	server.Host = cfg.Host
	server.Port = cfg.Port
	server.Username = cfg.Username
	server.Password = cfg.Password
	server.Encryption = cfg.Encryption
	server.KeepAlive = cfg.KeepAlive
	server.ConnectTimeout = cfg.ConnectTimeout
	server.SendTimeout = cfg.SendTimeout

	client, err := server.Connect()
	if err != nil {
		// без соединения слать нечего: client == nil, и email.Send(client) уронил бы процесс
		errorLog.Printf("can't connect to the mail server %s:%d: %s", cfg.Host, cfg.Port, err)
		return
	}

	from := m.From
	if from == "" {
		from = cfg.From
	}

	email := mail.NewMSG()
	email.SetFrom(from).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		path := filepath.Join(cfg.TemplatesDir, m.Template)
		data, err := os.ReadFile(path)
		if err != nil {
			app.ErrorLog.Printf("can't read the mail template %s: %s", path, err)
			email.SetBody(mail.TextHTML, m.Content)
		} else {
			msgToSend := strings.Replace(string(data), "[%body%]", m.Content, 1)
			email.SetBody(mail.TextHTML, msgToSend)
		}
	}

	if err = email.Send(client); err != nil {
		errorLog.Printf("can't send the mail to %s: %s", m.To, err)
		return
	}

	infoLog.Printf("mail sent to %s: %s", m.To, m.Subject)
}
