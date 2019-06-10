package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	//"crypto/rand"
	//"encoding/base64"
	"fmt"
	"time"
	"strconv"
	//"crypto/sha1"
	//"encoding/hex"
	"github.com/go-redis/redis"
)

var client *redis.Client

func main() {

	SizeJson, err := json.Marshal(32)

	req, err := http.NewRequest("GET", "http://localhost:8081/stg/token", bytes.NewBuffer(SizeJson))
	req.Header.Set("Content-Type", "application/json")

	client1 := &http.Client{}
	resp, err := client1.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("Response: ", string(body))
	SecureToken := string(body)
	resp.Body.Close()

	xtoken := SecureToken

	fmt.Println("generated Token is ", xtoken)

	TokenJson, err := json.Marshal(xtoken)

	req, err = http.NewRequest("POST", "http://localhost:8085/hasher", bytes.NewBuffer(TokenJson))
	req.Header.Set("Content-Type", "application/json")

	client1 = &http.Client{}
	resp, err = client1.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	//fmt.Println("Response: ", string(body))
	hash := string(body)
	resp.Body.Close()

	fmt.Println("Hash value is ", hash)

	//if (hash[0]=='0') {

		client = redis.NewClient(&redis.Options{
	   	Addr:     "localhost:6379",
		})

		if(client.LLen("dates").Val()==0) {
			client.RPush("dates",time.Now().Format("02-01-2006"))
			client.RPush("counts",0)
		}
		dates := client.RPop("dates").Val()
	    client.RPush("dates",dates)

		if(dates==time.Now().Format("02-01-2006")){
			fmt.Println("hi there")
			dlen := client.LLen("dates").Val()
			clen := client.LLen("counts").Val()
			if(dlen==clen){
				fmt.Println("hi fsffvfvdere")
				cur_count := client.RPop("counts").Val()
				i,err := strconv.Atoi(cur_count)
				if(err!=nil) {
					panic(err)
				}
				i++
				client.RPush("counts",i)
			}else {
				fmt.Println("dvfdere")
				client.RPush("counts",1)
			}

		}else {
			client.RPush("dates",time.Now().Format("02-01-2006"))
			client.RPush("counts",1)
		}

		val, err := client.LRange("LuckyHashes",0,-1).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("hash", val)

		defer client.Close()

		err = client.Publish("mychannel1", "hello").Err()
		if err != nil {
			panic(err)
		}
	//} else {
	//	fmt.Println("Not a LuckyHash value")
	//}

}
