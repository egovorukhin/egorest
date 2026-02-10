package egorest

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const VERSION = "0.3.10"

type Client struct {
	Config  Config
	BaseUrl *url.URL
	ctx     *context.Context
	client  *http.Client
}

func NewClient(config Config) (c *Client) {
	return &Client{
		Config: config,
	}
}

// WithContext Можем использовать контекст
func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = &ctx
	return c
}

// SetBasicAuth Устанавливаем Basic авторизацию
func (c *Client) SetBasicAuth(name, password string) *Client {
	c.Config.BasicAuth = SetBasicAuth(name, password)
	return c
}

// SetProxy Устанавливаем прокси сервер
func (c *Client) SetProxy(proxy string) *Client {
	c.Config.Proxy, _ = url.Parse(proxy)
	return c
}

// SetTimeout Устанавливаем таймаут соединения
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.Config.Timeout = timeout
	return c
}

// SetPath Добавляем маршрут, он будет подставляться перед маршрутами указанными в Request
func (c *Client) SetPath(path string) *Client {
	if c.BaseUrl != nil {
		c.BaseUrl.Path = path
	}
	return c
}

// SetHttpClient Установка http клиента
func (c *Client) SetHttpClient(client *http.Client) *Client {
	c.client = client
	return c
}

// Замена пустых значений на %20
func (*Client) trim(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}

// Формируем строк для http запроса
func (c *Client) url(path string) string {

	if path != "" {
		if path[0] != '/' && path[0] != '?' {
			path = "/" + path
		}
	}

	return c.trim(c.BaseUrl.String() + path)
}

// Send Отправляем запрос на сервер
func (c *Client) Send(r *Request) (resp *http.Response, err error) {

	//Преобразуем структуру в набор байт для отправки
	var body io.Reader
	if r.Data != nil {
		body, err = r.Data.marshal()
		if err != nil {
			return
		}
		r.setHeader(HeaderContentType, r.Data.ContentType)
	}

	// Устанавливаем базовый линк
	c.BaseUrl, err = c.Config.BaseUrl.getUrl()
	if err != nil {
		return
	}

	req, err := http.NewRequest(r.Method, c.url(r.Path), body)
	if err != nil {
		return
	}

	//На случай авторизации (Только Basic)
	if c.Config.BasicAuth != nil {
		req.SetBasicAuth(c.Config.BasicAuth.Name, c.Config.BasicAuth.Password)
	}

	//Добавляем заголовки
	req.Header.Add("User-Agent", "EgoRest "+VERSION)
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	// Если клиент не назначен, то создаем по конфигам
	if c.client == nil {
		write, read := 0, 0
		if c.Config.Buffers != nil {
			write, read = c.Config.Buffers.get()
		}
		c.client = &http.Client{
			Transport: &http.Transport{
				//Прокси сервер, если nil то не используем
				Proxy: http.ProxyURL(c.Config.Proxy),
				//Это на случай не подтверждённого сертификата (ОПАСНО)
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.Config.Secure,
				},
				WriteBufferSize: write,
				ReadBufferSize:  read,
			},
			Timeout: c.Config.Timeout,
		}
	}

	// Поехали с контекстом...
	if c.ctx != nil {
		ctx := *c.ctx
		resp, err = c.client.Do(req.WithContext(ctx))
		if err != nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default:
			}
		}
		return
	}

	// Просто поехали...
	return c.client.Do(req)
}

// Execute Отправка данных на сервер, ждём в ответе какую то структуру
func (c *Client) Execute(r *Request, responseBody interface{}, handler ...UnmarshalHandler) error {

	//Отправляем запрос
	resp, err := c.Send(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	//Получаем из ответа набор байт
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	//Если 200 значит десериализуем данные
	case 200:
		//Переводим все это дело в структуру,
		//но сначала находим в каком формате данные
		return ContentType(resp.Header.Get(HeaderContentType)).unmarshal(body, &responseBody, handler...)
	default:
		//Если другие ошибки, то body возвращаем в виде ошибки
		return errors.New(string(body))
	}
}

// Get вызывается функция Send, только с методом Get
func (c *Client) Get(r *Request) (*http.Response, error) {
	r.Method = MethodGet
	return c.Send(r)
}

// Post вызывается функция Send, только с методом Post
func (c *Client) Post(r *Request) (*http.Response, error) {
	r.Method = MethodPost
	return c.Send(r)
}

// Put вызывается функция Send, только с методом Put
func (c *Client) Put(r *Request) (*http.Response, error) {
	r.Method = MethodPut
	return c.Send(r)
}

// Delete вызывается функция Send, только с методом Delete
func (c *Client) Delete(r *Request) (*http.Response, error) {
	r.Method = MethodDelete
	return c.Send(r)
}

// ExecuteGet вызывается функция Execute, только с методом Get
func (c *Client) ExecuteGet(r *Request, responseBody interface{}, handler ...UnmarshalHandler) error {
	r.Method = MethodGet
	return c.Execute(r, responseBody, handler...)
}

// ExecutePost вызывается функция Execute, только с методом Post
func (c *Client) ExecutePost(r *Request, responseBody interface{}, handler ...UnmarshalHandler) error {
	r.Method = MethodPost
	return c.Execute(r, responseBody, handler...)
}

// ExecutePut вызывается функция Execute, только с методом Put
func (c *Client) ExecutePut(r *Request, responseBody interface{}, handler ...UnmarshalHandler) error {
	r.Method = MethodPut
	return c.Execute(r, responseBody, handler...)
}

// ExecuteDelete вызывается функция Execute, только с методом Delete
func (c *Client) ExecuteDelete(r *Request, responseBody interface{}, handler ...UnmarshalHandler) error {
	r.Method = MethodDelete
	return c.Execute(r, responseBody, handler...)
}
