package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post availability json", "/search-availability-json", "GET", []postData{}, http.StatusOK},
	{"post reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		switch test.method {
		case "GET":
			{
				resp, err := testServer.Client().Get(testServer.URL + test.url)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}

				if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("for %s expected %d, but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
				}
			}
		case "POST":
			{
				values := url.Values{}
				for _, param := range test.params {
					values.Add(param.key, param.value)
				}

				resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}

				if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("for %s expected %d, but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
				}
			}
		}
	}
}
