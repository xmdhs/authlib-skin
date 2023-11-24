package email

import (
	"bytes"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type EmailService struct {
	emailConfig []EmailConfig
	pri         *rsa.PrivateKey
	config      config.Config
	cache       cache.Cache
}

func NewEmail(pri *rsa.PrivateKey, c config.Config, cache cache.Cache) (*EmailService, error) {
	ec := lo.Map[config.SmtpUser, EmailConfig](c.Email.Smtp, func(item config.SmtpUser, index int) EmailConfig {
		return EmailConfig{
			Host: item.Host,
			Port: item.Port,
			SSL:  item.SSL,
			Name: item.Name,
			Pass: item.Pass,
		}
	})

	return &EmailService{
		emailConfig: ec,
		pri:         pri,
		config:      c,
		cache:       cache,
	}, nil
}

func (e EmailService) getRandEmailUser() (EmailConfig, error) {
	if len(e.emailConfig) == 0 {
		return EmailConfig{}, fmt.Errorf("没有可用的邮箱账号")
	}

	i := rand.Intn(len(e.emailConfig))
	return e.emailConfig[i], nil
}

func (e EmailService) SendEmail(ctx context.Context, to string, subject, body string) error {
	u, err := e.getRandEmailUser()
	if err != nil {
		return fmt.Errorf("SendRegVerify: %w", err)
	}
	m := mail.NewMsg()

	err = m.From(u.Name)
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

var emailTemplate = lo.Must(template.New("email").Parse(`<p>{{ .msg }}</p><a href="{{.url}}">{{ .url }}</a>`))

func (e EmailService) SendVerifyUrl(ctx context.Context, email string, interval int, host string) error {
	sendKey := []byte("SendEmail" + email)
	sendB, err := e.cache.Get(sendKey)
	if err != nil {
		return fmt.Errorf("SendVerifyUrl: %w", err)
	}
	if sendB != nil {
		return fmt.Errorf("SendVerifyUrl: %w", ErrSendLimit)
	}
	err = e.cache.Put(sendKey, []byte{1}, time.Now().Add(time.Second*time.Duration(interval)))
	if err != nil {
		return fmt.Errorf("SendVerifyUrl: %w", err)
	}

	code, err := newJwtToken(e.pri, email)
	if err != nil {
		return fmt.Errorf("SendVerifyUrl: %w", err)
	}

	q := url.Values{}
	q.Set("code", code)
	q.Set("email", email)

	u := url.URL{
		Host:   host,
		Scheme: "http",
		Path:   "/register",
	}
	u.RawQuery = q.Encode()

	if e.config.WebBaseUrl != "" {
		webBase, err := url.Parse(e.config.WebBaseUrl)
		if err != nil {
			return fmt.Errorf("SendVerifyUrl: %w", err)
		}
		u.Host = webBase.Host
		u.Scheme = webBase.Scheme
	}

	body := bytes.NewBuffer(nil)
	err = emailTemplate.Execute(body, map[string]any{
		"msg": "点击下方链接验证你的邮箱，1 天内有效",
		"url": u.String(),
	})
	if err != nil {
		return fmt.Errorf("SendVerifyUrl: %w", err)
	}

	err = e.SendEmail(ctx, email, "验证你的邮箱", body.String())
	if err != nil {
		return fmt.Errorf("SendVerifyUrl: %w", err)
	}
	return nil
}

var (
	ErrSendLimit    = errors.New("邮件发送限制")
	ErrTokenInvalid = errors.New("token 无效")
)

func (e EmailService) VerifyJwt(email, jwtStr string) error {
	token, err := jwt.ParseWithClaims(jwtStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return &e.pri.PublicKey, nil
	})
	if err != nil {
		return fmt.Errorf("VerifyJwt: %w", err)
	}
	sub, _ := token.Claims.GetSubject()
	if !token.Valid || sub != email {
		return fmt.Errorf("VerifyJwt: %w", ErrTokenInvalid)
	}
	return nil
}

func newJwtToken(jwtKey *rsa.PrivateKey, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   email,
		Issuer:    "authlib-skin email verification",
	})
	jwts, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("newJwtToken: %w", err)
	}
	return jwts, nil
}
