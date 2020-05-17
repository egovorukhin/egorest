package egorest

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const VERSION = "0.1.1"

const (
	GET = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"
	HEAD = "HEAD"
	PATCH = "PATCH"
	OPTIONS = "OPTIONS"
	TRACE = "TRACE"
	CONNECT = "CONNECT"
)

type Client struct {
	Hostname string
	Port int
	Secure bool
	Timeout int
	Proxy *url.URL
	BasicAuth *BasicAuth
	Request Request
}

//Создаём новый экземпляр Client
func NewClient(hostname string, port int, secure bool) Client {
	return Client{
		Hostname: hostname,
		Port:     port,
		Secure:   secure,
		Proxy: nil,
		BasicAuth: nil,
		Timeout: 30,
	}
}

//Устанавливаем заголовки
func (client Client) SetRequest(r Request) Client {
	client.Request = r
	return client
}

//Учтанавливаем Basic авторизацию
func (client Client) SetBasicAuth(name, password string) Client {
	client.BasicAuth = SetBasicAuth(name, password)
	return client
}

//Устанавливаем прокси сервер
func (client Client) SetProxy(proxy string) Client {
	client.Proxy, _ = url.Parse(proxy)
	return client
}

//Устанавливаем таймаут соединения
func (client Client) SetTimeout(timeout int) Client {
	client.Timeout = timeout
	return client
}

//Формируем строк для http запроса
func (client Client) url() string {
	s := "http"
	if client.Secure {
		s = "https"
	}
	return fmt.Sprintf("%s://%s:%d/", s, client.Hostname, client.Port)
}

//Отправка данных на сервер,
func (client Client) Send(route string, responseBody interface{}) error {

	//Преобразуем сируктуру в набор байт для отправки
	var body []byte
	if client.Request.Body != nil {
		var err error
		body, err = client.Request.Body.marshal()
		if err != nil {
			return err
		}
	}

	//Инициализируем клиента и передаём все необходимое
	httpClient := &http.Client{
		Timeout: time.Duration(client.Timeout) * time.Second,
		Transport: &http.Transport{
			//Это на случай не подтверждённого сертификата (ОПАСНО)
			TLSClientConfig:        &tls.Config{
				InsecureSkipVerify:          true,
			},
			//Прокси сервер, если nil то не используем
			Proxy: http.ProxyURL(client.Proxy),
		},
	}

	req, err := http.NewRequest(
		client.Request.Method,
		client.url() + route,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	//На случай аторизации (Только Basic)
	if client.BasicAuth != nil {
		req.SetBasicAuth(client.BasicAuth.Name, client.BasicAuth.Password)
	}

	//Добавляем заголовки
	req.Header.Add("User-Agent", "EgoRest/" + VERSION)
	for key, value := range client.Request.Headers {
		req.Header.Add(key, value)
	}

	//Поехали...
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	//Получаем из ответа набор байт
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Если 200 значит десериализуем данные
	switch resp.StatusCode {
	case 200:
		//Переводим все это дело в структуру
		err = GetFormatBody(resp.Header.Get("Content-Type")).unmarshal(body, &responseBody)
		if err != nil {
			return err
		}
		return nil
		//Если другие ошибки то body возвращаем в виде ошибки
	default:
		return errors.New(string(body))
	}
}
