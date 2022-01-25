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

	var city []City
	//Execute
	req := NewRequest(GET, "place/city?name=Алматы").
		SetHeader(SetHeader("Connection", "keep-alive"))

	var user string
	var password string

	err := NewClient("localhost", 443, true).
		SetBasicAuth(user, password).
		SetTimeout(15).
		SetRoute("api").
		Execute(req, &city)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Struct: %v\n", city)

	//Send Filter
	req.Data = &Data{
		ContentType: JSON,
		Body: City{
			Name: "Алматы",
		},
	}
	resp, err := NewClient("localhost", 443, true).
		SetTimeout(15).
		SetRoute("api/").
		// указываем контекст
		//WithContext(context.Background()).
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

	//By Uri
	client, err := NewClientByUri("https://localhost/api/place/city")
	if err != nil {
		t.Error(err)
	}
	client.SetTimeout(15)
	err = client.Execute(NewRequest(GET, ""), &city)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Struct: %v\n", city)

}
