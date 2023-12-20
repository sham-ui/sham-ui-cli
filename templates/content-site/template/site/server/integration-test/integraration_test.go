package integration

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"site/config"
	"site/internal/external_api/cms/proto"
	"site/pkg/asserts"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func runCMS(t *testing.T, cfg *config.Config, m cmsServer) {
	t.Helper()
	var lis net.Listener
	var err error
	switch {
	case strings.HasPrefix(cfg.API.Address, "unix:"):
		asserts.NoError(t, os.MkdirAll(path.Dir(cfg.API.Address), os.ModePerm))
		addr, err := net.ResolveUnixAddr("unix", cfg.API.Address[5:])
		asserts.NoError(t, err)
		lis, err = net.ListenUnix("unix", addr)
		asserts.NoError(t, err)
		asserts.NoError(t, os.Chmod(cfg.API.Address, 0o777))
	case strings.HasPrefix(cfg.API.Address, "/"):
		asserts.NoError(t, os.MkdirAll(path.Dir(cfg.API.Address), os.ModePerm))
		addr, err := net.ResolveUnixAddr("unix", cfg.API.Address)
		asserts.NoError(t, err)
		lis, err = net.ListenUnix("unix", addr)
		asserts.NoError(t, err)
		asserts.NoError(t, os.Chmod(cfg.API.Address, 0o777))
	default:
		lis, err = net.Listen("tcp", cfg.API.Address)
	}
	asserts.NoError(t, err)
	asserts.NoError(t, err)
	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})
	//nolint:exhaustruct
	proto.RegisterCMSServer(srv, testCmsServer{m: m})
	go func() {
		_ = srv.Serve(lis)
	}()
	time.Sleep(500 * time.Millisecond)
}

func runApp(t *testing.T) *exec.Cmd {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("./app")
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	asserts.NoError(t, cmd.Start())
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		time.Sleep(500 * time.Millisecond)
	})
	time.Sleep(500 * time.Millisecond)
	return cmd
}

func TestSmoke(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfiguration(testr.New(t).V(1), "config.cfg")
	asserts.NoError(t, err)

	// Action
	runCMS(t, cfg, newMockCmsServer(t))
	cmd := runApp(t)
	stdout := (cmd.Stdout).(*bytes.Buffer).Bytes()

	// Assert
	asserts.Contains(t, string(stdout), "start ssr sub process")
	asserts.Contains(t, string(stdout), "server started")
}

func Test_Pages(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		CMS          func(t mockConstructorTestingTnewMockCmsServer) *mockCmsServer
		ExpectedCode int
		ExpectedBody string
	}{
		{
			Name: "home",
			URL:  "/",
			CMS: func(t mockConstructorTestingTnewMockCmsServer) *mockCmsServer {
				m := newMockCmsServer(t)
				m.
					On("ArticleList", mock.Anything, &proto.ArticleListRequest{
						Offset: 0,
						Limit:  9,
					}).
					Return(&proto.ArticleListResponse{
						Articles: []*proto.ArticleListItem{
							{
								Title: "hello world",
								Slug:  "hello-world",
								Category: &proto.Category{
									Name: "world category",
									Slug: "world-category",
								},
								Content: "Hello world!",
							},
						},
						Total: 1,
					}, nil).
					Once()
				return m
			},
			ExpectedCode: 200,
			ExpectedBody: `<a href="/category/world-category/1/" class="categorie">world category</a><h5><a href="/hello-world">hello world</a></h5><p>Hello world!</p>`,
		},
		{
			Name:         "contact",
			URL:          "/contact",
			CMS:          newMockCmsServer,
			ExpectedCode: 200,
			ExpectedBody: `img src="/images/logo-dark.png"`,
		},
		{
			Name:         "contact strip slash",
			URL:          "/contact/",
			CMS:          newMockCmsServer,
			ExpectedCode: 200,
			ExpectedBody: `img src="/images/logo-dark.png"`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			// Arrange
			cfg, err := config.LoadConfiguration(testr.New(t).V(1), "config.cfg")
			asserts.NoError(t, err)
			runCMS(t, cfg, test.CMS(t))
			cmd := runApp(t)
			stdout := cmd.Stdout.(*bytes.Buffer)
			stdout.Reset()

			// Action
			//nolint:bodyclose,noctx
			resp, err := (&http.Client{Timeout: 5 * time.Second}).Get(cfg.Server.URL() + test.URL)
			asserts.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			asserts.NoError(t, err)

			// Assert
			asserts.Equals(t, test.ExpectedCode, resp.StatusCode, "status code")
			asserts.Contains(t, string(b), test.ExpectedBody)
			asserts.Contains(t, stdout.String(), test.URL)
		})
	}
}
