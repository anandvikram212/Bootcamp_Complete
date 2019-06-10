package main

import (
	"github.com/gorilla/mux"
	"bytes"
	"encoding/json"
	"fmt"
	//"strconv"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/rand"
	"encoding/base64"
)

var bits string
var Token string
var SecureToken string

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func SecureTokenGenerator(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("get params ", r.URL.Query())

	//randomBits := bits

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	var x int

	err = json.Unmarshal(jsn, &x)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	log.Printf("Received: %v\n", x)

 	//randomBits = r.URL.Query().Get("size")



	b, err := GenerateRandomBytes(x)
	tkn, err := base64.URLEncoding.EncodeToString(b), err

	GenTokenJson, err := json.Marshal(tkn)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(GenTokenJson)

}

func Server() {
	http.HandleFunc("/stg/token", SecureTokenGenerator)
	http.ListenAndServe(":8081", nil)

}

func Client() {

	//resp, err := http.Get("http://localhost:8088/stg/token?size=32")
	//if err != nil {
	// Serve an appropriately vague error to the
    // user, but log the details internally.
    //}
    SizeJson, err := json.Marshal(32)

	req, err := http.NewRequest("GET", "http://localhost:8081/stg/token?size=32", bytes.NewBuffer(SizeJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("Response: ", string(body))
	SecureToken = string(body)
	resp.Body.Close()

}

func main() {
	//go Server()
	//Client()
	r := mux.NewRouter()
	r.HandleFunc("/stg/token", SecureTokenGenerator).Methods("GET")
	if err := http.ListenAndServe(":8081", r); err != nil {
    	panic(err)
    }
}
