package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	ps "github.com/konstantin-kukharev/pureserver"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

var (
	Client HTTPClient
)

type TestResponse struct {
	A int `json:"A"`
	B int `json:"B"`
	C int `json:"C"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func HelloServer(w ps.ResponseWriter, req ps.HttpRequestInterface) {
	var body TestResponse
	_ = json.Unmarshal([]byte(req.GetBody()), &body)
	name, _ := req.GetParam("name")
	app, _ := strconv.Atoi(name)
	body.A += app
	body.B *= app
	body.C = body.C / app
	result, _ := json.Marshal(body)
	w.SetBody(result)
}

func HelloServer2(w ps.ResponseWriter, req ps.HttpRequestInterface) {
	var body TestResponse
	_ = json.Unmarshal([]byte(req.GetBody()), &body)
	name, _ := req.GetParam("name")
	app, _ := strconv.Atoi(name)
	body.A += app
	body.B *= app
	body.C = body.C / app
	result, _ := json.Marshal(body)
	w.SetBody(result)
}

func TestHttpServerRoutes(t *testing.T) {
	go serverUp(8080)

	incr := 10
	tr := TestResponse{1, 1, 1}
	test, _ := json.Marshal(tr)
	s, b, err := makeRequest(8080, `hello/`+strconv.Itoa(incr), test)
	var br TestResponse
	_ = json.Unmarshal(b, &br)
	if br.A != tr.A+incr || br.B != tr.B*incr || br.C != tr.C/incr {
		panic(
			fmt.Sprintf("%s, %s", s, err),
		)
	}
}

func serverUp(ports ...int) {
	mux := ps.NewMux()
	mux.Get("/hello/:name", ps.HandlerFunc(HelloServer))
	mux.Post("/hello/:name", ps.HandlerFunc(HelloServer))

	server := ps.NewHttp(mux)
	server.SetPort(ports...)

	fmt.Println(server.Serve())
}

func makeRequest(port int, route string, test []byte) (status string, body []byte, err error) {
	resp, err := Get(fmt.Sprintf("http://0.0.0.0:%d/%s", port, route), test, http.Header{})
	if err != nil {
		panic(
			fmt.Sprintf("%v", err),
		)
	}

	b, err := ioutil.ReadAll(resp.Body)

	return resp.Status, b, err
}

func Post(url string, body []byte, headers http.Header) (*http.Response, error) {
	Client = &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}

func Get(url string, body []byte, headers http.Header) (*http.Response, error) {
	Client = &http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}
