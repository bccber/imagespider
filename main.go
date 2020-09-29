package main

import (
	"container/list"
	"fmt"
	"imagespider/conf"
	"imagespider/spiders"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill)

	defer fmt.Println("所有爬虫采集完成")

	runtime.GOMAXPROCS(runtime.NumCPU())

	spiderList := list.New()
	for _, v := range conf.Config.Spider {
		if v.Enabled != true {
			continue
		}

		spider := spiders.CreateSpider(v.Name)
		if spider == nil {
			panic(v.Name + " is nil")
		}

		// 启动爬虫
		go spider.Start(v)
		spiderList.PushBack(spider)
	}

	<-exitChan

	// 退出
	for item := spiderList.Front(); item != nil; item = item.Next() {
		item.Value.(spiders.ISpider).Stop()
	}
}
