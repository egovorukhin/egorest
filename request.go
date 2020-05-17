package egorest

type Request struct {
	Method  string
	Headers map[string]string
	Data    *Data
}

type Header struct {
	Name  string
	Value string
}

//Возвращаем экземпляр запроса
func NewRequest(method string) Request {
	return Request{
		Method:  method,
		Headers: map[string]string{},
		Data:    nil,
	}
}

//Добавляем заголовки
func (r Request) AddHeader(name, value string) Request {
	r.Headers[name] = value
	return r
}

//Устанавливаем заголовки
func (r Request) SetHeader(headers ...Header) Request {
	for _, h := range headers {
		r.AddHeader(h.Name, h.Value)
	}
	return r
}

//Устанавливаем формат данных и структуру передаваемых данных
func (r Request) SetBody(formatBody FormatBody, body interface{}) Request {
	r.AddHeader("Accept", formatBody.String())
	r.Data = &Data{
		FormatBody: formatBody,
		Body:       body,
	}
	return r
}
