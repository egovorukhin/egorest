package main

import (
	"fmt"
	"github.com/egovorukhin/egorest"
	"log"
	"time"
)

func main() {

	cfg := egorest.Config{
		BaseUrl: egorest.BaseUrl{
			Url: "https://localhost:5005",
		},
		Secure:  true,
		Timeout: time.Second * 30,
	}

	r := egorest.NewRequest("api/rec/audio")
	err := r.AddFiles("audio", "d:\\audio\\02022022_1136_700043.wav", "d:\\audio\\21012022_0730_user2.wav")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := egorest.NewClient(cfg).Post(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
}
