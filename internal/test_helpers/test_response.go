package testhelpers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestResponse struct {
	Res *http.Response
}

func (r *TestResponse) Response() *http.Response {
	return r.Res
}

func (r *TestResponse) RequireText(t *testing.T) string {
	t.Helper()
	data, err := io.ReadAll(r.Res.Body)
	require.NoError(t, err, "failed to read body as bytes")
	return string(data)
}

func (r *TestResponse) RequireUnmarshal(t *testing.T, dst interface{}) {
	t.Helper()
	err := json.NewDecoder(r.Res.Body).Decode(dst)
	require.NoError(t, err)
}

func (r *TestResponse) RequireStatus(t *testing.T, status int) *TestResponse {
	t.Helper()
	require.NotNil(t, r.Res, "response is nil")
	require.Equal(t, status, r.Res.StatusCode, "invalid response status code")
	return r
}

func (r *TestResponse) RequireOk(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusOK)
	return r
}

func (r *TestResponse) RequireCreated(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusCreated)
	return r
}

func (r *TestResponse) RequireNoContent(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusNoContent)
	return r
}

func (r *TestResponse) RequireUnauthorized(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusUnauthorized)
	return r
}

func (r *TestResponse) RequireForbidden(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusForbidden)
	return r
}

func (r *TestResponse) RequireConflict(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusConflict)
	return r
}

func (r *TestResponse) RequireBadRequest(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusBadRequest)
	return r
}

func (r *TestResponse) RequireNotFound(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusNotFound)
	return r
}

func (r *TestResponse) RequireRedirect(t *testing.T, path string) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusFound)
	require.Equal(t, path, r.Res.Header.Get("Location"), "wrong redirect location")
	return r
}

func (r *TestResponse) RequireServerError(t *testing.T) *TestResponse {
	t.Helper()
	r.RequireStatus(t, http.StatusInternalServerError)
	return r
}
