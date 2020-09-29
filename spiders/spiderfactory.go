package spiders

import (
	"imagespider/models"
	"strings"
)

type ISpider interface {
	Start(models.SpiderInfo)
	Stop()
}

// 工厂方法
func CreateSpider(name string) ISpider {
	switch strings.ToLower(name) {
	case "sina":
		return &SinaSpider{}
	case "qiushibaike":
		return &QiuShiBaiKeSpider{}
	default:
		return nil
	}
}
