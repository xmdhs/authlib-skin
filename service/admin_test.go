package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmdhs/authlib-skin/model"
)

func TestWebService_Auth(t *testing.T) {
	ctx := context.Background()
	err := webService.Reg(ctx, model.UserReg{
		Email:        "TestWebService_Auth@xmdhs.com",
		Password:     "TestWebService_Auth",
		Name:         "TestWebService_Auth",
		CaptchaToken: "",
	}, "127.0.1.0/24", "127.0.1.0")
	require.Nil(t, err)

	l, err := webService.Login(ctx, model.Login{
		Email:        "TestWebService_Auth@xmdhs.com",
		Password:     "TestWebService_Auth",
		CaptchaToken: "",
	}, "0.0.0.0")
	require.Nil(t, err)

	token, err := webService.Auth(ctx, l.Token)
	require.Nil(t, err)

	assert.Equal(t, token.Subject, l.UUID)
	assert.Equal(t, token.Tid, "1")

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		w       *WebService
		args    args
		wantErr bool
	}{
		{
			name: "some string",
			w:    webService,
			args: args{
				ctx:   ctx,
				token: "123213",
			},
			wantErr: true,
		},
		{
			name: "valid jwt",
			w:    webService,
			args: args{
				ctx:   ctx,
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjpbeyJ0b29sdHQiOiJodHRwczovL3Rvb2x0dC5jb20ifV0sImlhdCI6MTY5NzEwMjMzOCwiZXhwIjoxNjk3MTI2Mzk5LCJhdWQiOiIiLCJpc3MiOiIiLCJzdWIiOiIifQ.JTQWl1PEX8u7PhVc4dTtv1DRS6e1PbMDZNWOAFJmVqE",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.w.Auth(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("WebService.Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
