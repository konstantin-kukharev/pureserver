package main

import (
	"encoding/json"
	"fmt"
	ps "github.com/konstantin-kukharev/pureserver"
	"strconv"
)

type TestResponse struct {
	A int `json:"A"`
	B int `json:"B"`
	C int `json:"C"`
}

func HelloServer(w ps.ResponseWriter, req ps.HttpRequestInterface) {
	var body TestResponse
	_ = json.Unmarshal([]byte(req.GetBody()), &body)
	name, _ := req.GetParam("increment")
	app, _ := strconv.Atoi(name)
	body.A += app
	body.B *= app
	body.C = body.C / app
	result, _ := json.Marshal(body)
	w.SetBody(result)
}

func main() {
	mux := ps.NewMux()
	mux.Post("/hello/:increment", ps.HandlerFunc(HelloServer))

	server := ps.NewHttp(mux)
	server.SetPort(8080)
	fmt.Println(server.Serve())
}
