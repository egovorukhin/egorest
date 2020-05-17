package egorest

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

type Body struct {
	FormatBody FormatBody
	Data interface{}
}

//Marshal JSON XML
func (body Body) marshal() ([]byte, error) {
	switch body.FormatBody {
	case JSON: return body.marshalJson()
	case XML: return body.marshalXml()
	}
	return nil, errors.New("Неизвестный формат данных")
}

func (body Body) marshalJson() ([]byte, error) {
	b, err := json.Marshal(&body.Data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (body Body) marshalXml() ([]byte, error) {
	b, err := xml.Marshal(&body.Data)
	if err != nil {
		return nil, err
	}
	return b, nil
}
