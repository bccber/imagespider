package models

type SpiderInfo struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	MaxGoroutine int    `json:"maxGoroutine"`
	MaxPage      int    `json:"maxPage"`
	URL          string `json:"url"`
}

type Image struct {
	Id    int
	MD5   string
	Url   string
	Title string
	State int
}

type SinaJsonItem struct {
	Name   string `json:"name"`
	ImgURL string `json:"img_url"`
}

type SinaJsonInfo struct {
	Data []SinaJsonItem `json:"data"`
}
