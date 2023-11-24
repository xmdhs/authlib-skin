package handle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/netip"

	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/model"
)

var HandelSet = wire.NewSet(NewUserHandel, NewAdminHandel, NewHandel)

func encodeJson[T any](w io.Writer, m model.API[T]) {
	json.NewEncoder(w).Encode(m)
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
