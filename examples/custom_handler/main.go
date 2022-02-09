package main

import (
	"encoding/json"
	"fmt"
	"github.com/egovorukhin/egorest"
	"log"
	"time"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	cfg := egorest.Config{
		BaseUrl: egorest.BaseUrl{
			Schema: "https",
			Host:   "localhost",
			Port:   5005,
			Path:   "api",
		},
		Secure:  true,
		Timeout: time.Second * 30,
		Buffers: &egorest.Buffers{
			Write: 4096,
			Read:  4096,
		},
		Proxy:     nil,
		BasicAuth: nil,
	}

	r := egorest.NewRequest("api/user")

	client := egorest.NewClient(cfg)
	user := User{}
	err := client.Execute(r, &user, func(contentType string, data []byte, v interface{}) error {
		if contentType == "application/json" {
			return json.Unmarshal(data, v)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", user)
}
