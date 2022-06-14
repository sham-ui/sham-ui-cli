package client

import (
	"bytes"
	"cms/test_helpers/asserts"
	"encoding/json"
	"fmt"
	"github.com/urfave/negroni"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ApiClient struct {
	n         *negroni.Negroni
	cookies   map[string]struct{}
	csrfToken string
}

func (c *ApiClient) Request(method string, url string, body interface{}) *ResponseWrapper {
	var bodyData io.Reader
	if nil != body {
		payload, err := json.Marshal(body)
		if nil != err {
			log.Fatalf("json marshal: %s", err)
		}
		bodyData = bytes.NewBuffer(payload)
	}
	req, err := http.NewRequest(method, url, bodyData)
	if nil != err {
		log.Fatalf("new request: %s", err)
	}
	c.setupCookies(req)
	req.Header.Set("X-Csrf-Token", c.csrfToken)
	responseRecorder := httptest.NewRecorder()
	c.n.ServeHTTP(responseRecorder, req)
	c.saveCookies(responseRecorder)

	return newResponseWrapper(responseRecorder)
}

func (c *ApiClient) GetCSRF() {
	resp := c.Request("GET", "/api/csrftoken", nil)
	c.csrfToken = resp.Response.Header().Get("X-Csrf-Token")
}

func (c *ApiClient) Login() {
	resp := c.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	if resp.Response.Code != http.StatusOK {
		log.Fatal("login response not OK")
	}
}

func (c *ApiClient) ResetCSRF() {
	c.csrfToken = ""
}

func (c *ApiClient) ResetCookies() {
	c.cookies = map[string]struct{}{}
}

func (c *ApiClient) saveCookies(resp *httptest.ResponseRecorder) {
	chunks := strings.Split(resp.Header().Get("Set-Cookie"), "; ")
	for _, chunk := range chunks {
		if "" != chunk {
			c.cookies[chunk] = struct{}{}
		}
	}
}

func (c *ApiClient) setupCookies(req *http.Request) {
	const separator = "; "
	var chunks []string
	for chunk := range c.cookies {
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}
	}
	req.Header.Set("Cookie", strings.Join(chunks, separator))
}

func (c *ApiClient) ExecuteTestCases(t *testing.T, testCases []TestCase) {
	for _, test := range testCases {
		resp := c.Request(test.Method, test.URL, test.Data)
		nameChunks := []string{test.Method, test.URL}
		if "" != test.Message {
			nameChunks = append(nameChunks, fmt.Sprintf("(%s)", test.Message))
		}
		testName := strings.Join(nameChunks, " ")
		asserts.Equals(t, test.ExpectedResponseStatusCode, resp.Response.Code, "code for "+testName)
		asserts.JSONEqualsWithoutSomeKeys(t, test.IgnoreKeys, test.ExpectedResponseJSON, resp.Text(), "body for "+testName)
	}
}

func NewApiClient(n *negroni.Negroni) *ApiClient {
	return &ApiClient{
		n:       n,
		cookies: map[string]struct{}{},
	}
}
