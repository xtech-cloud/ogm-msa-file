package model

import (
	"errors"
	"gorm.io/gorm"
)

type Object struct {
	Embedded gorm.Model `gorm:"embedded"`
	Filepath string     `gorm:"column:filepath;type:varchar(256);not null;unique"`
	URL      string     `gorm:"column:url;type:varchar(1024)"`
	MD5      string     `gorm:"column:md5;type:char(32)"`
	Size     uint64     `gorm:"column:size"`
}

var ErrObjectExits = errors.New("object exists")

func (Object) TableName() string {
	return "msa_file_object"
}

type ObjectDAO struct {
    conn *Conn
}

func NewObjectDAO(_conn *Conn) *ObjectDAO {
    conn := DefaultConn
    if nil != _conn {
        conn = _conn
    }
    return &ObjectDAO{
        conn: conn,
    }
}
