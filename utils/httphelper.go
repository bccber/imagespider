package utils

import (
	"errors"
	"fmt"
	"imagespider/conf"
	"imagespider/db"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// 下载HTML
func DownUrlHtml(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return ""
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(buf)
}

// 下载图片保存到本地
func DownImage(downUrl string) (md5 string, err error) {
	if downUrl == "" {
		return "", errors.New("downUrl == ''")
	}
	md5 = MD5(downUrl)
	if md5 == "" {
		return "", errors.New("md5 == ''")
	}

	strExt := path.Ext(downUrl)
	saveFile := fmt.Sprintf("%s%c/%c/%s%s", conf.Config.SavePath, md5[0], md5[1], md5, strExt)
	if PathIsExist(saveFile) {
		// 如果文件存在，判断数据库是否存在
		b, err := db.CheckMD5(md5)
		if b && err != nil {
			// 如果数据库存在，则不需要下载
			return md5, nil
		}

		// 如果数据库不存在，则返回md5和nil，重新入库
		return md5, nil
	}

	savePath := filepath.Dir(saveFile)
	if !PathIsExist(savePath) {
		os.MkdirAll(savePath, os.ModePerm)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", downUrl, nil)

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return "", err
	}

	defer resp.Body.Close()
	file, _ := os.Create(saveFile)
	io.Copy(file, resp.Body)

	return md5, nil
}
