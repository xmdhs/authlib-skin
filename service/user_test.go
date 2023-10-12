package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/migrate"
	"github.com/xmdhs/authlib-skin/model"
)

var webService *WebService

func TestMain(m *testing.M) {
	ctx := context.Background()

	clean := initWebService(ctx)
	code := m.Run()

	clean()

	os.Exit(code)
}

func initWebService(ctx context.Context) func() {
	c := lo.Must(ent.Open("mysql", "root:root@tcp(127.0.0.1)/test"))
	lo.Must0(c.Schema.Create(context.Background(), migrate.WithForeignKeys(false), migrate.WithDropIndex(true), migrate.WithDropColumn(true)))
	rsa4 := lo.Must(rsa.GenerateKey(rand.Reader, 4096))
	webService = NewWebService(config.Default(), c, &http.Client{}, cache.NewFastCache(100000), rsa4)

	return func() {
		c.User.Delete().Exec(ctx)
		c.Texture.Delete().Exec(ctx)
		c.UserProfile.Delete().Exec(ctx)
		c.UserTexture.Delete().Exec(ctx)
		c.UserToken.Delete().Exec(ctx)
	}
}

func TestWebService_Reg(t *testing.T) {
	ctx := context.Background()
	webService.config.MaxIpUser = 1
	type args struct {
		ctx      context.Context
		u        model.UserReg
		ipPrefix string
		ip       string
	}
	tests := []struct {
		name    string
		w       *WebService
		args    args
		wantErr bool
	}{
		{
			name: "1",
			w:    webService,
			args: args{
				ctx: ctx,
				u: model.UserReg{
					Email:        "1@xmdhs.com",
					Password:     "123456",
					Name:         "111",
					CaptchaToken: "",
				},
				ipPrefix: "127.0.0.0/24",
				ip:       "127.0.0.1",
			},
			wantErr: false,
		},
		{
			name: "email duplicate",
			w:    webService,
			args: args{
				ctx: ctx,
				u: model.UserReg{
					Email:        "1@xmdhs.com",
					Password:     "123456",
					Name:         "111",
					CaptchaToken: "",
				},
				ipPrefix: "127.0.0.0/24",
				ip:       "127.0.0.1",
			},
			wantErr: true,
		},
		{
			name: "name duplicate",
			w:    webService,
			args: args{
				ctx: ctx,
				u: model.UserReg{
					Email:        "2@xmdhs.com",
					Password:     "123456",
					Name:         "111",
					CaptchaToken: "",
				},
				ipPrefix: "127.0.0.0/24",
				ip:       "127.0.0.1",
			},
			wantErr: true,
		},
		{
			name: "MaxIpUser",
			w:    webService,
			args: args{
				ctx: ctx,
				u: model.UserReg{
					Email:        "3@xmdhs.com",
					Password:     "123456",
					Name:         "333",
					CaptchaToken: "",
				},
				ipPrefix: "127.0.0.0/24",
				ip:       "127.0.0.1",
			},
			wantErr: true,
		},
		{
			name: "MaxIpUser",
			w:    webService,
			args: args{
				ctx: ctx,
				u: model.UserReg{
					Email:        "4@xmdhs.com",
					Password:     "123456",
					Name:         "444",
					CaptchaToken: "",
				},
				ipPrefix: "127.0.0.2/24",
				ip:       "127.0.0.1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.Reg(tt.args.ctx, tt.args.u, tt.args.ipPrefix, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("WebService.Reg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	webService.config.MaxIpUser = 0
}
