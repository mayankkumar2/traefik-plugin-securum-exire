package traefik_plugin_securum_exire

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Config struct {
}


func CreateConfig() *Config {
	return &Config{}
}

type SecurumExire struct {
	http.Handler
	next http.Handler
	name string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Println("OK runner")
	return &SecurumExire{
		next: next,
		name: name,
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
	e.contentLength = len(b)
	return e.buffer.Write(b)
}

func (e *SecurumExireWriter) WriteHeader(statusCode int) {
	e.statusCode = statusCode
}

func (e *SecurumExire) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	respWr := &SecurumExireWriter{
		ResponseWriter: rw,
		p:              e,
	}
	e.next.ServeHTTP(respWr, req)
	cl := respWr.contentLength
	var b = respWr.buffer.Bytes()
	if len(b) > 100 {
		fmt.Println("Length more than 100")
		responseString := []byte("Length More than 100")
		fmt.Println("length written: ", len(responseString))
		cl = len(responseString)
	} else {
		fmt.Println("Length less than 100")
		l, _ := rw.Write(b)
		fmt.Println("length written: ", l)
	}
	fmt.Println(rw.Header().Values("content-length"))
	rw.Header().Del("content-length")
	rw.Header().Set("content-length", strconv.Itoa(cl))
	fmt.Println(rw.Header().Values("content-length"))
	rw.WriteHeader(respWr.statusCode)
	responseString := []byte("Length More than 100")
	_, _ = rw.Write(responseString)
}