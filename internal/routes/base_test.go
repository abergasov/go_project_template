package routes_test

import (
	testhelpers "go_project_template/internal/test_helpers"
	"testing"
)

func TestCheckPing(t *testing.T) {
	// given
	container := testhelpers.GetClean(t)
	srv := testhelpers.NewTestServer(t, container)

	t.Run("200 on /", func(t *testing.T) {
		// when, then
		srv.Get(t, "/").RequireOk(t)
	})
	t.Run("404 on unknown path", func(t *testing.T) {
		// when, then
		srv.Get(t, "/unknown").RequireStatus(t, 404)
	})
}
