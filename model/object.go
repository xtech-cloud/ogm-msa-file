package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Object struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Bucket    string `gorm:"column:bucket;type:char(32);not null"`
	Filepath  string `gorm:"column:filepath;type:varchar(256);not null"`
	URL       string `gorm:"column:url;type:varchar(1024)"`
	MD5       string `gorm:"column:md5;type:char(32)"`
	Size      uint64 `gorm:"column:size"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrObjectExists = errors.New("object exists")
var ErrObjectNotFound = errors.New("object not found")

func (Object) TableName() string {
	return "msa_file_object"
}

type ObjectDAO struct {
	conn *Conn
}

type ObjectQuery struct {
	Filepath string
	MD5      string
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

func (this *ObjectDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Object{}).Count(&count).Error
	return count, err
}

func (this *ObjectDAO) CountOfBucket(_bucket string) (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("bucket = ?", _bucket).Count(&count).Error
	return count, err
}

func (this *ObjectDAO) CountOfMD5(_bucket string, _md5 string) (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("bucket = ? AND md5 = ?", _bucket, _md5).Count(&count).Error
	return count, err
}

func (this *ObjectDAO) Insert(_object *Object) error {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("filepath= ?", _object.Filepath).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrObjectExists
	}

	return this.conn.DB.Create(_object).Error
}

func (this *ObjectDAO) Update(_object *Object) error {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("filepath = ?", _object.Filepath).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrObjectNotFound
	}

	return this.conn.DB.Updates(_object).Error
}

func (this *ObjectDAO) Delete(_filepath string) error {
	var count int64
	err := this.conn.DB.Model(&Object{}).Where("filepath = ?", _filepath).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrObjectNotFound
	}

	return this.conn.DB.Where("filepath = ?", _filepath).Delete(&Object{}).Error
}

func (this *ObjectDAO) List(_offset int64, _count int64) ([]*Object, error) {
	var objects []*Object
	res := this.conn.DB.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&objects)
	return objects, res.Error
}

func (this *ObjectDAO) QueryOne(_query *ObjectQuery) (*Object, error) {
	db := this.conn.DB.Model(&Object{})
	hasWhere := false
	if "" != _query.Filepath {
		db = db.Where("fillepath = ?", _query.Filepath)
		hasWhere = true
	}
	if "" != _query.MD5 {
		db = db.Where("md5 = ?", _query.MD5)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrObjectNotFound
	}

	var object Object
	err := db.First(&object).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrObjectNotFound
	}
	return &object, err
}
