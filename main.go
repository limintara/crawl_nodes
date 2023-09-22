package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sanzhang007/crawl_nodes/controlers/xml"
	"github.com/sanzhang007/crawl_nodes/db"
	"github.com/sanzhang007/crawl_nodes/parse"
	clash "github.com/sanzhang007/webgin/protocol"
	"github.com/xxf098/lite-proxy/web"
	"github.com/xxf098/lite-proxy/web/render"
	"gopkg.in/yaml.v3"
)

var urls Urls

// var C = make(chan db.Node1)
var nodes []db.Node1

type Urls struct {
	Urls []string `yaml:"urls"`
}

func loadConfig() {
	// b, _ := embed_urls.ReadFile("urls.yaml")
	b, err := os.ReadFile("urls.yaml")
	if err != nil {
		log.Panic(err)
	}
	yaml.Unmarshal(b, &urls)
}

var test = flag.String("test", "", "test from command line with subscription link or file")

func main() {
	flag.Parse()

	loadConfig()
	s := xml.XmlUrls("config.json")
	// s = append(s, parse.GetChangfengpaths()...)
	for _, item := range s {
		exist := func(s1 string, s2 []string) bool {
			for _, item := range s2 {
				if s1 == item {
					return true
				}
			}
			return false
		}
		if !exist(item, urls.Urls) {
			urls.Urls = append(urls.Urls, item)
		}
	}
	urls.Urls = append(urls.Urls, parse.GetChangfengpaths()...)
	fmt.Println(urls.Urls)
	fmt.Printf("count total url is %d\n", len(urls.Urls))
	time.Sleep(time.Second * 5)
	liteSpeed()

}

func WriteFile(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func liteSpeed() {
	var wg sync.WaitGroup
	// urls.Urls = append(urls.Urls, parse.GetChangfengpaths()...)

	for _, url := range urls.Urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Printf("url: %v\n", url)
			nodes, _ := web.ParseLinks(url)

			for _, link := range nodes {
				if strings.HasPrefix(link, "vless://") {
					WriteFile("tmp.yaml", []byte(link+"\n"), 0666)

				}
				node := db.Node1{
					Url:        url,
					Link:       link,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				var dbNode db.Node1
				// db.DB.Where("link= ?", node.Link).Find(&dbNode)

				b, err := GetPotocolServerAndPortAndPassword(node.Link)
				if err != nil {
					continue
				}
				if b == nil {
					continue
				}

				node.Link1 = string(b)
				// db.DB.Where("link1= ?", string(b)).Find(&dbNode)
				db.DB.Where("link1= ?", string(b)).Find(&dbNode)
				// b := []byte("ss")

				if dbNode.Id == 0 {
					// fmt.Printf("db.Node1: %v\n", dbNode.Link1)
					db.DB.Create(&node)
				}
			}
		}(ChangeDate(url))

	}
	wg.Wait()
	configPath := "configPing.json"

	if *test == "all" {
		// db.DB.Where("(fail_count = 0 and success_count = 0) or fail_count<20 ").Find(&nodes)
		db.DB.Find(&nodes)
		fmt.Printf("本次测速数量: %v\n", len(nodes))
		time.Sleep(time.Second * 5)
		testAsyncContext(nodes, &configPath)
	} else {
		db.DB.Where("(fail_count = 0 and success_count = 0) or (fail_count<5 and success_count>=1)").Find(&nodes)
		testFromCMD(nodes, &configPath)
	}
}

func testFromCMD(nodes []db.Node1, configPath *string) {
	for index, node := range nodes {
		n, err := web.TestFromCMD(node.Link, configPath)
		if n == nil {
			continue
		}
		if err != nil {
			log.Printf("[error %s]:  %v", err, node)
			writeMsg(index, node, n[0], false)
			continue
		}
		if n[0].Ping == "0" {
			writeMsg(index, node, n[0], false)
			continue
		}
		writeMsg(index, node, n[0], true)
	}
}

func testAsyncContext(nodes []db.Node1, configPath *string) {
	ctx := context.Background()
	var builder strings.Builder
	for _, node := range nodes {
		builder.WriteString(node.Link)
		builder.WriteString(" ")
	}
	options, _ := web.ReadConfig(*configPath)
	options.Subscription = builder.String()
	// c, _, err := web.TestAsyncContext(ctx, *options)
	n, err := web.TestContext(ctx, *options, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for _, item := range n {
		var dbNode db.Node1
		db.DB.Where("link= ?", item.Link).Find(&dbNode)
		if item.Ping == "0" {
			writeMsg(0, dbNode, item, false)
		} else {
			writeMsg(0, dbNode, item, true)
		}
		fmt.Printf("item: %v\n", item)
		if item.Id == len(nodes)-1 {
			break
		}
	}
}

var mux sync.Mutex

func writeMsg(index int, node db.Node1, n render.Node, m bool) {
	mux.Lock()
	fmt.Println("###########################################################################################################################################")
	fmt.Printf("index: %v\n", index)
	fmt.Println("Id		Ping	AvgSpeed	MaxSpeed	FailCount	SuccessCount	CreateTime")
	fmt.Println(node.Id, "\t\t", node.Ping, "\t", node.AvgSpeed, "\t", node.MaxSpeed, "\t\t", node.FailCount, "\t", node.SuccessCount, "\t", node.CreateTime.Format("2006/01/02 15:04:05"))
	if m {
		node.SuccessCount += 1
	} else {
		node.FailCount += 1
	}
	db.DB.Model(&node).Update("max_speed", n.MaxSpeed).Update("success_count", node.SuccessCount).Update("fail_count", node.FailCount).Update("ping", n.Ping).Update("avg_speed", n.AvgSpeed).Update("update_time", time.Now())
	db.DB.Find(&node)
	fmt.Println(node.Id, "\t\t", node.Ping, "\t", node.AvgSpeed, "\t", node.MaxSpeed, "\t\t", node.FailCount, "\t", node.SuccessCount, "\t", node.CreateTime.Format("2006/01/02 15:04:05"))
	fmt.Println("###########################################################################################################################################")
	mux.Unlock()
}

func ChangeDate(str string) string {

	now := time.Now()

	str1 := strings.Replace(str, "{year}", now.Format("2006"), -1)
	str2 := strings.Replace(str1, "{month}", now.Format("01"), -1)
	str3 := strings.Replace(str2, "{day}", now.Format("02"), -1)
	str4 := strings.Replace(str3, "{month_s}", now.Format("1"), -1)
	str5 := strings.Replace(str4, "{day_s}", now.Format("2"), -1)
	return str5
}

type ServerAndPortAndPassword struct {
	Server   string
	Port     string
	Password string
}

func GetPotocolServerAndPortAndPassword(link string) ([]byte, error) {
	mux.Lock()
	c, err := clash.ClashParse(&link)
	if err != nil {
		return nil, nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, nil
	}
	var s ServerAndPortAndPassword
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, nil
	}
	mux.Unlock()
	return json.Marshal(s)
}

func Testxxx(t *testing.T) {
	b, err := GetPotocolServerAndPortAndPassword(`ss://YWVzLTI1Ni1jZmI6cXdlclJFV1FAQEAyMjEuMTUwLjEwOS41Ojk1NTU=#%F0%9F%87%B0%F0%9F%87%B7KR_439
	`)
	fmt.Println(string(b))
	fmt.Println(err)

}
