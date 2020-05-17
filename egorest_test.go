package egorest

import (
	"fmt"
	"testing"
)

type ResponseBody struct {
	Status string `json:"status"`
	Country string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region string `json:"region"`
	RegionName string `json:"regionName"`
	City string `json:"city"`
	Zip string `json:"zip"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Timezone string `json:"timezone"`
	Isp string `json:"isp"`
	Org string `json:"org"`
	As string `json:"as"`
	Query string `json:"query"`
}

func TestClient_Send(t *testing.T) {

	responseBody := ResponseBody{}
	err := NewClient("ip-api.com", 80, false).
		SetTimeout(15).
		SetRequest(NewRequest(GET).AddHeader("Accept", "*/*")).
		Send("/json/145.255.163.43", &responseBody)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(responseBody)
}
