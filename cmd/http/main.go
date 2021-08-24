package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type TestResponse struct {
	A int `json:"A"`
	B int `json:"B"`
	C int `json:"C"`
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/hello/{increment:[0-9]+}", HelloHandler).Methods("POST")

	http.Handle("/", rtr)
	fmt.Println("Server started at port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	var body TestResponse
	_ = json.NewDecoder(r.Body).Decode(&body)
	params := mux.Vars(r)
	name := params["name"]
	app, _ := strconv.Atoi(name)
	body.A += app
	body.B *= app
	body.C = body.C / app
	result, _ := json.Marshal(body)
	w.WriteHeader(200)
	w.Write(result)
}
