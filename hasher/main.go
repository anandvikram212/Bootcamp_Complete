package main

import (
	"github.com/gorilla/mux"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/sha1"
	"encoding/hex"
)

var Hash string
var Token string

func GetSHA256Hash(text string) string {
	hasher := sha1.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}
func GetHashValue(w http.ResponseWriter, r *http.Request) {

	tkn := Token

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}
	err = json.Unmarshal(jsn, &tkn)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	log.Printf("Received: %v\n", tkn)

	Hash = GetSHA256Hash(tkn)

	GetHashJson, err := json.Marshal(Hash)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(GetHashJson)

}

func Server() {
	http.HandleFunc("/hasher", GetHashValue)
	http.ListenAndServe(":8080", nil)

}

func Client(x string)  {

	//resp, err := http.Get("http://localhost:8088/stg/token?size=32")
	//if err != nil {
	// Serve an appropriately vague error to the
    // user, but log the details internally.
    //}
    TokenJson, err := json.Marshal(x)

	req, err := http.NewRequest("POST", "http://localhost:8080/hasher", bytes.NewBuffer(TokenJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("Response: ", string(body))
	resp.Body.Close()

}

func main() {
	//go Server()
	//Client("token")
	r := mux.NewRouter()
	r.HandleFunc("/hasher", GetHashValue).Methods("POST")
	if err := http.ListenAndServe(":8085",r); err != nil {
    	panic(err)
    }
}
