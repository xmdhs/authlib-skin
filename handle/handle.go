package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
)

type Handel struct {
	webService *service.WebService
	validate   *validator.Validate
	config     config.Config
	logger     *slog.Logger
}

func NewHandel(webService *service.WebService, validate *validator.Validate, config config.Config, logger *slog.Logger) *Handel {
	return &Handel{
		webService: webService,
		validate:   validate,
		config:     config,
		logger:     logger,
	}
}

func encodeJson[T any](w io.Writer, m model.API[T]) {
	json.NewEncoder(w).Encode(m)
}

func (h *Handel) getTokenbyAuthorization(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		h.logger.DebugContext(ctx, "缺少 Authorization")
		handleError(ctx, w, "缺少 Authorization", model.ErrAuth, 401)
		return ""
	}
	al := strings.Split(auth, " ")
	if len(al) != 2 || al[0] != "Bearer" {
		h.logger.DebugContext(ctx, "Authorization 格式错误")
		handleError(ctx, w, "Authorization 格式错误", model.ErrAuth, 401)
		return ""
	}
	return al[1]
}

func getPrefix(ip string) (string, error) {
	ipa, err := netip.ParseAddr(ip)
	if err != nil {
		return "", fmt.Errorf("getPrefix: %w", err)
	}
	if ipa.Is6() {
		return lo.Must1(ipa.Prefix(48)).String(), nil
	}
	return lo.Must1(ipa.Prefix(24)).String(), nil
}
