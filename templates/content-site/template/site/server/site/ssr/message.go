package ssr

type nodejsRequest struct {
	ID      string `json:"id"`
	Origin  string `json:"origin"`
	URL     string `json:"url"`
	API     string `json:"api"`
	Cookies string `json:"cookies"`
}

type nodejsResponse struct {
	ID    string `json:"id"`
	HTML  string `json:"html"`
	Error string `json:"error"`
}
