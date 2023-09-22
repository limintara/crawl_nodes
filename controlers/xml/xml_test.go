package xml_test

import (
	"fmt"
	"net/url"
	"testing"

	"time"

	"github.com/sanzhang007/crawl_nodes/config"
	x "github.com/sanzhang007/crawl_nodes/controlers/xml"
)

func TestXML(t *testing.T) {
	c := x.Context{}

	ct, _ := config.LoadConfig("config.json")
	c.Url = ct.Urls[0].Url
	c.Proxy, _ = url.Parse(ct.Proxy)
	c.Get()
	c.XmlUnmarshal()
	s := c.FilterCustom(func(s time.Time) bool {
		return s.After(time.Now().Add(-time.Hour * 48))
	})
	fmt.Printf("s: %v\n", s)
}
