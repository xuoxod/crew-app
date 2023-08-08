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
	{"dashboard", "/dashboard", "GET", []postData{}, http.StatusOK},
	{"register", "/register", "GET", []postData{}, http.StatusOK},
	{"signin", "/signin", "POST", []postData{
		{key: "fname", value: "anita"},
		{key: "lname", value: "bathe"},
		{key: "email", value: "ani@google.net"},
		{key: "phone", value: "213-777-9311"},
		{key: "pwd1", value: "123"},
		{key: "pwd2", value: "123"},
	}, http.StatusOK},
	{"register", "/register", "POST", []postData{
		{key: "email", value: "ani@google.net"},
		{key: "password", value: "123"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			response, err := ts.Client().Get(ts.URL + e.url)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("\n\tfor %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			response, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("\n\tfor %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
