package yggdrasil

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	yggdrasilM "github.com/xmdhs/authlib-skin/model/yggdrasil"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

type PubRsaKey string

type Yggdrasil struct {
	logger           *slog.Logger
	validate         *validator.Validate
	yggdrasilService *yggdrasilS.Yggdrasil
	config           config.Config
	pubkey           PubRsaKey
}

func NewYggdrasil(logger *slog.Logger, validate *validator.Validate, yggdrasilService *yggdrasilS.Yggdrasil, config config.Config, pubkey PubRsaKey) *Yggdrasil {
	return &Yggdrasil{
		logger:           logger,
		validate:         validate,
		yggdrasilService: yggdrasilService,
		config:           config,
		pubkey:           pubkey,
	}
}

func getAnyModel[K any](ctx context.Context, w http.ResponseWriter, r io.Reader, validate *validator.Validate, slog *slog.Logger) (K, bool) {
	a, err := utils.DeCodeBody[K](r, validate)
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
		return a, false
	}
	return a, true
}

func (y *Yggdrasil) YggdrasilRoot() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var host string
		if y.config.TextureBaseUrl != "" {
			u := lo.Must(url.Parse(y.config.TextureBaseUrl))
			host = u.Hostname()
		} else {
			host, _ = lo.TryOr[string](func() (string, error) {
				h, _, err := net.SplitHostPort(r.Host)
				return h, err
			}, r.Host)
		}
		w.Write(lo.Must1(json.Marshal(yggdrasilM.Yggdrasil{
			Meta: yggdrasilM.YggdrasilMeta{
				ImplementationName:    "authlib-skin",
				ImplementationVersion: "0.0.1",
				Links: yggdrasilM.YggdrasilMetaLinks{
					Homepage: y.config.HomepageUrl,
					Register: y.config.RegisterUrl,
				},
				ServerName:       y.config.ServerName,
				EnableProfileKey: true,
			},
			SignaturePublickey: string(y.pubkey),
			SkinDomains:        []string{host},
		})))

	}
}

func (y *Yggdrasil) TextureAssets() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "image/png")
		http.StripPrefix("/texture/", http.FileServer(http.Dir(y.config.TexturePath))).ServeHTTP(w, r)
	}
}
