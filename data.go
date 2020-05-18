package egorest

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

type Data struct {
	ContentType ContentType
	Body        interface{}
}

type ContentType string

func (f ContentType) String() string {
	return string(f)
}

const (
	NONE     ContentType = "none"
	JSON     ContentType = "application/json"
	XML      ContentType = "application/xml"
	TEXT_XML ContentType = "text/xml"
)

//Unmarshal JSON XML
func (f ContentType) unmarshal(data []byte, v interface{}) error {
	switch f {
	case JSON:
		return f.unmarshalJson(data, v)
	case XML:
		return f.unmarshalXml(data, v)
	default:
		return errors.New("Неизвестный формат данных")
	}
}

func (f ContentType) unmarshalJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func (f ContentType) unmarshalXml(data []byte, v interface{}) error {
	err := xml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func getFormatBody(s string) ContentType {

	//JSON
	if strings.Contains(s, JSON.String()) {
		return JSON
	}

	//XML
	if strings.Contains(s, XML.String()) || strings.Contains(s, TEXT_XML.String()) {
		return XML
	}

	return NONE
}

//Marshal JSON XML
func (data Data) marshal() ([]byte, error) {
	switch data.ContentType {
	case JSON:
		return data.marshalJson()
	case XML:
		return data.marshalXml()
	default:
		return nil, errors.New("Неизвестный формат данных")
	}
}

func (data Data) marshalJson() ([]byte, error) {
	b, err := json.Marshal(&data.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (data Data) marshalXml() ([]byte, error) {
	b, err := xml.Marshal(&data.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
