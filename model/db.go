package model

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"ogm-msa-file/config"

	"github.com/asim/go-micro/v3/logger"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var base64Coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

type Conn struct {
	DB *gorm.DB
}

var DefaultConn *Conn

func Setup() {
	var err error
	var db *gorm.DB

	if "sqlite" == config.Schema.Database.Driver {
		dsn := config.Schema.Database.SQLite.Path
		logger.Warnf("!!! Database is lite mode, file at %v", dsn)
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else if "mysql" == config.Schema.Database.Driver {
		mysql_addr := config.Schema.Database.MySQL.Address
		mysql_user := config.Schema.Database.MySQL.User
		mysql_passwd := config.Schema.Database.MySQL.Password
		mysql_db := config.Schema.Database.MySQL.DB
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", mysql_user, mysql_passwd, mysql_addr, mysql_db)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		logger.Fatal("the driver of database is missing")
	}

	if nil != err {
		logger.Fatal(err)
	}
	DefaultConn = &Conn{
		DB: db,
	}
}

func Cancel() {
}

func AutoMigrateDatabase() {
	err := DefaultConn.DB.AutoMigrate(&Bucket{})
	if nil != err {
		logger.Fatal(err)
	}
	err = DefaultConn.DB.AutoMigrate(&Object{})
	if nil != err {
		logger.Fatal(err)
	}
}

func NewUUID() string {
	guid := uuid.NewV4()
	h := md5.New()
	h.Write(guid.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func ToUUID(_content string) string {
	h := md5.New()
	h.Write([]byte(_content))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(_content string) string {
	h := md5.New()
	h.Write([]byte(_content))
	return hex.EncodeToString(h.Sum(nil))
}

func ToBase64(_content []byte) string {
	return base64Coder.EncodeToString(_content)
}
