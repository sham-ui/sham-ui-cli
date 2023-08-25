package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/negroni"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"{{shortName}}/test_helpers/asserts"
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
	return c.makeRequest(req)
}

func (c *ApiClient) RequestMultiPart(method string, url string, r *MulipartRequest) *ResponseWrapper {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(r.FileField, r.FilePath)
	if nil != err {
		log.Fatalf("create form file: %s", err)
	}
	_, err = io.Copy(part, bytes.NewBuffer(r.Content))
	if nil != err {
		log.Fatalf("init buffer: %s", err)
	}
	err = writer.Close()
	if nil != err {
		log.Fatalf("close writer: %s", err)
	}
	req, err := http.NewRequest(method, url, body)
	if nil != err {
		log.Fatalf("new request: %s", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return c.makeRequest(req)
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

func (c *ApiClient) makeRequest(req *http.Request) *ResponseWrapper {
	c.setupCookies(req)
	req.Header.Add("X-Csrf-Token", c.csrfToken)
	responseRecorder := httptest.NewRecorder()
	c.n.ServeHTTP(responseRecorder, req)
	c.saveCookies(responseRecorder)
	return newResponseWrapper(responseRecorder)
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

func (c *ApiClient) saveCookies(resp *httptest.ResponseRecorder) {
	chunks := strings.Split(resp.Header().Get("Set-Cookie"), "; ")
	for _, chunk := range chunks {
		if "" != chunk {
			c.cookies[chunk] = struct{}{}
		}
	}
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
