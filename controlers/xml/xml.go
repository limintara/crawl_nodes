package xml

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

	"github.com/sanzhang007/crawl_nodes/config"
	"github.com/sanzhang007/crawl_nodes/controlers"
)

type Context struct {
	controlers.Context
	Custom
}

func (c *Context) XmlUnmarshal() {
	err := xml.Unmarshal(c.Buf, &(c.Custom))
	c.Err = err
}

func (c *Context) FilterCustom(condition func(time.Time) bool) []string {
	// var s []string
	s := make([]string, len(c.Contents))
	i := 0
	for _, item := range c.Contents {
		if condition(item.LastModified) {
			s[i] = c.Url + item.Key
			i++
		}
	}
	return s[:i]
}

type Custom struct {
	XMLName  xml.Name
	Contents []Content
}

type Content struct {
	Key          string
	LastModified time.Time
	Size         int
}

func XmlUrls(configFile string) []string {
	ct, _ := config.LoadConfig(configFile)
	var ss []string
	for _, item := range ct.Urls {

		c := Context{}
		c.Url = item.Url
		c.Proxy, _ = url.Parse(ct.Proxy)
		c.Get()
		c.XmlUnmarshal()
		s := c.FilterCustom(func(s time.Time) bool {
			return s.After(time.Now().Add(-time.Hour * 48))
		})
		fmt.Printf("s: %v\n", s)
		ss = append(ss, s...)
	}
	return ss
}
