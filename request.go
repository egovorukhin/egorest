package egorest

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

//Возвращаем экземпляр запроса
func NewRequest(method, route string) *Request {
	return &Request{
		Method:  method,
		Headers: map[string]string{},
		Route:   route,
		Data:    nil,
	}
}

//Добавляем заголовки
func (r *Request) addHeader(name, value string) {
	r.Headers[name] = value
}

func (r *Request) SetHeader(headers ...Header) *Request {
	for _, h := range headers {
		r.addHeader(h.Name, h.Value)
	}
	return r
}

//Устанавливаем формат данных и структуру передаваемых данных
func (r *Request) setBody(contentType ContentType, body interface{}) *Request {
	r.addHeader("Accept", contentType.String())
	r.addHeader("Content-Type", contentType.String())
	r.Data = &Data{
		ContentType: contentType,
		Body:        body,
	}
	return r
}

//Body в формате Json
func (r Request) Json(body interface{}) *Request {
	return r.setBody(JSON, body)
}

//Body в формате Xml
func (r Request) Xml(body interface{}) *Request {
	return r.setBody(XML, body)
}
