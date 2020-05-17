package egorest

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

type Data struct {
	FormatBody FormatBody
	Body       interface{}
}

//Marshal JSON XML
func (data Data) marshal() ([]byte, error) {
	switch data.FormatBody {
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
