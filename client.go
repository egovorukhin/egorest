package egorest

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const VERSION = "0.2.14"

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
	CONNECT = "CONNECT"
)

type Client struct {
	Url        *url.URL
	Hostname   string
	Port       int
	Secure     bool
	Timeout    int
	Route      string
	Proxy      *url.URL
	BasicAuth  *BasicAuth
	ctx        *context.Context
	httpClient *http.Client
}

// NewClient Создаём новый экземпляр Client
func NewClient(hostname string, port int, secure bool) *Client {
	return &Client{
		Hostname:  hostname,
		Port:      port,
		Secure:    secure,
		Proxy:     nil,
		BasicAuth: nil,
		Timeout:   30,
	}
}

func NewClientByUri(uri string) (*Client, error) {
	Url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return &Client{
		Url:       Url,
		Timeout:   30,
		Proxy:     nil,
		BasicAuth: nil,
	}, nil
}

// WithContext Можем использовать контекст
func (client *Client) WithContext(ctx context.Context) *Client {
	client.ctx = &ctx
	return client
}

func (client *Client) SetHttpClient(c *http.Client) *Client {
	client.httpClient = c
	return client
}

// SetBasicAuth Устанавливаем Basic авторизацию
func (client *Client) SetBasicAuth(name, password string) *Client {
	client.BasicAuth = SetBasicAuth(name, password)
	return client
}

// SetProxy Устанавливаем прокси сервер
func (client *Client) SetProxy(proxy string) *Client {
	client.Proxy, _ = url.Parse(proxy)
	return client
}

// SetTimeout Устанавливаем таймаут соединения
func (client *Client) SetTimeout(timeout int) *Client {
	client.Timeout = timeout
	return client
}

// SetRoute Добавляем маршрут, он будет подставляться перед маршрутами указанными в Request
func (client *Client) SetRoute(route string) *Client {
	client.Route = route
	return client
}

func (Client) trim(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}

//Формируем строк для http запроса
func (client *Client) url(route string) string {

	if client.Url != nil {
		return client.trim(client.Url.String())
	}

	s := "http"
	if client.Secure {
		s = "https"
	}
	if route != "" {
		if route[0] != '/' && route[0] != '?' {
			route = "/" + route
		}
	}
	if client.Route != "" {
		if client.Route[0] != '/' {
			client.Route = "/" + client.Route
		}
		if len(client.Route) > 0 && client.Route[len(client.Route)-1] == '/' && route != "" {
			client.Route = client.Route[:len(client.Route)-1]
		}
	}
	//Пробел меняем на спец символ
	route = client.trim(route)
	client.Route = client.trim(client.Route)

	return fmt.Sprintf("%s://%s:%d%s%s", s, client.Hostname, client.Port, client.Route, route)
}

// Send Отправляем запрос на сервер
func (client *Client) Send(r *Request) (*http.Response, error) {

	//Преобразуем сируктуру в набор байт для отправки
	var body []byte
	if r.Data != nil {
		var err error
		body, err = r.Data.marshal()
		if err != nil {
			return nil, err
		}
	}

	// Инициализируем клиента и передаём все необходимое.
	// Если http клиент пустой, то определяем его
	if client.httpClient == nil {
		client.httpClient = &http.Client{
			Transport: &http.Transport{
				//Это на случай не подтверждённого сертификата (ОПАСНО)
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				//Прокси сервер, если nil то не используем
				Proxy: http.ProxyURL(client.Proxy),
			},
			Timeout: time.Duration(client.Timeout) * time.Second,
		}
	}

	req, err := http.NewRequest(
		r.Method,
		client.url(r.Route),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	//На случай аторизации (Только Basic)
	if client.BasicAuth != nil {
		req.SetBasicAuth(client.BasicAuth.Name, client.BasicAuth.Password)
	}

	//Добавляем заголовки
	req.Header.Add("User-Agent", "EgoRest/"+VERSION)
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	//Поехали...
	if client.ctx != nil {
		ctx := *client.ctx
		resp, err := client.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default:
			}
		}
		return resp, err
		//return ctxhttp.Do(*client.ctx, client.httpClient, req)
	}
	return client.httpClient.Do(req)
}

// Execute Отправка данных на сервер, ждём в ответе какую то структуру
func (client *Client) Execute(r *Request, responseBody interface{}) error {

	//Отправляем запрос
	resp, err := client.Send(r)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	//Получаем из ответа набор байт
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Если 200 значит десериализуем данные
	switch resp.StatusCode {
	case 200:
		//Переводим все это дело в структуру,
		//но сначала находим в каком формате данные
		err = getFormatBody(resp.Header.Get("Content-Type")).unmarshal(body, &responseBody)
		if err != nil {
			return err
		}
		return nil
		//Если другие ошибки то body возвращаем в виде ошибки
	default:
		return errors.New(string(body))
	}
}
