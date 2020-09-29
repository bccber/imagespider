package conf

import (
	"encoding/json"
	"imagespider/models"
	"io/ioutil"
)

var Config *conf

type conf struct {
	MySQLConn string              `json:"mysqlConn"`
	SavePath  string              `json:"savePath"`
	Spider    []models.SpiderInfo `json:"spider"`
}

func init() {
	buf, err := ioutil.ReadFile("./conf/conf.json")
	if err != nil {
		panic("conf.json 有误")
	}

	if err = json.Unmarshal(buf, &Config); err != nil {
		panic(err)
	}
}
