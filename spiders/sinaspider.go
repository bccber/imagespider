package spiders

import (
	"encoding/json"
	"fmt"
	"imagespider/db"
	"imagespider/models"
	"imagespider/utils"
	"time"
)

// 新浪gif爬虫
type SinaSpider struct {
	bExit   bool
	maxChan chan int
}

// 采集
func (p *SinaSpider) collect(spiderInfo models.SpiderInfo, pageIndex int) {
	defer func() {
		<-p.maxChan
		fmt.Printf("SinaSpider 第 %d 页采集完毕\n", pageIndex)
	}()

	fmt.Printf("SinaSpider 正在采集第 %d 页\n", pageIndex)

	strURL := fmt.Sprintf(spiderInfo.URL, pageIndex)
	strHtml := utils.DownUrlHtml(strURL)
	if strHtml == "" {
		fmt.Println("SinaSpider strHtml == ''")
		return
	}

	info := models.SinaJsonInfo{}
	err := json.Unmarshal([]byte(strHtml), &info)
	if err != nil || info.Data == nil || len(info.Data) <= 0 {
		fmt.Println("SinaSpider info.Data == nil")
		return
	}

	for _, item := range info.Data {
		// 下载图片
		strImgUrl := "https:" + item.ImgURL
		md5, err := utils.DownImage(strImgUrl)
		if md5 == "" || err != nil {
			fmt.Println("SinaSpider util.DownImage ", err.Error())
			continue
		}

		newInfo := models.Image{
			MD5:   md5,
			Url:   strImgUrl,
			Title: item.Name,
			State: 0,
		}

		err = db.AppendImage(newInfo)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// 启动
func (p *SinaSpider) Start(spiderInfo models.SpiderInfo) {
	p.bExit = false
	max := spiderInfo.MaxGoroutine
	if max <= 0 || max >= 20 {
		max = 10
	}
	p.maxChan = make(chan int, max)

	for {
		fmt.Println("SinaSpider开始采集")
		for pageIndex := 1; pageIndex <= spiderInfo.MaxPage; pageIndex++ {
			p.maxChan <- 1
			go p.collect(spiderInfo, pageIndex)
			time.Sleep(2 * time.Second)
		}
		if p.bExit {
			break
		}

		time.Sleep(5 * time.Second)
	}
}

// 停止
func (p *SinaSpider) Stop() {
	p.bExit = true
	fmt.Println("SinaSpider采集完成")
}
