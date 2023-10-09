package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/xmdhs/authlib-skin/model"
)

func (w *WebService) GetConfig(ctx context.Context) model.Config {
	return model.Config{
		Captcha: model.Captcha{
			Type:    w.config.Captcha.Type,
			SiteKey: w.config.Captcha.SiteKey,
		},
		AllowChangeName: !w.config.OfflineUUID,
	}
}

type turnstileRet struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

var ErrCaptcha = errors.New("验证码错误")

type ErrTurnstile struct {
	ErrorCodes []string
}

func (e ErrTurnstile) Error() string {
	return strings.Join(e.ErrorCodes, " ")
}

func (w *WebService) verifyCaptcha(ctx context.Context, token, ip string) error {
	if w.config.Captcha.Type != "turnstile" {
		return nil
	}
	bw := &bytes.Buffer{}
	err := json.NewEncoder(bw).Encode(turnstileResponse{
		Secret:   w.config.Captcha.Secret,
		Response: token,
		Remoteip: ip,
	})
	if err != nil {
		return fmt.Errorf("verifyTurnstile: %w", err)
	}
	reqs, err := http.NewRequestWithContext(ctx, "POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", bw)
	if err != nil {
		return fmt.Errorf("verifyTurnstile: %w", err)
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("Content-Type", "application/json")
	rep, err := w.httpClient.Do(reqs)
	if err != nil {
		return fmt.Errorf("verifyTurnstile: %w", err)
	}
	defer rep.Body.Close()

	var t turnstileRet
	err = json.NewDecoder(rep.Body).Decode(&t)
	if err != nil {
		return fmt.Errorf("verifyTurnstile: %w", err)
	}

	if !t.Success {
		return fmt.Errorf("verifyTurnstile: %w", errors.Join(ErrTurnstile{
			ErrorCodes: t.ErrorCodes,
		}, ErrCaptcha))
	}
	return nil
}

type turnstileResponse struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	Remoteip string `json:"remoteip"`
}
