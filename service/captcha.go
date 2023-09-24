package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/xmdhs/authlib-skin/model"
)

func (w *WebService) GetCaptcha(ctx context.Context) model.Captcha {
	return model.Captcha{
		Type:    w.config.Captcha.Type,
		SiteKey: w.config.Captcha.SiteKey,
	}
}

type turnstileRet struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

type ErrTurnstile struct {
	ErrorCodes []string
}

func (e ErrTurnstile) Error() string {
	return strings.Join(e.ErrorCodes, " ")
}

func (w *WebService) verifyTurnstile(ctx context.Context, token, ip string) error {
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
		return fmt.Errorf("verifyTurnstile: %w", ErrTurnstile{
			ErrorCodes: t.ErrorCodes,
		})
	}
	return nil
}

type turnstileResponse struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	Remoteip string `json:"remoteip"`
}
