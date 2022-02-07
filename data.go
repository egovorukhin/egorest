package egorest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
)

type Data struct {
	ContentType string
	Body        interface{}
}

type ContentType string

type UnmarshalHandler func([]byte, interface{}) error
type MarshalHandler func(interface{}) ([]byte, error)

// MIME types that are commonly used
const (
	MIMETextXML               = "text/xml"
	MIMETextHTML              = "text/html"
	MIMETextPlain             = "text/plain"
	MIMEApplicationXML        = "application/xml"
	MIMEApplicationJSON       = "application/json"
	MIMEApplicationJavaScript = "application/javascript"
	MIMEApplicationForm       = "application/x-www-form-urlencoded"
	MIMEOctetStream           = "application/octet-stream"
	MIMEMultipartForm         = "multipart/form-data"

	MIMETextXMLCharsetUTF8               = "text/xml; charset=utf-8"
	MIMETextHTMLCharsetUTF8              = "text/html; charset=utf-8"
	MIMETextPlainCharsetUTF8             = "text/plain; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8        = "application/xml; charset=utf-8"
	MIMEApplicationJSONCharsetUTF8       = "application/json; charset=utf-8"
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

// Unmarshal Content-Type
func (c ContentType) unmarshal(data []byte, v interface{}, handler ...UnmarshalHandler) error {
	switch string(c) {
	case MIMEApplicationJSON,
		MIMEApplicationJSONCharsetUTF8:
		return c.json(data, v)
	case MIMEApplicationXML,
		MIMEApplicationXMLCharsetUTF8:
		return c.xml(data, v)
	case MIMETextXML,
		MIMETextHTML,
		MIMETextPlain,
		MIMETextXMLCharsetUTF8,
		MIMETextHTMLCharsetUTF8,
		MIMETextPlainCharsetUTF8:
		return c.text(data, v)
	}
	if len(handler) > 0 {
		return handler[0](data, v)
	}
	return c.none(data, v)
}

// Unmarshal JSON
func (c ContentType) json(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal XML
func (c ContentType) xml(data []byte, v interface{}) error {
	err := xml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal TEXT or Unknown -> Text
func (c ContentType) text(data []byte, v interface{}) error {
	v = string(data)
	return nil
}

// Unmarshal none
func (c ContentType) none(data []byte, v interface{}) error {
	return c.text(data, v)
}

// Marshal данных
func (data Data) marshal(handler ...MarshalHandler) (io.Reader, error) {
	var body []byte
	err := errors.New("неизвестный формат данных")
	switch ContentType(data.ContentType) {
	case MIMEApplicationJSON,
		MIMEApplicationJSONCharsetUTF8:
		body, err = data.json()
		break
	case MIMEApplicationXML,
		MIMEApplicationXMLCharsetUTF8:
		body, err = data.xml()
		break
	default:
		if len(handler) > 0 {
			body, err = handler[0](data.Body)
			break
		}
		if buf, ok := data.Body.(*bytes.Buffer); ok {
			return buf, nil
		}
	}
	return bytes.NewBuffer(body), err
}

func (data Data) json() ([]byte, error) {
	b, err := json.Marshal(&data.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (data Data) xml() ([]byte, error) {
	b, err := xml.Marshal(&data.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
