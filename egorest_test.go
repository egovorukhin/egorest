package egorest

import (
	"testing"
	"time"
)

type Provider struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func TestClient_Execute(t *testing.T) {

	//http://ip-api.com/json/
	cfg := Config{
		BaseUrl: BaseUrl{
			Url:    "",
			Schema: "http",
			Host:   "ip-api.com",
			Port:   80,
			Path:   "json",
		},
		Secure:  false,
		Timeout: time.Second * 30,
		Buffers: &Buffers{
			Write: 4096,
			Read:  4096,
		},
		Proxy:     nil,
		BasicAuth: nil,
	}

	// Структура ответа
	p := Provider{}
	// Новый запрос
	r := NewRequest("172.19.151.219")
	// Инициализация клиента
	err := NewClient(cfg).ExecuteGet(r, &p, nil)
	if err != nil {
		t.Fatal(err)
	}

}
