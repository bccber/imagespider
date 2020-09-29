package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"imagespider/conf"
	"imagespider/models"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", conf.Config.MySQLConn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

// 增加图片
func AppendImage(info models.Image) error {
	sql := `INSERT IGNORE INTO t_images(md5,url,title,state) VALUE(?,?,?,?);`
	_, err := db.Exec(sql, info.MD5, info.Url, info.Title, info.State)

	return err
}

// 判断图片是否存在
func CheckMD5(md5 string) (bool, error) {
	sql := `SELECT COUNT(0) FROM t_images WHERE md5=?;`
	var count int64
	err := db.QueryRow(sql, md5).Scan(&count)

	return count > 0, err
}
