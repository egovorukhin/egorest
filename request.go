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
	Path    string
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
func NewRequest(path string, method ...string) *Request {
	r := &Request{
		Method: MethodGet,
		Path:   path,
	}
	if len(method) > 0 {
		r.Method = method[0]
	}
	return r
}

// Добавляем заголовки
func (r *Request) addHeader(name, value string) {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
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
func (r *Request) setBody(contentType string, body interface{}) *Request {
	r.addHeader(HeaderContentType, contentType)
	r.Data = &Data{
		ContentType: contentType,
		Body:        body,
	}
	return r
}

// JSON Body в формате Json
func (r *Request) JSON(body interface{}) *Request {
	return r.setBody(MIMEApplicationJSONCharsetUTF8, body)
}

// XML Body в формате Xml
func (r *Request) XML(body interface{}) *Request {
	return r.setBody(MIMEApplicationXMLCharsetUTF8, body)
}

// AddFiles Отправка файла multipart
func (r *Request) AddFiles(fieldName string, files ...string) (*Request, error) {

	if len(files) == 0 {
		return r, nil
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, file := range files {
		err := r.openFile(fieldName, file, writer)
		if err != nil {
			return nil, err
		}
	}
	// Закрываем Writer
	_ = writer.Close()
	r.setBody(writer.FormDataContentType(), &body)

	return r, nil
}

// Открытие файла и запись в multipart
func (r *Request) openFile(fieldName, file string, writer *multipart.Writer) error {
	// Открываем файл для чтения
	f, err := os.Open(file)
	if err != nil {
		return err
	}
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
