package egorest

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

type FormatBody string

func (f FormatBody) String() string {
	return string(f)
}

const (
	NONE FormatBody = "none"
	JSON FormatBody = "application/json"
	XML  FormatBody = "application/xml"
)

//Unmarshal JSON XML
func (f FormatBody) unmarshal(data []byte, v interface{}) error {
	switch f {
	case JSON:
		return f.unmarshalJson(data, v)
	case XML:
		return f.unmarshalXml(data, v)
	default:
		return errors.New("Неизвестный формат данных")
	}
}

func (f FormatBody) unmarshalJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func (f FormatBody) unmarshalXml(data []byte, v interface{}) error {
	err := xml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func getFormatBody(s string) FormatBody {

	//JSON
	if strings.Contains(s, JSON.String()) {
		return JSON
	}

	//XML
	if strings.Contains(s, XML.String()) {
		return XML
	}
	return NONE
}
