package egorest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

type Data struct {
	ContentType string
	Body        interface{}
}

type ContentType string

func (f ContentType) String() string {
	return string(f)
}

const (
	NONE    ContentType = "none"
	JSON    ContentType = "application/json"
	XML     ContentType = "application/xml"
	TextXml ContentType = "text/xml"
)

// Unmarshal JSON XML
func (f ContentType) unmarshal(data []byte, v interface{}) error {
	switch f {
	case JSON:
		return f.unmarshalJson(data, v)
	case XML:
		return f.unmarshalXml(data, v)
	default:
		return f.unmarshalNone(data, v)
	}
}

// Unmarshal JSON
func (f ContentType) unmarshalJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal XML
func (f ContentType) unmarshalXml(data []byte, v interface{}) error {
	err := xml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal TEXT or Unknown -> Text
func (f ContentType) unmarshalNone(data []byte, v interface{}) error {
	v = string(data)
	return nil
}

func getFormatBody(s string) ContentType {

	//JSON
	if strings.Contains(s, JSON.String()) {
		return JSON
	}
	//XML
	if strings.Contains(s, XML.String()) || strings.Contains(s, TextXml.String()) {
		return XML
	}

	return NONE
}

// Marshal JSON XML
func (data Data) marshal() (bytes.Buffer, error) {
	var body []byte
	err := errors.New("неизвестный формат данных")
	switch ContentType(data.ContentType) {
	case JSON:
		body, err = data.marshalJson()
		break
	case XML:
		body, err = data.marshalXml()
		break
	default:
		if buf, ok := data.Body.(bytes.Buffer); ok {
			return buf, nil
		}
	}
	return *bytes.NewBuffer(body), err
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
