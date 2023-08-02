package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sanzhang007/crawl_nodes/db"

	"github.com/xxf098/lite-proxy/web"
	"github.com/xxf098/lite-proxy/web/render"
	"gopkg.in/yaml.v3"
)

//go:embed "urls.yaml"
var embed_urls embed.FS

// var m map[string]any
var urls Urls

// var C = make(chan db.Nodes)
var nodes []db.Nodes

type Urls struct {
	Urls []string `yaml:"urls"`
}

func init() {
	b, _ := embed_urls.ReadFile("urls.yaml")
	yaml.Unmarshal(b, &urls)
}

var test = flag.String("test", "", "test from command line with subscription link or file")

func main() {
	flag.Parse()
	for {
		liteSpeed()
	}

}

func liteSpeed() {
	var wg sync.WaitGroup
	for _, url := range urls.Urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Printf("url: %v\n", url)
			nodes, _ := web.ParseLinks(url)
			// fmt.Printf("nodes: %v\n", nodes)
			for _, link := range nodes {
				node := db.Nodes{
					Url:        url,
					Link:       link,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				var dbNode db.Nodes
				db.DB.Where("link= ?", node.Link).Find(&dbNode)
				if dbNode.Id == 0 {
					db.DB.Create(&node)
				}
			}

		}(ChangeDate(url))

	}
	wg.Wait()

	configPath := "config.json"

	if *test == "all" {
		db.DB.Where("(fail_count = 0 and success_count = 0) or (fail_count<5 and success_count>=1)").Find(&nodes)
		fmt.Printf("本次测速数量: %v\n", len(nodes))
		time.Sleep(time.Second * 5)
		testAsyncContext(nodes, &configPath)
	} else {
		db.DB.Where("(fail_count = 0 and success_count = 0) or (fail_count<5 and success_count>=1)").Find(&nodes)
		testFromCMD(nodes, &configPath)
	}
}

func testFromCMD(nodes []db.Nodes, configPath *string) {
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

func testAsyncContext(nodes []db.Nodes, configPath *string) {
	ctx := context.Background()
	var builder strings.Builder
	for _, node := range nodes {
		builder.WriteString(node.Link)
		builder.WriteString(" ")
	}
	options, _ := web.ReadConfig(*configPath)
	options.Subscription = builder.String()
	c, _, err := web.TestAsyncContext(ctx, *options)
	// c, err := web.TestContext(ctx, *options, &web.OutputMessageWriter{})
	// fmt.Printf("s: %v\n", s)
	// fmt.Printf("s: %v\n", err)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for item := range c {
		var dbNode db.Nodes
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

func writeMsg(index int, node db.Nodes, n render.Node, m bool) {
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
