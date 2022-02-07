package egorest

import (
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	BaseUrl   BaseUrl
	Secure    bool
	Timeout   time.Duration
	Buffers   *Buffers
	Proxy     *url.URL
	BasicAuth *BasicAuth
}

type Buffers struct {
	Write int
	Read  int
}

type BaseUrl struct {
	Url    string
	Schema string
	Host   string
	Port   int
	Path   string
}

func (b BaseUrl) getUrl() (*url.URL, error) {
	if b.Url != "" {
		return url.Parse(b.Url)
	}

	port := ""
	if b.Port > 0 {
		port = ":" + strconv.Itoa(b.Port)
	}

	return &url.URL{
		Scheme: b.Schema,
		Host:   b.Host + port,
		Path:   b.Path,
	}, nil
}

func (b *Buffers) get() (int, int) {
	return b.Write, b.Read
}
