package testhelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_project_template/internal/logger"
	"go_project_template/internal/routes"
	"net/http"
	"testing"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	appPort  int
	client   http.Client
	authUser string
}

func NewTestServer(t *testing.T, container *TestContainer) *TestServer {
	appPort, err := freeport.GetFreePort()
	require.NoError(t, err, "failed to get free port for app")
	srv := &TestServer{
		appPort: appPort,
		client:  *http.DefaultClient,
	}

	appLog := logger.NewAppSLogger("")
	appHTTPServer := routes.InitAppRouter(appLog, container.ServiceSampler, fmt.Sprintf(":%d", srv.appPort))
	t.Cleanup(func() {
		require.NoError(t, appHTTPServer.Stop())
	})
	go appHTTPServer.Run()
	return srv
}

func (ts *TestServer) AuthUser(mail string) {
	ts.authUser = mail
}

func (ts *TestServer) Get(t *testing.T, path string) *TestResponse {
	t.Helper()
	return ts.Request(t, http.MethodGet, path, nil, nil)
}

func (ts *TestServer) Post(t *testing.T, path string, body any) *TestResponse {
	t.Helper()
	return ts.Request(t, http.MethodPost, path, body, nil)
}

func (ts *TestServer) Put(t *testing.T, path string, body any) *TestResponse {
	t.Helper()
	return ts.Request(t, http.MethodPut, path, body, nil)
}

func (ts *TestServer) Delete(t *testing.T, path string, body any) *TestResponse {
	t.Helper()
	return ts.Request(t, http.MethodDelete, path, body, nil)
}

func (ts *TestServer) Request(t *testing.T, method string, path string, body interface{}, headers map[string]string) *TestResponse {
	t.Helper()

	var b []byte
	var err error
	if body != nil {
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["Content-Type"] = "application/json"
		b, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}

	u := fmt.Sprintf("http://localhost:%d%s", ts.appPort, path)
	req, err := http.NewRequest(method, u, bytes.NewBuffer(b))
	require.NoError(t, err, "failed to construct new request for url %s: %s", u, err)
	if err != nil {
		t.Fatal(err)
	}

	if len(headers) > 0 {
		for headerKey, headerVal := range headers {
			req.Header.Add(headerKey, headerVal)
		}
	}
	if ts.authUser != "" {
		req.Header.Add("Authorization", fmt.Sprint("Bearer ", ts.CreateToken(t, ts.authUser)))
	}

	res, err := ts.client.Do(req)
	require.NoError(t, err, "failed to make request to %s: %s", u, err)
	t.Cleanup(func() {
		require.NoError(t, res.Body.Close())
	})
	return &TestResponse{Res: res}
}

func (ts *TestServer) CreateToken(_ *testing.T, _ string) string {
	return "test_token"
}
