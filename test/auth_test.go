package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmdhs/authlib-skin/model"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	rep, err := http.Get("http://127.0.0.1:8081/api/v1/user")
	require.Nil(t, err)
	defer rep.Body.Close()

	require.Equal(t, rep.StatusCode, 401)

	reqs, err := http.NewRequest("GET", "http://127.0.0.1:8081/api/v1/user", nil)
	require.Nil(t, err)
	reqs.Header.Add("Authorization", "Bearer aaaaaaaaaaa")
	rep1, err := http.DefaultClient.Do(reqs)
	require.Nil(t, err)
	defer rep1.Body.Close()

	assert.Equal(t, rep.StatusCode, 401)

	var api model.API[any]
	require.Nil(t, json.NewDecoder(rep1.Body).Decode(&api))

	assert.Equal(t, int(api.Code), int(5))
}
