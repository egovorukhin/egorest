package egorest

import (
	"fmt"
	"testing"
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

type Response struct {
	Code    int    `json:"code"`
	Message []City `json:"message"`
}

type City struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
}

func TestClient_Send(t *testing.T) {

	responseBody := Response{}
	err := NewClient("dls.hq.bc", 80, false).
		SetTimeout(15).
		SetRequest(NewRequest(GET)).
		Execute("/api/place/city", &responseBody)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(responseBody)
}
