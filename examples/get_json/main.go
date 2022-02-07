package main

import (
	"fmt"
	"github.com/egovorukhin/egorest"
	"log"
	"time"
)

type App struct {
	Hostname       string     `json:"hostname"`
	Login          string     `json:"login"`
	Version        string     `json:"version"`
	Started        bool       `json:"started"`
	Domain         string     `json:"domain"`
	Keys           []string   `json:"keys"`
	ConnectTime    time.Time  `json:"connect_time"`
	DisconnectTime *time.Time `json:"disconnect_time,omitempty"`
	IpAddress      string     `json:"ip_address"`
}

type Apps []App

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

	// Структура ответа
	apps := Apps{}
	// Новый запрос
	r := egorest.NewRequest("chrome/apps")
	// Инициализация клиента
	err := egorest.NewClient(cfg).ExecuteGet(r, &apps)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v", apps)
}
