package mysql

import (
	"bluebell/setting"
	"fmt"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

// Init 初始化MySQL连接
//func Init(cfg *setting.MySQLConfig) (err error) {
//	// "user:password@tcp(host:port)/dbname"
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
//	db, err = sqlx.Connect("mysql", dsn)
//	if err != nil {
//		return
//	}
//	db.SetMaxOpenConns(cfg.MaxOpenConns)
//	db.SetMaxIdleConns(cfg.MaxIdleConns)
//	return
//}

func Init(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s", "./server_nodes.db")
	db, err = sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}
