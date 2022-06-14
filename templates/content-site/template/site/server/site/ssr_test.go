package main

import (
	"net/http"
	"site/ssr"
	"site/test_helpers"
	"site/test_helpers/asserts"
	"testing"
)

func TestSSR(t *testing.T) {
	render := ssr.NewServerSideRender()
	env := test_helpers.NewTestEnv(render)
	revert := env.Default()
	defer render.Stop()
	defer revert()

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
			ExpectedResponse:           "<!doctype html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"theme-color\" content=\"\"><link rel=\"icon\" href=\"/favicon.ico\"><title>{{ logoText }} | Contact</title><link rel=\"stylesheet\" href=\"/bundle.css\" /></head><body class=\"animation-stopped\"><nav class=\"navbar navbar-expand-lg fixed-top\"><div class=\"container-fluid\"><div class=\"logo\"><a href=\"\"><img src=\"/images/logo-dark.png\" alt=\"logo\"></img></a></div><!--0--><div class=\"collapse navbar-collapse\"><ul class=\"navbar-nav ml-auto mr-auto\"><li class=\"nav-item\"><a href=\"\" class=\"nav-link\">Home</a></li><li class=\"nav-item\"><a href=\"/contact/\" class=\"nav-link\">Contact</a></li></ul></div><div class=\"navbar-right ml-auto\"><div class=\"theme-switch-wrapper\"><label for=\"checkbox\" class=\"theme-switch\"><input type=\"checkbox\" id=\"checkbox\"><div class=\"slider round\"></div></label></div><!--1--><div class=\"search-icon\"><i class=\"icon-search\"></i></div><button type=\"button\" aria-expanded=\"false\" aria-label=\"Toggle navigation\" class=\"navbar-toggler\"><span class=\"navbar-toggler-icon\"></span></button></div></div></nav><!--0--><div><section class=\"section pt-55\"><div class=\"container\"><div class=\"row\"><div class=\"container-fluid\"><div class=\"row\"><div class=\"categorie-title\"><h3>Contact us</h3></div></div></div></div><div class=\"row\"><div class=\"col-lg-10 offset-lg-1 mb-20\"> For contact with us send email to <img src=\"/images/email.png\" alt=\"contact-email\"></img></div></div></div></section><!--0--></div><!--0--><!--0--><!--0--><!--1--><footer class=\"footer\"><div class=\"container-fluid\"><div class=\"row\"><div class=\"col-lg-12\"><div class=\"back\"><a href=\"#\" class=\"back-top\"><i class=\"icon-up\"></i></a></div></div></div></div></footer><!--2--><!--3--><!--0--><!--0--><script>window.data={\"0\":{\"C\":[\"1\"],\"N\":[10]},\"1\":{\"C\":[\"2\",\"6\",\"7\"],\"N\":[1,6,8,9]},\"2\":{\"C\":[\"3\",\"4\",\"5\"],\"N\":[0,[0,0],[0,0,1],[0,0,2],[0,0,3],[0,0,3,1],[0,0,3,2],[0,0,3,2,0],[0,0,3,3],[0,0,3,3,0]]},\"3\":{\"N\":[0,[0,0],[0,0,0]]},\"4\":{\"N\":[0,[0,0],[0,0,0],[0,0,0,0],[0,1],[0,1,0],[0,1,0,0]]},\"5\":{\"N\":[0,[0,0],[0,0,0],[0,0,1]]},\"6\":{\"N\":[7,[7,0],[7,0,0],[7,0,0,0],[7,0,0,0,0],[7,0,0,0,0,0],[7,0,0,0,0,0,0]]},\"7\":{\"C\":[\"8\"],\"N\":[5]},\"8\":{\"C\":[\"9\"]},\"9\":{\"C\":[\"10\"],\"N\":[4]},\"10\":{\"C\":[\"11\"],\"N\":[3]},\"11\":{\"C\":[\"12\"],\"N\":[2,[2,1]]},\"12\":{\"N\":[0,[0,0],[0,0,0],[0,0,0,0],[0,0,0,0,0],[0,0,0,0,0,0],[0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],[0,0,1],[0,0,1,0],[0,0,1,0,0],[0,0,1,0,1]]},\"app\":{\"C\":[\"0\"],\"N\":[11]}};</script><script src=\"/s.min.js\"></script><script>System.import('/bundle.js');</script></body></html>",
		},
		{
			Name:                       "home",
			URL:                        "/",
			ExpectedResponseStatusCode: http.StatusOK,
			ExpectedResponse:           "<!doctype html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"theme-color\" content=\"\"><link rel=\"icon\" href=\"/favicon.ico\"><title>{{ logoText }} | Home</title><link rel=\"stylesheet\" href=\"/bundle.css\" /></head><body class=\"animation-stopped\"><nav class=\"navbar navbar-expand-lg fixed-top\"><div class=\"container-fluid\"><div class=\"logo\"><a href=\"\"><img src=\"/images/logo-dark.png\" alt=\"logo\"></img></a></div><!--0--><div class=\"collapse navbar-collapse\"><ul class=\"navbar-nav ml-auto mr-auto\"><li class=\"nav-item\"><a href=\"\" class=\"nav-link active\">Home</a></li><li class=\"nav-item\"><a href=\"/contact/\" class=\"nav-link\">Contact</a></li></ul></div><div class=\"navbar-right ml-auto\"><div class=\"theme-switch-wrapper\"><label for=\"checkbox\" class=\"theme-switch\"><input type=\"checkbox\" id=\"checkbox\"><div class=\"slider round\"></div></label></div><!--1--><div class=\"search-icon\"><i class=\"icon-search\"></i></div><button type=\"button\" aria-expanded=\"false\" aria-label=\"Toggle navigation\" class=\"navbar-toggler\"><span class=\"navbar-toggler-icon\"></span></button></div></div></nav><!--0--><div><section class=\"section pt-55\"><!--0--><ul class=\"errors-panel\"><li>Fail load articles list</li></ul><!--0--><!--1--><div class=\"container-fluid\"><div class=\"row\"><!--2--><div class=\"col-lg-12\"><div class=\"pagination mt--10\"><ul class=\"list-inline\"><!--0--><li><a href=\"/page/1/\" class=\"active\">1</a></li><!--1--></ul></div></div></div></div></section><!--0--></div><!--0--><!--0--><!--0--><!--0--><!--1--><footer class=\"footer\"><div class=\"container-fluid\"><div class=\"row\"><div class=\"col-lg-12\"><div class=\"back\"><a href=\"#\" class=\"back-top\"><i class=\"icon-up\"></i></a></div></div></div></div></footer><!--2--><!--3--><!--0--><!--0--><script>window.data={\"0\":{\"C\":[\"1\"],\"N\":[11]},\"1\":{\"C\":[\"2\",\"6\",\"7\"],\"N\":[1,7,9,10]},\"2\":{\"C\":[\"3\",\"4\",\"5\"],\"N\":[0,[0,0],[0,0,1],[0,0,2],[0,0,3],[0,0,3,1],[0,0,3,2],[0,0,3,2,0],[0,0,3,3],[0,0,3,3,0]]},\"3\":{\"N\":[0,[0,0],[0,0,0]]},\"4\":{\"N\":[0,[0,0],[0,0,0],[0,0,0,0],[0,1],[0,1,0],[0,1,0,0]]},\"5\":{\"N\":[0,[0,0],[0,0,0],[0,0,1]]},\"6\":{\"N\":[8,[8,0],[8,0,0],[8,0,0,0],[8,0,0,0,0],[8,0,0,0,0,0],[8,0,0,0,0,0,0]]},\"7\":{\"C\":[\"8\"],\"N\":[6]},\"8\":{\"C\":[\"9\"]},\"9\":{\"C\":[\"10\"],\"N\":[5]},\"10\":{\"O\":{\"10\":\"home.page\",\"13\":true,\"15\":[\"Fail load articles list\"]},\"C\":[\"11\"],\"N\":[4]},\"11\":{\"C\":[\"14\"],\"N\":[3]},\"14\":{\"C\":[\"15\"],\"N\":[2,[2,1]]},\"15\":{\"C\":[\"16\",\"17\"],\"N\":[0,[0,0],[0,3],[0,4],[0,4,0],[0,4,0,0],[0,4,0,1]]},\"16\":{\"C\":[\"19\"],\"N\":[2]},\"17\":{\"C\":[\"18\"],\"N\":[0,[0,0]]},\"18\":{\"N\":[0,1,[1,0],[1,0,0],2]},\"19\":{\"O\":{\"15\":[\"Fail load articles list\"]},\"C\":[\"20\"],\"N\":[1]},\"20\":{\"N\":[0,[0,0]]},\"app\":{\"C\":[\"0\"],\"N\":[12]}};</script><script src=\"/s.min.js\"></script><script>System.import('/bundle.js');</script></body></html>",
		},
	}

	for _, testCase := range testCases {
		resp := env.API.Request(http.MethodGet, testCase.URL, nil)
		asserts.Equals(t, testCase.ExpectedResponseStatusCode, resp.Response.Code, testCase.Name+" response code")
		asserts.Equals(t, testCase.ExpectedResponse, resp.Text(), testCase.Name+" response body")
	}
}
