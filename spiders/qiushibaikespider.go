package spiders

import "imagespider/models"

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"imagespider/db"
	"imagespider/utils"
	"strings"
	"time"
)

// 糗事百科爬虫
type QiuShiBaiKeSpider struct {
	bExit   bool
	maxChan chan int
}

// 采集
func (p *QiuShiBaiKeSpider) collect(spiderInfo models.SpiderInfo, pageIndex int) {
	defer func() {
		<-p.maxChan
		fmt.Printf("QiuShiBaiKeSpider 第 %d 页采集完毕\n", pageIndex)
	}()

	fmt.Printf("QiuShiBaiKeSpider 正在采集第 %d 页\n", pageIndex)

	strUrl := fmt.Sprintf(spiderInfo.URL, pageIndex)
	doc, err := htmlquery.LoadURL(strUrl)
	if err != nil {
		return
	}

	divs := htmlquery.Find(doc, `//div[starts-with(@id,"qiushi_tag_")]`)
	for _, div := range divs {

		content := htmlquery.FindOne(div, `.//div[@class="content"]//span`)
		if content == nil {
			continue
		}

		img := htmlquery.FindOne(div, `.//img[@class="illustration"]`)
		if img == nil {
			continue
		}

		strTitle := htmlquery.InnerText(content)
		strTitle = strings.TrimSpace(strTitle)

		imgUrl := htmlquery.SelectAttr(img, "src")

		// 下载图片
		strImgUrl := "https:" + imgUrl
		md5, err := utils.DownImage(strImgUrl)
		if md5 == "" || err != nil {
			fmt.Println("QiuShiBaiKeSpider util.DownImage ", err.Error())
			continue
		}

		newInfo := models.Image{
			MD5:   md5,
			Url:   strImgUrl,
			Title: strTitle,
			State: 0,
		}

		err = db.AppendImage(newInfo)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// 启动
func (p *QiuShiBaiKeSpider) Start(spiderInfo models.SpiderInfo) {
	p.bExit = false
	max := spiderInfo.MaxGoroutine
	if max <= 0 || max >= 20 {
		max = 10
	}
	p.maxChan = make(chan int, max)

	for {
		fmt.Println("QiuShiBaiKeSpider开始采集")
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
func (p *QiuShiBaiKeSpider) Stop() {
	p.bExit = true
	fmt.Println("QiuShiBaiKeSpider采集完成")
}
