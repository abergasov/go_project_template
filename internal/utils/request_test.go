package utils_test

import (
	"context"
	"fmt"
	"go_project_template/internal/utils"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

type TestResponse struct {
	ID1 uuid.UUID `json:"id"`
	ID2 uuid.UUID `json:"id_2"`
	ID3 uuid.UUID `json:"id_3"`
}

type TestRequest struct {
	ID1 uuid.UUID `json:"id"`
	ID2 uuid.UUID `json:"id_2"`
	ID3 uuid.UUID `json:"id_3"`
}

func TestGetCurl(t *testing.T) {
	// given
	header := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	sampleResult := TestResponse{
		ID1: uuid.New(),
		ID2: uuid.New(),
		ID3: uuid.New(),
	}
	headers := map[string]string{
		"Custom": header,
	}
	srv, address := newTestServer(t)
	srv.registerHandler(http.MethodGet, "/test", func(ctx *fiber.Ctx) error {
		require.Equal(t, []string{header}, ctx.GetReqHeaders()["Custom"])
		return ctx.JSON(sampleResult)
	})
	srv.Start(t)

	t.Run("should serve error in case of non existing address", func(t *testing.T) {
		// when
		res, code, err := utils.GetCurl[TestResponse](ctx, "http://127.0.0.1:1/test", headers)

		// then
		require.ErrorContains(t, err, "connect: connection refused")
		require.Zerof(t, code, "should return 0 code")
		require.Zero(t, res, "should return empty response")
	})
	t.Run("should serve correct request", func(t *testing.T) {
		// when
		res, code, err := utils.GetCurl[TestResponse](ctx, fmt.Sprintf("%s/test", address), headers)

		// then
		require.NoError(t, err, "should not return error")
		require.Equal(t, http.StatusOK, code, "should return 200 code")
		require.Equal(t, &sampleResult, res, "should return correct response")
	})
}

func TestPatchCurl(t *testing.T) {
	// given
	header := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	sampleResult := TestResponse{
		ID1: uuid.New(),
		ID2: uuid.New(),
		ID3: uuid.New(),
	}
	sampleRequest := TestRequest{
		ID1: uuid.New(),
		ID2: uuid.New(),
		ID3: uuid.New(),
	}
	headers := map[string]string{
		"Custom": header,
	}
	srv, address := newTestServer(t)
	srv.registerHandler(http.MethodPatch, "/test", func(ctx *fiber.Ctx) error {
		require.Equal(t, []string{header}, ctx.GetReqHeaders()["Custom"])
		var req TestRequest
		require.NoError(t, ctx.BodyParser(&req))
		require.Equal(t, sampleRequest, req)
		return ctx.JSON(sampleResult)
	})
	srv.Start(t)

	t.Run("should serve error in case of non existing address", func(t *testing.T) {
		// when
		res, code, err := utils.PatchCurl[TestResponse](ctx, "http://127.0.0.1:1/test", sampleRequest, headers)

		// then
		require.ErrorContains(t, err, "connect: connection refused")
		require.Zerof(t, code, "should return 0 code")
		require.Zero(t, res, "should return empty response")
	})
	t.Run("should serve correct request", func(t *testing.T) {
		// when
		res, code, err := utils.PatchCurl[TestResponse](ctx, fmt.Sprintf("%s/test", address), sampleRequest, headers)

		// then
		require.NoError(t, err, "should not return error")
		require.Equal(t, http.StatusOK, code, "should return 200 code")
		require.Equal(t, &sampleResult, res, "should return correct response")
	})
}

func TestPostCurl(t *testing.T) {
	// given
	header := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	sampleResult := TestResponse{
		ID1: uuid.New(),
		ID2: uuid.New(),
		ID3: uuid.New(),
	}
	sampleRequest := TestRequest{
		ID1: uuid.New(),
		ID2: uuid.New(),
		ID3: uuid.New(),
	}
	headers := map[string]string{
		"Custom": header,
	}
	srv, address := newTestServer(t)
	srv.registerHandler(http.MethodPost, "/test", func(ctx *fiber.Ctx) error {
		require.Equal(t, []string{header}, ctx.GetReqHeaders()["Custom"])
		var req TestRequest
		require.NoError(t, ctx.BodyParser(&req))
		require.Equal(t, sampleRequest, req)
		return ctx.JSON(sampleResult)
	})
	srv.Start(t)
	t.Run("should serve error in case of non existing address", func(t *testing.T) {
		// when
		res, code, err := utils.PostCurl[TestResponse](ctx, "http://127.0.0.1:1/test", sampleRequest, headers)

		// then
		require.ErrorContains(t, err, "connect: connection refused")
		require.Zerof(t, code, "should return 0 code")
		require.Zero(t, res, "should return empty response")
	})
	t.Run("should serve correct request", func(t *testing.T) {
		// when
		res, code, err := utils.PostCurl[TestResponse](ctx, fmt.Sprintf("%s/test", address), sampleRequest, headers)

		// then
		require.NoError(t, err, "should not return error")
		require.Equal(t, http.StatusOK, code, "should return 200 code")
		require.Equal(t, &sampleResult, res, "should return correct response")
	})
}

type testServer struct {
	address    string
	httpEngine *fiber.App
}

// newTestServer creates a new test server for testing http requests.
func newTestServer(t *testing.T) (fakeServer *testServer, address string) {
	appPort, err := freeport.GetFreePort()
	require.NoError(t, err, "failed to get free port for app")
	fakeServer = &testServer{
		address: fmt.Sprintf(":%d", appPort),
		httpEngine: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
	}
	fakeServer.httpEngine.Use(recover.New())
	return fakeServer, fmt.Sprintf("http://127.0.0.1:%d", appPort)
}

func (ts *testServer) registerHandler(method, path string, handler func(ctx *fiber.Ctx) error) {
	ts.httpEngine.Add(method, path, handler)
}

func (ts *testServer) Start(t *testing.T) {
	go func() {
		require.NoError(t, ts.httpEngine.Listen(ts.address))
	}()
	t.Cleanup(func() {
		require.NoError(t, ts.httpEngine.Shutdown())
	})
}
