package traefik_plugin_securum_exire

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Url string `json:"url" yaml:"url" toml:"url"`
}


func CreateConfig() *Config {
	return &Config{}
}

type SecurumExire struct {
	http.Handler
	next http.Handler
	name string
	config *Config
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// fmt.Println(config)
	return &SecurumExire{
		next: next,
		name: name,
		config: config,
	}, nil
}

type SecurumExireWriter struct {
	http.ResponseWriter
	buffer     bytes.Buffer
	overridden bool
	p *SecurumExire
	statusCode int
	contentLength int
}

func (e *SecurumExireWriter) Header() http.Header {
	return e.ResponseWriter.Header()
}

func (e *SecurumExireWriter) Write(b []byte) (int, error)  {
	e.contentLength += len(b)
	return e.buffer.Write(b)
}

func (e *SecurumExireWriter) WriteHeader(statusCode int) {
	e.statusCode = statusCode
}

func(e *SecurumExire) CheckIfLeak(endpoint string, b io.Reader) (bool, error) {
	if e.config.Url == "" {
		return false, errors.New("error: plugin not configured properly")
	}
	urlEndpoint, err := url.Parse(fmt.Sprintf("http://%s/check", e.config.Url))
	if e.config.Url == "" {
		return false, errors.New("error: plugin not configured properly")
	}
	var req = &http.Request{
		URL: urlEndpoint,
		Header: map[string][]string{
			"endpoint": {endpoint},
			"content-type": {"application/json"},
		},
		Body: ioutil.NopCloser(b),
		Method: http.MethodPost,
	}
	cl := http.Client{
		Timeout: time.Minute*5,
	}
	response, err := cl.Do(req)
	if err != nil {
		return false, err
	} else {
		return response.StatusCode == http.StatusForbidden, nil
	}
}

func (e *SecurumExire) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	isBlocked, err := e.CheckIfBlocked(req.URL.Path)
	if err != nil {
		log.Println(err)
		e.next.ServeHTTP(rw, req)
		return
	}
	respWr := &SecurumExireWriter{
		ResponseWriter: rw,
		p:              e,
	}
	if !isBlocked {
		e.next.ServeHTTP(respWr, req)
	} else {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	cl := respWr.contentLength
	buff := new(bytes.Buffer)
	tee := io.TeeReader(&respWr.buffer, buff)
	var response []byte
	var code int
	isLeak, _ := e.CheckIfLeak(req.URL.Path, tee)
	if isLeak {
		response = []byte("{}")
		code = http.StatusForbidden
	} else {
		response, _ = ioutil.ReadAll(buff)
		code = respWr.statusCode
	}
	cl = len(response)
	//fmt.Println("length written: ", cl)
	//fmt.Println(rw.Header().Values("content-length"))
	rw.Header().Del("content-length")
	rw.Header().Set("content-length", strconv.Itoa(cl))
	// fmt.Println(rw.Header().Values("content-length"))
	_, _ = rw.Write(response)
	rw.WriteHeader(code)
}

func (e *SecurumExire) CheckIfBlocked(endpoint string) (bool, error) {
	if e.config.Url == "" {
		return false, errors.New("error: plugin not configured properly")
	}
	urlEndpoint, _ := url.Parse(fmt.Sprintf("http://%s/check_endpoint", e.config.Url))
	var req = &http.Request{
		URL: urlEndpoint,
		Header: map[string][]string{
			"endpoint": {endpoint},
		},
		Method: http.MethodGet,
	}
	cl := http.Client{
		Timeout: time.Minute*5,
	}
	response, err := cl.Do(req)
	if err != nil {
		return false, err
	} else {
		return response.StatusCode == http.StatusForbidden, nil
	}
}