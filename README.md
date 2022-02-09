# Rest клиент для проектов на golang
- [Описание](README.md#описание)
- [Установка](README.md#установка)
- Использование
    - [Конфигурация](README.md#конфигурация)
    - [Инициализация клиента](README.md#инициализация-клиента)
    - [Инициализация запроса](README.md#инициализация-запроса)
    - [Отправка запроса](README.md#отправка-запроса)
    - [Десириализация ответа по заголовку Content-Type](README.md#десириализация-ответа-по-заголовку-content-type)
    - [Использование http.Client](README.md#использование-httpclient)
    - [Отправка файла](README.md#отправка-файла)
- [Примеры](README.md#примеры)
- [Лицензия](README.md#лицензия)

## Описание
__EgoRest__ - это http клиент, который подключается в виде пакета в проект на golang. В запросе нужно передать `url` и тело запроса при необходимости, и отправьте запрос. Есть возможность как получить ответ в виде структуры `http.*Response`, так и десериализовать тело ответа в структуру основываясь на заголовке `Content-Type`

## Установка
```
go get -u github.com/egovorukhin/egorest
```

## Использование

### Конфигурация
```golang
	cfg := egorest.Config{
		BaseUrl: egorest.BaseUrl{
			Schema: "https",
			Host:   "localhost",
			Port:   5005,
			Path:   "api",
		},
		Secure:  true,
		Timeout: time.Second * 30,
		Buffers: &egorest.Buffers{
			Write: 4096,
			Read:  4096,
		},
		Proxy:     nil,
		BasicAuth: nil,
	}
```
`BaseUrl` - формирование базового адреса запроса.  
`Secure` - флаг от которого зависит неужно ли проверять сертификат TLS.  
`Timeout` - таймаут запроса по истечении которого соединение закроется не зависимо от результата.  
`Buffers` - входной и выходной буфер соединения.  
`Proxy` - нужно указать структуру `*url.Url`.  
`BasicAuth` - базовая авторизация запроса, указывается имя и пароль.
### Инициализация клиента
```golang
    client := egorest.NewClient(cfg)
```
`NewClient` - инициализирует клиента, который получет структуру `Config` описанную выше.
### Инициализация запроса
```golang
    request := egorest.NewRequest("api/rec/audio")
    request = egorest.NewRequest("api/rec/audio", "GET")
```
`NewRequest` - инициализирует запрос, принимается переменная `path `часть пути адреса в виде строки.  
`Method` - необязательный аргумент, по умолчанию ставится `GET`.
### Отправка запроса
```golang
    client := egorest.NewClient(cfg)
    request := egorest.NewRequest("api/rec/audio")
    resp, _ := client.Send(request)
    resp, _ := client.Get(request)
    resp, _ = client.Post(request)
    resp, _ = client.Put(request)
    resp, _ = client.Delete(request)
```
`Send` - отправить запрос.  
`Get` - отправить запрос используя метод `GET`.  
`Post` - отправить запрос используя метод `POST`.  
`Put` - отправить запрос используя метод `PUT`.  
`Delete` - отправить запрос используя метод `DELETE`.  
`SetBody` - установка тела запроса.  
`JSON` - установка тела запроса в формате `json` и установкой заголовка `Content-Type: "application/json; charset=utf-8"`.  
`XML` - установка тела запроса в формате `xml` и установкой заголовка `Content-Type: "application/xml; charset=utf-8"`.  

### Десириализация ответа по заголовку Content-Type
 ```golang
    user := User{}
	request.SetBody(egorest.MIMEApplicationJSON, user)
	resp, _ = client.Post(request)
	requset.JSON(user)
	resp, _ = client.Post(request)
	request.XML(user)
	resp, _ = client.Post(request)
	_ = client.Execute(request, &user)
	_ = client.ExecuteGet(request, &user)
	_ = client.ExecutePost(request, &user)
	_ = client.ExecutePut(request, &user)
	_ = client.ExecuteDelete(request, &user)
    _ = client.Execute(request, &user, func(contentType string, data []byte, v interface{}) error {
		if contentType == "application/json" {
			return json.Unmarshal(data, v)
		}
		return nil
	})
 ```
`Execute` - отправрка запроса и десериализация тела ответа на основе заголовка `Content-Type`.  
`ExecuteGet` - отправрка запроса с методом `GET` и десериализация тела ответа на основе заголовка `Content-Type`.  
`ExecutePost` - отправрка запроса с методом `POST` и десериализация тела ответа на основе заголовка `Content-Type`.  
`ExecutePut` - отправрка запроса с методом `PUT` и десериализация тела ответа на основе заголовка `Content-Type`.  
`ExecuteDelete` - отправрка запроса с методом `DELETE` и десериализация тела ответа на основе заголовка `Content-Type`.  
В функции `Execute...` есть необязательный параметр, который принимает функцию `UnmarshalHandler`. Функция должна реализовать десериализацию данных не специфичных форматов, имеет вид `func(contentType string, data []byte, v interface{}) error`, где `contentType string` возвращает значение заголовка `Content-Type`, срез байт `data []byte` и структуру `v interface{}` которую нужно заполнить.
### Использование http.Client
```golang
    httpclient := &http.Client{}
    client := NewClient(cfg).SetHttpClient(httpclient)
```
При необходимости можно настроить `*http.Client` самостоятельно и передать в `egorest.Client` с помощью функции `SetHttpClient(client *http.Client) *Client` и тогда запрос будет отправляться вашим клиентом.
### Отправка файла
```golang
	r := egorest.NewRequest("api/images")
	err := r.AddFiles("image", "/images/img1.png", "/images/img2.png")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = NewClient(cfg).Post(r)
```
Для передачи файлов на сервер используйте функцию `AddFiles`, которая принимает `fieldName string` имя поля и список файлов в виде пути `files ...string`. Если возникли проблемы с открытием или чтением хоть одного из файлов вернется ошибка

## Примеры
- [Все](https://github.com/egovorukhin/egorest/tree/master/examples)
- [Десериализовать тело ответа](https://github.com/egovorukhin/egorest/tree/master/examples/get_json)
- [Отправить запрос с телом](https://github.com/egovorukhin/egorest/tree/master/examples/send_body)
- [Отправить файл](https://github.com/egovorukhin/egorest/tree/master/examples/add_file)
- [Обработчик для десериализации](https://github.com/egovorukhin/egorest/tree/master/examples/custom_handler)

## Лицензия
Пользуйтесь на здоровье
