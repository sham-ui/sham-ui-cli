package main

import (
	"net/http"
	"regexp"
	"site/ssr"
	"site/test_helpers"
	"site/test_helpers/asserts"
	"testing"
)

func TestSSR(t *testing.T) {
	render := ssr.NewServerSideRender()
	env := test_helpers.NewTestEnv(render)
	revert := env.Default()
	render.Start()
	defer render.Stop()
	defer revert()

	rehydrateDataRegexp := regexp.MustCompile(`<script>.+<\/script>`)

	testCases := []struct {
		Name                       string
		URL                        string
		ExpectedResponseStatusCode int
		ExpectedResponse           string
	}{
		{
			Name:                       "contact",
			URL:                        "/contact",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponse:           "<!doctype html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"theme-color\" content=\"\"><link rel=\"icon\" href=\"/favicon.ico\"><title>{{ logoText }} | Contact</title><link rel=\"stylesheet\" href=\"/bundle.css\" /></head><body class=\"animation-stopped\"><nav class=\"navbar navbar-expand-lg fixed-top\"><div class=\"container-fluid\"><div class=\"logo\"><a href=\"\"><img src=\"/images/logo-dark.png\" alt=\"logo\"></img></a></div><!--0--><div class=\"collapse navbar-collapse\"><ul class=\"navbar-nav ml-auto mr-auto\"><li class=\"nav-item\"><a href=\"\" class=\"nav-link\">Home</a></li><li class=\"nav-item\"><a href=\"/contact/\" class=\"nav-link\">Contact</a></li></ul></div><div class=\"navbar-right ml-auto\"><div class=\"theme-switch-wrapper\"><label for=\"checkbox\" class=\"theme-switch\"><input type=\"checkbox\" id=\"checkbox\"><div class=\"slider round\"></div></label></div><!--1--><div class=\"search-icon\"><i class=\"icon-search\"></i></div><button type=\"button\" aria-expanded=\"false\" aria-label=\"Toggle navigation\" class=\"navbar-toggler\"><span class=\"navbar-toggler-icon\"></span></button></div></div></nav><!--0--><div><section class=\"section pt-55\"><div class=\"container\"><div class=\"row\"><div class=\"container-fluid\"><div class=\"row\"><div class=\"categorie-title\"><h3>Contact us</h3></div></div></div></div><div class=\"row\"><div class=\"col-lg-10 offset-lg-1 mb-20\"> For contact with us send email to <img src=\"/images/email.png\" alt=\"contact-email\"></img></div></div></div></section><!--0--></div><!--0--><!--0--><!--0--><!--1--><footer class=\"footer\"><div class=\"container-fluid\"><div class=\"row\"><div class=\"col-lg-12\"><div class=\"back\"><a href=\"#\" class=\"back-top\"><i class=\"icon-up\"></i></a></div></div></div></div></footer><!--2--><!--3--><!--0--><!--0--></body></html>",
		},
		{
			Name:                       "home",
			URL:                        "/",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponse:           "<!doctype html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"theme-color\" content=\"\"><link rel=\"icon\" href=\"/favicon.ico\"><title>{{ logoText }} | Home</title><link rel=\"stylesheet\" href=\"/bundle.css\" /></head><body class=\"animation-stopped\"><nav class=\"navbar navbar-expand-lg fixed-top\"><div class=\"container-fluid\"><div class=\"logo\"><a href=\"\"><img src=\"/images/logo-dark.png\" alt=\"logo\"></img></a></div><!--0--><div class=\"collapse navbar-collapse\"><ul class=\"navbar-nav ml-auto mr-auto\"><li class=\"nav-item\"><a href=\"\" class=\"nav-link active\">Home</a></li><li class=\"nav-item\"><a href=\"/contact/\" class=\"nav-link\">Contact</a></li></ul></div><div class=\"navbar-right ml-auto\"><div class=\"theme-switch-wrapper\"><label for=\"checkbox\" class=\"theme-switch\"><input type=\"checkbox\" id=\"checkbox\"><div class=\"slider round\"></div></label></div><!--1--><div class=\"search-icon\"><i class=\"icon-search\"></i></div><button type=\"button\" aria-expanded=\"false\" aria-label=\"Toggle navigation\" class=\"navbar-toggler\"><span class=\"navbar-toggler-icon\"></span></button></div></div></nav><!--0--><div><section class=\"section pt-55\"><!--0--><ul class=\"errors-panel\"><li>Fail load articles list</li></ul><!--0--><!--1--><div class=\"container-fluid\"><div class=\"row\"><!--2--><div class=\"col-lg-12\"><div class=\"pagination mt--10\"><ul class=\"list-inline\"><!--0--><li><a href=\"/page/1/\" class=\"active\">1</a></li><!--1--></ul></div></div></div></div></section><!--0--></div><!--0--><!--0--><!--0--><!--0--><!--1--><footer class=\"footer\"><div class=\"container-fluid\"><div class=\"row\"><div class=\"col-lg-12\"><div class=\"back\"><a href=\"#\" class=\"back-top\"><i class=\"icon-up\"></i></a></div></div></div></div></footer><!--2--><!--3--><!--0--><!--0--></body></html>",
		},
	}

	for _, testCase := range testCases {
		resp := env.API.Request(http.MethodGet, testCase.URL, nil)
		// Remove rehydrator data, because it's changes after each build
		html := rehydrateDataRegexp.ReplaceAllString(resp.Text(), "")
		asserts.Equals(t, testCase.ExpectedResponseStatusCode, resp.Response.Code, testCase.Name+" response code")
		asserts.Equals(t, testCase.ExpectedResponse, html, testCase.Name+" response body")
	}
}
