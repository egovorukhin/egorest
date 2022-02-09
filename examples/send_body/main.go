package main

import (
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
	}
	user := User{
		Id:   1,
		Name: "User",
	}
	// Custom format
	r := egorest.NewRequest("/api/user").SetBody(egorest.MIMEApplicationJSON, user)
	resp, err := egorest.NewClient(cfg).Post(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)

	// JSON
	r = egorest.NewRequest("/api/user").JSON(user)
	resp, err = egorest.NewClient(cfg).Post(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)

	// XML
	r = egorest.NewRequest("/api/user").XML(user)
	resp, err = egorest.NewClient(cfg).Post(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)
}
