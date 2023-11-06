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

func (r *TestResponse) RequireStatus(t *testing.T, status int) {
	t.Helper()
	require.NotNil(t, r.Res, "response is nil")
	require.Equal(t, status, r.Res.StatusCode, "invalid response status code")
}

func (r *TestResponse) RequireOk(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusOK)
}

func (r *TestResponse) RequireCreated(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusCreated)
}

func (r *TestResponse) RequireNoContent(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusNoContent)
}

func (r *TestResponse) RequireUnauthorized(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusUnauthorized)
}

func (r *TestResponse) RequireForbidden(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusForbidden)
}

func (r *TestResponse) RequireConflict(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusConflict)
}

func (r *TestResponse) RequireBadRequest(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusBadRequest)
}

func (r *TestResponse) RequireNotFound(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusNotFound)
}

func (r *TestResponse) RequireRedirect(t *testing.T, path string) {
	t.Helper()
	r.RequireStatus(t, http.StatusFound)
	require.Equal(t, path, r.Res.Header.Get("Location"), "wrong redirect location")
}

func (r *TestResponse) RequireServerError(t *testing.T) {
	t.Helper()
	r.RequireStatus(t, http.StatusInternalServerError)
}
