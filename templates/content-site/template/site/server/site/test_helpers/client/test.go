package client

type TestCase struct {
	Message                    string
	Method                     string
	URL                        string
	Data                       map[string]interface{}
	ExpectedResponseStatusCode int
	ExpectedResponseJSON       map[string]interface{}
}
