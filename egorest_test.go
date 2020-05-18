package egorest

import (
	"fmt"
	"io/ioutil"
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
	req := NewRequest(GET, "api/place/city").
		SetHeader(SetHeader("Connection", "Keep-Alive"))

	err := NewClient("dls.hq.bc", 80, false).
		SetTimeout(15).
		Execute(req, &responseBody)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Struct: %v\n", responseBody)

	resp, err := NewClient("dls.hq.bc", 80, false).
		SetTimeout(15).
		Send(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Response.Body: %s\n", body)
}
