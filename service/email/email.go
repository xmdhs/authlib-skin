package email

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/samber/lo"
	"github.com/wneessen/go-mail"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
)

type EmailConfig struct {
	Host string
	Port int
	SSL  bool
	Name string
	Pass string
}

type Email struct {
	emailConfig []EmailConfig
	pri         *rsa.PrivateKey
	config      config.Config
	cache       cache.Cache
}

func NewEmail(emailConfig []EmailConfig, pri *rsa.PrivateKey, config config.Config, cache cache.Cache) Email {
	return Email{
		emailConfig: emailConfig,
		pri:         pri,
		config:      config,
		cache:       cache,
	}
}

func (e Email) getRandEmailUser() EmailConfig {
	i := rand.Intn(len(e.emailConfig))
	return e.emailConfig[i]
}

func (e Email) SendEmail(ctx context.Context, to string, subject, body string) error {
	u := e.getRandEmailUser()
	m := mail.NewMsg()

	err := m.From(u.Name)
	if err != nil {
		return fmt.Errorf("SendRegVerify: %w", err)
	}

	err = m.To(to)
	if err != nil {
		return fmt.Errorf("SendRegVerify: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	c, err := mail.NewClient(u.Host, mail.WithPort(u.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(u.Name), mail.WithPassword(u.Pass))
	if err != nil {
		return fmt.Errorf("SendRegVerify: %w", err)
	}
	if u.SSL {
		c.SetSSL(true)
	}
	defer c.Close()

	err = c.DialAndSendWithContext(ctx, m)
	if err != nil {
		return fmt.Errorf("SendRegVerify: %w", err)
	}
	return nil
}

func (e Email) SendVerifyCode(ctx context.Context, email string, interval int) error {
	sendKey := []byte("SendEmail" + email)
	sendB, err := e.cache.Get(sendKey)
	if err != nil {
		return fmt.Errorf("SendRegVerifyCode: %w", err)
	}
	if sendB == nil {
		return fmt.Errorf("SendRegVerifyCode: %w", ErrSendLimit)
	}
	err = e.cache.Put(sendKey, []byte{1}, time.Now().Add(time.Second*time.Duration(interval)))
	if err != nil {
		return fmt.Errorf("SendRegVerifyCode: %w", err)
	}

	code := lo.RandomString(8, append(lo.NumbersCharset, lo.UpperCaseLettersCharset...))

	err = e.cache.Put([]byte("VerifyCode"+email), []byte(code), time.Now().Add(5*time.Minute))
	if err != nil {
		return fmt.Errorf("SendRegVerifyCode: %w", err)
	}
	err = e.SendEmail(ctx, email, "验证你的邮箱", fmt.Sprintf("验证码：%v，五分钟内有效", code))
	if err != nil {
		return fmt.Errorf("SendRegVerifyCode: %w", err)
	}
	return nil
}

var (
	ErrCodeNotValid = errors.New("验证码无效")
	ErrSendLimit    = errors.New("邮件发送限制")
)

func (e Email) VerifyCode(email, code string) error {
	key := []byte("VerifyCode" + email)
	codeb, err := e.cache.Get(key)
	if err != nil {
		return fmt.Errorf("VerifyCode: %w", err)
	}
	if string(codeb) != code {
		err := e.cache.Del(key)
		if err != nil {
			return fmt.Errorf("VerifyCode: %w", err)
		}
		return ErrCodeNotValid
	}
	return nil
}
