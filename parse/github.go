package parse

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/sanzhang007/NodeShare/speed/urlsp"
)

func changfeng(u string, reg string) []string {
	c := urlsp.Context{
		Url: u,
		Proxy: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:7890",
		},
	}
	c.Get()
	r := regexp.MustCompile(reg)
	return r.FindAllString(string(c.Buf), -1)

}

func GetChangfengpaths() []string {
	s2 := time.Now().Format("2006_01_02")
	s := changfeng(fmt.Sprintf("https://github.com/changfengoss/pub/tree/main/data/%s", s2),
		fmt.Sprintf(`{"name":"[a-zA-Z0-9]+.yaml","path":"data/%s/[a-zA-Z0-9]+.(yaml|txt)","contentType":"file"}`, s2))
	var m map[string]string
	temps := make([]string, len(s))
	for index, item := range s {
		json.Unmarshal([]byte(item), &m)
		temps[index] = "https://raw.githubusercontent.com/changfengoss/pub/main/" + m["path"]
	}
	fmt.Println(temps)
	return temps
}
