package service

import (
	"context"

	"github.com/xmdhs/authlib-skin/model"
)

func (w *WebService) GetCaptcha(ctx context.Context) model.Captcha {
	return model.Captcha{
		Type:    w.config.Captcha.Type,
		SiteKey: w.config.Captcha.SiteKey,
	}
}
