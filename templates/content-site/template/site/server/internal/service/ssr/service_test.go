package ssr

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"site/pkg/tracing"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace/noop"

	"site/pkg/asserts"

	"github.com/go-logr/logr/testr"
)

const echoScript = `process.on( 'message', function( msg ) {
	process.send( {
		id: msg.id,
		html: "api: " + msg.api + ", origin: " + msg.origin + ", url: " + msg.url + ", cookies: " + JSON.stringify( msg.cookies )
	} );
} );`

const errorScript = `process.on( 'message', function( msg ) {
	process.send( {
		id: msg.id,
		error: 'test error'
	} );
} );`

const exitOnMessageScript = `process.on( 'message', function( msg ) {
	if ( msg.url.includes( "/response" ) ) {
		process.send( {
			id: msg.id,
			html: "success"
		} )
		return
	}
	process.exit( 1 );
} );`

func TestService_Render(t *testing.T) {
	testCases := []struct {
		name             string
		script           string
		url              *url.URL
		cookies          []*http.Cookie
		expectedResponse string
		expectedError    error
	}{
		{
			name:   "echo",
			script: echoScript,
			url: func() *url.URL {
				u, _ := url.Parse("http://127.0.0.1:8080/echo?a=1")
				return u
			}(),
			cookies: []*http.Cookie{
				{
					Name:  "foo1",
					Value: "bar1",
				},
				{
					Name:  "foo2",
					Value: "bar2",
				},
			},
			expectedResponse: "api: http://127.0.0.1:9991, origin: http://127.0.0.1:8080, url: http://127.0.0.1:8080/echo?a=1, cookies: \"foo1=bar1; foo2=bar2\"",
		},
		{
			name:   "error",
			script: errorScript,
			url: func() *url.URL {
				u, _ := url.Parse("http://127.0.0.1:8080/error")
				return u
			}(),
			cookies: []*http.Cookie{
				{
					Name:  "foo1",
					Value: "bar1",
				},
				{
					Name:  "foo2",
					Value: "bar2",
				},
			},
			expectedError: errors.New("ssr respond with error: test error"),
		},
		{
			name:   "exit",
			script: exitOnMessageScript,
			url: func() *url.URL {
				u, _ := url.Parse("http://127.0.0.1:8080/exit")
				return u
			}(),
			cookies: []*http.Cookie{
				{
					Name:  "foo1",
					Value: "bar1",
				},
				{
					Name:  "foo2",
					Value: "bar2",
				},
			},
			expectedError: errors.New("context deadline exceeded"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			srv := New(
				testr.New(t).V(10),
				noop.NewTracerProvider(),
				tracing.NewPropagator(),
				"http://127.0.0.1:9991",
				[]byte(test.script),
			)
			asserts.NoError(t, srv.Start())
			t.Cleanup(srv.Stop)

			// Action
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			resp, err := srv.Render(ctx, test.url, test.cookies)

			// Assert
			asserts.Equals(t, test.expectedResponse, string(resp), "response")
			asserts.ErrorsEqual(t, test.expectedError, err)
		})
	}
}

func TestService_restart(t *testing.T) {
	srv := New(
		testr.New(t).V(10),
		noop.NewTracerProvider(),
		tracing.NewPropagator(),
		"http://127.0.0.1:9991",
		[]byte(exitOnMessageScript),
	)
	asserts.NoError(t, srv.Start())
	t.Cleanup(srv.Stop)

	// Fail response
	{
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		u, _ := url.Parse("http://127.0.0.1:8080")
		resp, err := srv.Render(ctx, u, nil)
		asserts.ErrorsEqual(t, errors.New("context deadline exceeded"), err)
		asserts.Equals(t, "", string(resp))
	}

	// Wait for restart
	time.Sleep(500 * time.Millisecond)

	// Success response after restart
	{
		u, _ := url.Parse("http://127.0.0.1:8080/response")
		resp, err := srv.Render(context.Background(), u, nil)
		asserts.NoError(t, err)
		asserts.Equals(t, "success", string(resp), "response")
	}
}
