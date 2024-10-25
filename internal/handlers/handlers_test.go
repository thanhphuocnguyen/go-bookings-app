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

var testCases = []struct {
	name             string
	url              string
	method           string
	params           []postData
	expectStatusCode int
}{
	{"home", "/", "GET", []postData{}, 200},
	{"about", "/about", "GET", []postData{}, 200},
	{"contact", "/contact", "GET", []postData{}, 200},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, 200},
	{"majors-suite", "/majors-suite", "GET", []postData{}, 200},
	{"make-reservation", "/make-reservation", "GET", []postData{}, 200},
	{"post-make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Thanh Phuoc"},
		{key: "last_name", value: "Nguyen"},
		{key: "email", value: "testing@example.com"},
		{key: "phone", value: "123456789"},
	}, 200},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, 200},
	{"search-availability", "/search-availability", "GET", []postData{}, 200},
	{"post-search-availability", "/search-availability", "POST", []postData{}, 200},
	{"search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2021-01-01"},
		{key: "end", value: "2021-01-02"},
	}, 200},
}

func TestHandlers(t *testing.T) {
	// arrange
	routes := getRoutes()
	// httptest.NewTLSServer is used to create a new test server with TLS
	ts := httptest.NewTLSServer(routes)

	// defer is used to close the server after the test is done
	defer ts.Close()

	for _, tt := range testCases {
		// t.Run is used to run subtests and sub-benchmarks in a single test
		t.Run(tt.name, func(t *testing.T) {
			// act
			var res *http.Response
			var err error
			// ts.Client() is used to create a new client to make requests to the server
			client := ts.Client()
			if tt.method == "GET" {
				res, err = client.Get(ts.URL + tt.url)
			} else if tt.method == "POST" {
				values := url.Values{}
				for _, v := range tt.params {
					values.Add(v.key, v.value)
				}
				res, err = client.PostForm(ts.URL+tt.url, values)
			} else {
				t.Fatal("Unsupported method")
			}
			if err != nil {
				t.Fatal(err)
			}
			// assert
			if res.StatusCode != tt.expectStatusCode {
				t.Errorf("for %s, expected %d but got %d", tt.name, tt.expectStatusCode, res.StatusCode)
			}
		})
	}

}
