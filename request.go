package egorest

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Request struct {
	Method  string
	Headers map[string]string
	Route   string
	Data    *Data
}

type Header struct {
	Name  string
	Value string
}

func SetHeader(name, value string) Header {
	return Header{
		Name:  name,
		Value: value,
	}
}

// NewRequest Возвращаем экземпляр запроса
func NewRequest(method, route string) *Request {
	return &Request{
		Method:  method,
		Headers: map[string]string{},
		Route:   route,
		Data:    nil,
	}
}

// Добавляем заголовки
func (r *Request) addHeader(name, value string) {
	r.Headers[name] = value
}

// SetHeader установка заголовков
func (r *Request) SetHeader(headers ...Header) *Request {
	for _, h := range headers {
		r.addHeader(h.Name, h.Value)
	}
	return r
}

// Устанавливаем формат данных и структуру передаваемых данных
/*func (r *Request) setContentTypeAndBody(contentType string, body interface{}) *Request {
	r.addHeader("Accept", contentType)
	r.addHeader("Content-Type", contentType)
	r.Data = &Data{
		ContentType: contentType,
		Body:        body,
	}
	return r
}*/

// Устанавливаем формат данных и структуру передаваемых данных
func (r *Request) setBody(contentType string, body interface{}) *Request {
	//r.addHeader("Accept", contentType)
	r.addHeader("Content-Type", contentType)
	r.Data = &Data{
		ContentType: contentType,
		Body:        body,
	}
	return r
}

// Json Body в формате Json
func (r *Request) Json(body interface{}) *Request {
	return r.setBody(JSON.String(), body)
}

// Xml Body в формате Xml
func (r *Request) Xml(body interface{}) *Request {
	return r.setBody(XML.String(), body)
}

// AddFiles Отправка файла multipart
func (r *Request) AddFiles(fieldName string, files ...string) (*Request, error) {

	if len(files) == 0 {
		return r, nil
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	defer writer.Close()
	for _, file := range files {
		err := r.openFile(fieldName, file, writer)
		if err != nil {
			return nil, err
		}
	}
	r.setBody(writer.FormDataContentType(), body)

	return r, nil
}

// Открытие файла и запись в multipart
func (r *Request) openFile(fieldName, file string, writer *multipart.Writer) error {
	// Открываем файл для чтения
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	fw, err := writer.CreateFormFile(fieldName, filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		return err
	}
	return nil
}
