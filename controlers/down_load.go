package controlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type Context struct {
	Url   string
	Proxy *url.URL
	Err   error
	Buf   []byte
	Link  string
}

var EmptyContext Context = Context{}

func (c *Context) Get() {
	var client http.Client
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(c.Proxy),
	}
	r, err := client.Get(c.Url)
	if err != nil {
		c.Err = err
		return
	}
	body := r.Body
	defer body.Close()
	c.Buf, c.Err = io.ReadAll(body)
}

func (c Context) String() string {
	return string(string(c.Buf))
}

func (c *Context) Parse(regStr, content string) {
	if c.Buf == nil {
		return
	}
	r := regexp.MustCompile(regStr)
	if len(r.FindAllString(content, 1)) == 1 {
		c.Link = r.FindAllString(content, 1)[0]
	} else {
		c.Err = errors.New("indexOut")
	}
}
