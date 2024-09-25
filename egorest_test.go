package egorest

import (
	"fmt"
	"io"
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

func TestSetFormData(t *testing.T) {

	cfg := Config{
		BaseUrl: BaseUrl{
			Url: "http://localhost:7474",
		},
		Secure:  true,
		Timeout: time.Second * 30,
	}
	incidentId := "2854711"
	values := map[string]interface{}{
		"incidentId": incidentId,
		"userLogin":  "govorukhin_35893",
		"files":      []string{"C:\\downloads\\[new-bucket-1b5f4695]test.txt"},
	}
	r := NewRequest("/api/file")
	err := r.FormData(values)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := NewClient(cfg).Post(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))
}
